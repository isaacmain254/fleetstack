package tools

import (
	"fmt"
	"os/exec"
	// "bufio"
	// "io"
)


func RailpackPlan(projectPath string) (output []byte, err error) {
	cmd := exec.Command("railpack", "plan", projectPath)

	output, err = cmd.CombinedOutput()

	return output, err
}


// func RailpackBuild(projectPath string) (output []byte, err error) {
// 	cmd := exec.Command("railpack", "build", projectPath)
// 	fmt.Println("Running command:", cmd.String())

// 	output, err = cmd.CombinedOutput()

// 	return output, err
// }

func RailpackBuild(projectPath string) (output []byte, err error) {
	cmd := exec.Command("railpack", "build", projectPath)
	fmt.Println("Running command:", cmd.String())

	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("railpack output:\n%s\n", output)
	}

	  // Get pipes for real-time reading
    // stdout, _ := cmd.StdoutPipe()
    // stderr, _ := cmd.StderrPipe()
    
    // // Combine streams if desired
    // stream := io.MultiReader(stdout, stderr)
    
    // cmd.Start()
    
    // scanner := bufio.NewScanner(stream)

    // for scanner.Scan() {
    //     fmt.Println(scanner.Text()) 
	// 	// Output appears in real-time
	// 	if err := scanner.Err(); err != nil {
	// 		fmt.Printf("scanner error: %v\n", err)
	// 	}
    // }
    
    // cmd.Wait()

	return output, err
}