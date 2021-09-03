package main

import (
	"context"
	"github.com/caos/zitadel/cmd/backupctl/cmds/backup"
	"github.com/caos/zitadel/cmd/backupctl/cmds/restore"
	"os"
	"os/signal"

	"github.com/caos/orbos/mntr"

	"github.com/caos/zitadel/cmd/backupctl/cmds"
)

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	monitor := mntr.Monitor{
		OnInfo:         mntr.LogMessage,
		OnChange:       mntr.LogMessage,
		OnError:        mntr.LogError,
		OnRecoverPanic: mntr.LogPanic,
	}

	defer func() { monitor.RecoverPanic(recover()) }()

	rootCmd := cmds.RootCommand()
	backupCmd := cmds.BackupCommand(monitor)
	backupCmd.AddCommand(
		backup.S3Command(ctx, monitor),
		backup.GCSCommand(ctx, monitor),
	)
	rootCmd.AddCommand(backupCmd)
	restoreCmd := cmds.RestoreCommand(monitor)
	restoreCmd.AddCommand(
		restore.S3Command(ctx, monitor),
		restore.GCSCommand(ctx, monitor),
	)
	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(cmds.RestoreCommand(monitor))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
