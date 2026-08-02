package github

import (
	"context"
	"log"
	"net/http"
	"time"
	"os/exec"
	"path"
	"path/filepath"

	"fmt"
	"io"
	"os"

	gogithub "github.com/google/go-github/v89/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/isaacmain254/fleetstack/backend/internal/config"
	"github.com/isaacmain254/fleetstack/backend/internal/nixpacks"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, getRedirectURL(), http.StatusTemporaryRedirect)
}

func (h *Handler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	callbackHandler(w, r)
}


func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	// state := r.URL.Query().Get("state")
	conf := getOAuthConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Printf("unable to exchange code for token: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	repos, err := getCurrentUserRepo(token.AccessToken)
	if err != nil {
		log.Printf("unable to get user repositories: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, repo := range repos {
		log.Printf("Repository: %s, URL: %s", *repo.Name, *repo.HTMLURL)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Repositories fetched successfully. Check server logs for details."))

}

func getOAuthConfig() *oauth2.Config {
	conf := config.Load()
	return &oauth2.Config{
		ClientID:     conf.GitHub.ClientID,
		ClientSecret: conf.GitHub.ClientSecret,
		// RedirectURL:  "http://localhost:8080/callback",
		Scopes:   []string{"user:email", "repo:public_repo"},
		Endpoint: github.Endpoint,
	}
}

func getRedirectURL() string {
	config := getOAuthConfig()
	authURL := config.AuthCodeURL("state")
	return authURL
}


func getCurrentUserRepo(accessToken string) ([]*gogithub.Repository, error) {
	log.Printf("Access Token: %s", accessToken)
	client, err := gogithub.NewClient(gogithub.WithAuthToken(accessToken))
	if err != nil {
		log.Printf("unable to create GitHub client: %v", err)
		return nil, err
	}
	opt := &gogithub.RepositoryListByAuthenticatedUserOptions{Affiliation: "owner", Sort: "updated", Direction: "desc"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	repos, _, err := client.Repositories.ListByAuthenticatedUser(ctx, opt)
	if err != nil {
		log.Printf("unable to list repositories: %v", err)
		return nil, err
	}
	return repos, nil
}


func (h *Handler) CloneRepoHandler(w http.ResponseWriter, r *http.Request) {
	repoURL := r.FormValue("repo_url")
	if repoURL == "" {
		http.Error(w, "Missing repo_url parameter", http.StatusBadRequest)
		return
	}
	output, err := cloneRepo(repoURL)
	if err != nil {
		http.Error(w, "Failed to clone repository: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf(`{"message": "Repository cloned successfully.", "output": %q}`, output))
}

func cloneRepo(repoURL string) (string, error) {
	// store clone repository in FLEETSTACK_STORAGE=/opt/fleetstack
	storageRoot := os.Getenv("FLEETSTACK_STORAGE")
	if storageRoot == "" {
		home, _ := os.UserHomeDir()
		storageRoot = filepath.Join(home, ".fleetstack")
	}

	repoName := path.Base(repoURL)
	clonePath := filepath.Join(
		storageRoot,
		"projects",
		repoName,
	)
	err := os.MkdirAll(clonePath, os.ModePerm)
	if err != nil {
		log.Printf("unable to create directory for cloning: %v", err)
		return "", err
	}
	// tmpDir, err := os.MkdirTemp("", "fleetstack_projects_*")
	// if err != nil {
	// 	log.Printf("unable to create temporary directory: %v", err)
	// 	return "", err
	// }
	// clonePath := path.Join(tmpDir, repoName)
	// cloned repository should be stored in backend/tmp folder
	// clonePath := path.Join("tmp", path.Base(repoURL))
	// err := os.MkdirAll(clonePath, os.ModePerm)
	// if err != nil {
	// 	log.Printf("unable to create directory for cloning: %v", err)
	// 	return "", err
	// }
	cmd := exec.Command("git", "clone", repoURL, clonePath)
	err = cmd.Run()
	if err != nil {
		log.Printf("unable to clone repository: %v", err)
		return "", err
	}

	// cmd = exec.Command(
	// 	"nixpacks",
	// 	"plan",
	// 	clonePath,
	// )

	// output, err := cmd.CombinedOutput()

	// fmt.Println(string(output))
	// fmt.Println(err)

	// return string(output), nil
	output, err := nixpacks.NixPacksPlan(clonePath)
	if err != nil {
		log.Printf("unable to run nixpacks plan: %v", err)
		return "", err
	}

	return string(output), nil
}