package tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	// "bufio"
	"io"
	"strings"
	"time"
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

func RailpackBuild(projectPath string, imageName string) (output []byte, builtImage string, err error) {
	start := time.Now()
	args := []string{"--verbose", "build", "--progress", "plain"}
	if strings.TrimSpace(imageName) != "" {
		args = append(args, "--name", imageName)
	}
	args = append(args, projectPath)
	cmd := exec.Command("railpack", args...)
	fmt.Println("Running command:", cmd.String())

	var buf bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &buf) // prints live AND captures

	cmd.Stdout = mw
	cmd.Stderr = mw

	err = cmd.Run()
	fmt.Printf("Railpack build finished in %s\n", time.Since(start).Round(time.Second))
	builtImage = parseRailpackLoadedImage(buf.String())

	return buf.Bytes(), builtImage, err
}

func parseRailpackLoadedImage(output string) string {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Loaded image:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Loaded image:"))
		}
	}
	return ""
}

// func RailpackBuild(projectPath string) (output []byte, err error) {
// 	cmd := exec.Command("railpack", "build", projectPath)
// 	fmt.Println("Running command:", cmd.String())

// 	output, err = cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Printf("railpack output:\n%s\n", output)
// 	}

// 	//   Get pipes for real-time reading
//     stdout, _ := cmd.StdoutPipe()
//     stderr, _ := cmd.StderrPipe()

//     // Combine streams if desired
//     stream := io.MultiReader(stdout, stderr)

//     cmd.Start()

//     scanner := bufio.NewScanner(stream)

//     for scanner.Scan() {
//         fmt.Println(scanner.Text())
// 		// Output appears in real-time
// 		if err := scanner.Err(); err != nil {
// 			fmt.Printf("scanner error: %v\n", err)
// 		}
//     }

//     cmd.Wait()

// 	return output, err
// }
