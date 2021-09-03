package rsync

import (
	"context"
	"os/exec"
)

func GetCommand(
	ctx context.Context,
	configPath string,
	sourceName string,
	sourceFilePath string,
	destinationName string,
	destinationFilePath string,
) *exec.Cmd {
	return exec.CommandContext(
		ctx,
		"rclone",
		"--no-check-certificate",
		"--config="+configPath,
		"sync",
		sourceName+":"+sourceFilePath,
		destinationName+":"+destinationFilePath,
	)
}
