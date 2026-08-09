package projects

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"github.com/gorilla/mux"

	"github.com/isaacmain254/fleetstack/backend/internal/tools"
)

type Handler struct {
	db *sql.DB
}

type Project struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	RepoURL       string `json:"repo_url"`
	Branch        string `json:"branch"`
	RootDirectory string `json:"root_directory"`
	ClonePath     string `json:"clone_path"`
	CreatedAt     string `json:"created_at"`
}

func NewHandler(db *sql.DB) *Handler {
	// You can store the db reference in the handler if needed
	return &Handler{db: db}
}

// create a new project
func (h *Handler) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var input Project

	// Decode JSON body into the input struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Basic validation (optional but recommended)
	if input.Name == "" || input.RepoURL == "" {
		http.Error(w, "name and repo_url are required", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO projects (name, repo_url, branch, root_directory, clone_path) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	var id int
	var createdAt string

	err := h.db.QueryRow(
		query,
		input.Name,
		input.RepoURL,
		input.Branch,
		input.RootDirectory,
		input.ClonePath,
	).Scan(&id, &createdAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	project := Project{
		ID:            id,
		Name:          input.Name,
		RepoURL:       input.RepoURL,
		Branch:        input.Branch,
		RootDirectory: input.RootDirectory,
		ClonePath:     input.ClonePath,
		CreatedAt:     createdAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(project); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) ListProjectsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`SELECT id, name, repo_url, branch, root_directory, clone_path, created_at FROM projects`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.RepoURL, &p.Branch, &p.RootDirectory, &p.ClonePath, &p.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetProjectByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rows, err := h.db.Query("SELECT id, name, repo_url, branch, root_directory, clone_path, created_at FROM projects WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	if rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.RepoURL, &p.Branch, &p.RootDirectory, &p.ClonePath, &p.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "project not found", http.StatusNotFound)
	}
}


// Project Path should be dynamic
func (h *Handler) BuildProjectHandler(w http.ResponseWriter, r *http.Request) {
	// home directory: ~/.fleetstack/projects
	_, err := BuildProject("/home/cursor/.fleetstack/projects/Workpay")
	if err != nil {
		log.Printf("Failed to Build project: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"message": "Project built successfully."}`)
}

func BuildProject(projectPath string) (output []byte, err error) {

	output, err = tools.NixPacksBuild(projectPath)
	if err != nil {
		log.Printf("Error building the project: %v", err)
		return nil, err
	}
	return output, nil
}
