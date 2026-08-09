package projects

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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

type DeploymentStatus string

const (
	StatusQueued    DeploymentStatus = "queued"
	StatusBuilding  DeploymentStatus = "building"
	StatusDeploying DeploymentStatus = "deploying"
	StatusRunning   DeploymentStatus = "running"
	StatusFailed    DeploymentStatus = "failed"
)

type Deployment struct {
	ID          int              `json:"id"`
	ProjectID   int              `json:"project_id"`
	Status      DeploymentStatus `json:"status"`
	CommitSha   string           `json:"commit_sha"`
	ImageName   string           `json:"image_name"`
	ContainerID string           `json:"container_id"`
	StartedAt   string           `json:"started_at"`
	FinishedAt  string           `json:"finished_at"`
	CreatedAt   string           `json:"created_at"`
}

func NewHandler(db *sql.DB) *Handler {
	// You can store the db reference in the handler if needed
	return &Handler{db: db}
}

func normalizeDeploymentStatus(status DeploymentStatus) DeploymentStatus {
	switch status {
	case StatusQueued, StatusBuilding, StatusDeploying, StatusRunning, StatusFailed:
		return status
	default:
		return StatusQueued
	}
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= '0' && r <= '9':
			return r
		default:
			return '-'
		}
	}, value)
	value = strings.Trim(value, "-")
	if value == "" {
		return "app"
	}
	return value
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

	// Apply defaults
	if input.Branch == "" {
		input.Branch = "main"
	}

	if input.RootDirectory == "" {
		input.RootDirectory = "/home/cursor/.fleetstack/projects/"
	}

	if input.ClonePath == "" {
		input.ClonePath = "/home/cursor/.fleetstack/projects/" + slugify(input.Name)
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

func (h *Handler) DeployProjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Fetch the project details from the database using the provided ID
	var project Project
	err := h.db.QueryRow("SELECT id, name, repo_url, branch, root_directory, clone_path, created_at FROM projects WHERE id = $1", id).
		Scan(&project.ID, &project.Name, &project.RepoURL, &project.Branch, &project.RootDirectory, &project.ClonePath, &project.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "project not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to fetch project: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Create a deployment for the project and return it immediately.
	var input Deployment

	// Decode JSON body into the input struct, but allow empty bodies to default to queued.
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil && err != io.EOF {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	status := normalizeDeploymentStatus(input.Status)

	query := `INSERT INTO deployments (project_id, status, commit_sha, image_name, container_id, started_at, finished_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`

	var deploymentID int
	var createdAt string

	err = h.db.QueryRow(
		query,
		project.ID,
		status,
		input.CommitSha,
		input.ImageName,
		input.ContainerID,
		input.StartedAt,
		input.FinishedAt,
	).Scan(&deploymentID, &createdAt)

	if err != nil {
		http.Error(w, "failed to create deployment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	deployment := Deployment{
		ID:          deploymentID,
		ProjectID:   project.ID,
		Status:      status,
		CommitSha:   input.CommitSha,
		ImageName:   input.ImageName,
		ContainerID: input.ContainerID,
		StartedAt:   input.StartedAt,
		FinishedAt:  input.FinishedAt,
		CreatedAt:   createdAt,
	}

	// Return the build output as a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(deployment); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) ProcessPendingDeployments() {
	for {
		deployment, err := h.nextQueuedDeployment()
		if err != nil {
			log.Printf("failed to fetch queued deployment: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if deployment == nil {
			time.Sleep(2 * time.Second)
			continue
		}

		if err := h.processDeployment(deployment); err != nil {
			log.Printf("deployment %d failed: %v", deployment.ID, err)
		}
	}
}

func (h *Handler) nextQueuedDeployment() (*Deployment, error) {
	var deployment Deployment
	err := h.db.QueryRow(`
		SELECT id, project_id, status, commit_sha, image_name, container_id, started_at, finished_at, created_at
		FROM deployments
		WHERE status = $1
		ORDER BY created_at ASC, id ASC
		LIMIT 1
	`, StatusQueued).Scan(
		&deployment.ID,
		&deployment.ProjectID,
		&deployment.Status,
		&deployment.CommitSha,
		&deployment.ImageName,
		&deployment.ContainerID,
		&deployment.StartedAt,
		&deployment.FinishedAt,
		&deployment.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	_, err = h.db.Exec(`UPDATE deployments SET status=$1 WHERE id=$2`, StatusBuilding, deployment.ID)
	if err != nil {
		return nil, err
	}
	deployment.Status = StatusBuilding
	return &deployment, nil
}

func (h *Handler) processDeployment(deployment *Deployment) error {
	project, err := h.getProjectByID(deployment.ProjectID)
	if err != nil {
		return fmt.Errorf("fetch project: %w", err)
	}

	projectPath := project.ClonePath
	if projectPath == "" {
		projectPath = filepath.Join(os.Getenv("HOME"), ".fleetstack", "projects", slugify(project.Name))
	}
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		if mkdirErr := os.MkdirAll(projectPath, 0o755); mkdirErr != nil {
			return fmt.Errorf("create project path: %w", mkdirErr)
		}
	}

	if _, err := tools.RailpackBuild(projectPath); err != nil {
		if updateErr := h.updateDeploymentStatus(deployment.ID, StatusFailed, "", ""); updateErr != nil {
			log.Printf("failed to update deployment %d after build failure: %v", deployment.ID, updateErr)
		}
		return fmt.Errorf("railpack build: %w", err)
	}

	if err := h.updateDeploymentStatus(deployment.ID, StatusDeploying, "", ""); err != nil {
		return fmt.Errorf("update deployment status to deploying: %w", err)
	}

	imageName := deployment.ImageName
	if imageName == "" {
		imageName = fmt.Sprintf("fleetstack/%s:%d", slugify(project.Name), deployment.ID)
	}
	containerName := fmt.Sprintf("fleetstack-%s-%d", slugify(project.Name), deployment.ID)
	cmd := exec.Command("docker", "run", "-d", "--name", containerName, imageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if updateErr := h.updateDeploymentStatus(deployment.ID, StatusFailed, imageName, ""); updateErr != nil {
			log.Printf("failed to update deployment %d after container failure: %v", deployment.ID, updateErr)
		}
		return fmt.Errorf("docker run failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	containerID := strings.TrimSpace(string(output))
	if containerID == "" {
		containerID = containerName
	}
	if err := h.updateDeploymentStatus(deployment.ID, StatusRunning, imageName, containerID); err != nil {
		return fmt.Errorf("update deployment status to running: %w", err)
	}
	return nil
}

func (h *Handler) getProjectByID(projectID int) (*Project, error) {
	var project Project
	err := h.db.QueryRow("SELECT id, name, repo_url, branch, root_directory, clone_path, created_at FROM projects WHERE id = $1", projectID).
		Scan(&project.ID, &project.Name, &project.RepoURL, &project.Branch, &project.RootDirectory, &project.ClonePath, &project.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (h *Handler) updateDeploymentStatus(id int, status DeploymentStatus, imageName, containerID string) error {
	_, err := h.db.Exec(`
		UPDATE deployments
		SET status = $1, image_name = COALESCE(NULLIF($2, ''), image_name), container_id = COALESCE(NULLIF($3, ''), container_id), finished_at = CASE WHEN $1 IN ('failed', 'running') THEN NOW() ELSE finished_at END
		WHERE id = $4
	`, status, imageName, containerID, id)
	return err
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
