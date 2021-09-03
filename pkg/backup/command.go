package backup

import (
	"fmt"
	"os/exec"
	"sync"
)

func runCommand(cmd *exec.Cmd, wg *sync.WaitGroup) error {
	errs := make(chan error)
	go func() {
		out, err := cmd.CombinedOutput()
		if err != nil {
			errs <- fmt.Errorf("%s: %s", err, out)
			wg.Done()
			return
		}
		wg.Done()
		errs <- nil
	}()

	return <-errs
}
