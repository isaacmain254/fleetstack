package projects

import (
	"io"
	"log"
	"net/http"

	"github.com/isaacmain254/fleetstack/backend/internal/nixpacks"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
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

// type Project struct {
// 	name       string
// 	repo_url   string
// 	created_at string
// }

// func createProject(ctx context.Context, db *sql.DB, project Project) error {
// 	query := `
// 		INSERT INTO projects (
// 			name, repo_url, created_at
// 		) VALUES ($1, $2, $3)
// 	`

// 	_, err := db.ExecContext(ctx, query, project.name, project.repo_url, project.created_at)

// 	if err != nil {
// 		log.Printf("Error creating a project: %v", err)
// 		return err
// 	}
// 	return nil
// }

func  BuildProject(projectPath string) (output []byte, err error) {

	output, err = nixpacks.NixPacksBuild(projectPath)
	if err != nil {
		log.Printf("Error building the project: %v", err)
		return nil, err
	}
	return output, nil
}