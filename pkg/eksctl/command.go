package eksctl

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func Command(args []string, stdin string) error {
	if err := CheckVersion(); err != nil {
		return err
	}

	command := exec.Command("eksctl", args...)
	command.Stdin = strings.NewReader(string(stdin))
	stdout, _ := command.StdoutPipe()
	stderr, _ := command.StderrPipe()

	command.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	slurp, _ := io.ReadAll(stderr)
	err := command.Wait()
	if err != nil {
		return fmt.Errorf("eksctl failed: %s", string(slurp))
	}

	return nil
}

func CommandWithResult(args []string, stdin string) (string, error) {
	if err := CheckVersion(); err != nil {
		return "", err
	}

	command := exec.Command("eksctl", args...)
	command.Stdin = strings.NewReader(string(stdin))
	stdout, _ := command.StdoutPipe()
	stderr, _ := command.StderrPipe()

	command.Start()

	result, _ := io.ReadAll(stdout)
	slurp, _ := io.ReadAll(stderr)
	err := command.Wait()
	if err != nil {
		return "", fmt.Errorf("eksctl failed: %s", string(slurp))
	}

	return string(result), nil
}
