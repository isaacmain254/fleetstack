package nixpacks

import (
	"os/exec"
)

// type BuildInfo struct {
// 	Message        string `json:"message"`
// 	Framework      string `json:"framework"`
// 	BuildImage     string `json:"buildImage"`
// 	InstallCommand string `json:"installCommand"`
// 	BuildCommand   string `json:"buildCommand"`
// 	StartCommand   string `json:"startCommand"`
// }


func NixPacksBuild(projectPath string) (output []byte, err error) {
	cmd := exec.Command("nixpacks", "build", projectPath)

	output, err = cmd.CombinedOutput()

	return output, err
}	