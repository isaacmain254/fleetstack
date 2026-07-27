package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	gogithub "github.com/google/go-github/v65/github"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)
var clientID string
var clientSecret string
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	r := mux.NewRouter()
	r.HandleFunc("/redirect", redirectHandler)
	r.HandleFunc("/callback", callbackHandler)
	r.HandleFunc("/clone", handleCloneRepo).Methods("POST")

	log.Println("Starting server on :8080")

	log.Fatal(http.ListenAndServe(":8080", r))

}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, getRedirectURL(), http.StatusTemporaryRedirect)
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
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		// RedirectURL:  "http://localhost:8080/callback",
		Scopes:      []string{"user:email", "repo:public_repo"},
		Endpoint:    github.Endpoint,
	}
}

func getRedirectURL() string {
	config := getOAuthConfig()
	authURL := config.AuthCodeURL("state")
	return authURL
}

func getCurrentUserRepo(accessToken string) ([]*gogithub.Repository, error) {
	client := gogithub.NewClient(nil).WithAuthToken(accessToken)
	opt := &gogithub.RepositoryListByAuthenticatedUserOptions{Affiliation: "owner"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	repos, _, err := client.Repositories.ListByAuthenticatedUser(ctx, opt)
	if err != nil {
		log.Printf("unable to list repositories: %v", err)
		return nil, err
	}
	return repos, nil
}

func handleCloneRepo(w http.ResponseWriter, r *http.Request) {
	repoURL := r.FormValue("repo_url")
	if repoURL == "" {
		http.Error(w, "Missing repo_url parameter", http.StatusBadRequest)
		return
	}
	err := cloneRepo(repoURL)
	if err != nil {
		http.Error(w, "Failed to clone repository", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Repository cloned successfully."))
}

func cloneRepo(repoURL string) error {
	cmd := exec.Command("git", "clone", repoURL)
	err := cmd.Run()
	if err != nil {
		log.Printf("unable to clone repository: %v", err)
		return err
	}
	return nil
}
