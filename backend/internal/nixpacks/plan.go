package nixpacks

import (
	"fmt"
	"os/exec"
	"encoding/json"
)

type NixpacksPlan struct {
	Variables  map[string]string `json:"variables"`
	BuildImage string            `json:"buildImage"`
	Phases     struct {
		Install struct {
			Cmds []string `json:"cmds"`
		} `json:"install"`

		Build struct {
			Cmds []string `json:"cmds"`
		} `json:"build"`
	} `json:"phases"`

	Start struct {
		Cmd string `json:"cmd"`
	} `json:"start"`
}

type BuildInfo struct {
	Message        string `json:"message"`
	Framework      string `json:"framework"`
	BuildImage     string `json:"buildImage"`
	InstallCommand string `json:"installCommand"`
	BuildCommand   string `json:"buildCommand"`
	StartCommand   string `json:"startCommand"`
}

func NixPacksPlan(projectPath string) (output []byte, err error) {
	cmd := exec.Command("nixpacks", "plan", projectPath)

	output, err = cmd.CombinedOutput()

	var plan NixpacksPlan

	if err := json.Unmarshal(output, &plan); err != nil {
		return nil, err
	}

	fmt.Println(string(output))
	fmt.Println(err)

	return output, err
}
