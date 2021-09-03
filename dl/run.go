package dl

import (
	"fmt"
	"os"
	"os/exec"
)

func Run(procs map[string][]string) error {
	for proc, args := range procs {
		cmd := exec.Command(proc, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			list := append([]string{proc}, args...)
			return fmt.Errorf("failed when executing: %v -> %v", list, err)
		}
	}
	return nil
}
