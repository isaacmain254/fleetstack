package tools

import (
	"fmt"
	"os/exec"
	"encoding/json"
)


// type NixpacksPlan struct {
// 	Variables  map[string]string `json:"variables"`
// 	BuildImage string            `json:"buildImage"`
// 	Phases     struct {
// 		Install struct {
// 			Cmds []string `json:"cmds"`
// 		} `json:"install"`

// 		Build struct {
// 			Cmds []string `json:"cmds"`
// 		} `json:"build"`
// 	} `json:"phases"`

// 	Start struct {
// 		Cmd string `json:"cmd"`
// 	} `json:"start"`
// }

type NixpacksPlan struct {
	Variables  map[string]string `json:"variables"`
	
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


func NixPacksBuild(projectPath string) (output []byte, err error) {
	cmd := exec.Command("nixpacks", "build", projectPath, "--name my-website", "--env NIXPACKS_NODE_VERSION=22")
	fmt.Println("Running command:", cmd.String())

	output, err = cmd.CombinedOutput()

	return output, err
}	