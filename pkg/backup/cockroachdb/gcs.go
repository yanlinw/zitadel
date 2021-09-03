package cockroachdb

import (
	"context"
	"encoding/base64"
	"os/exec"
	"strings"
)

func GetBackupToGCS(
	ctx context.Context,
	certsFolder string,
	host string,
	port string,
	bucketName string,
	filePath string,
	serviceAccount []byte,
) *exec.Cmd {
	saStr := strings.ReplaceAll(base64.StdEncoding.EncodeToString(serviceAccount), "\n", "")
	return exec.CommandContext(
		ctx,
		"cockroach",
		"sql",
		"--certs-dir="+certsFolder,
		"--host="+host,
		"--port="+port,
		`-e`,
		`BACKUP TO 'gs://`+bucketName+`/`+filePath+`?AUTH=specified&CREDENTIALS=`+saStr+`';`,
	)
}

func GetRestoreFromGCS(
	ctx context.Context,
	certsFolder string,
	host string,
	port string,
	bucketName string,
	filePath string,
	serviceAccount []byte,
) *exec.Cmd {
	saStr := strings.ReplaceAll(base64.StdEncoding.EncodeToString(serviceAccount), "\n", "")
	return exec.CommandContext(
		ctx,
		"cockroach",
		"sql",
		"--certs-dir="+certsFolder,
		"--host="+host,
		"--port="+port,
		"-e",
		`RESTORE FROM 'gs://`+bucketName+`/`+filePath+`?AUTH=specified&CREDENTIALS=`+saStr+`';`,
	)
}
