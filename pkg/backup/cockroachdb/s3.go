package cockroachdb

import (
	"context"
	"os/exec"
	"strings"
)

func GetBackupToS3(
	ctx context.Context,
	certsFolder string,
	host string,
	port string,
	bucketName string,
	filePath string,
	accessKeyID []byte,
	secretAccessKey []byte,
	sessionToken []byte,
	endpoint string,
	region string,
) *exec.Cmd {
	parameters := []string{
		"AWS_ACCESS_KEY_ID=" + string(accessKeyID),
		"AWS_SECRET_ACCESS_KEY=" + string(secretAccessKey),
		"AWS_SESSION_TOKEN=" + string(sessionToken),
		"AWS_ENDPOINT=" + endpoint,
	}
	if region != "" {
		parameters = append(parameters, "AWS_REGION="+region)
	}

	return exec.CommandContext(
		ctx,
		"cockroach",
		"sql",
		"--certs-dir="+certsFolder,
		"--host="+host,
		"--port="+port,
		"-e",
		`BACKUP TO 's3://`+bucketName+`/`+filePath+`?`+strings.Join(parameters, "&")+`';`,
	)
}

func GetRestoreFromS3(
	ctx context.Context,
	certsFolder string,
	host string,
	port string,
	bucketName string,
	filePath string,
	accessKeyID []byte,
	secretAccessKey []byte,
	sessionToken []byte,
	endpoint string,
	region string,
) *exec.Cmd {
	parameters := []string{
		"AWS_ACCESS_KEY_ID=" + string(accessKeyID),
		"AWS_SECRET_ACCESS_KEY=" + string(secretAccessKey),
		"AWS_SESSION_TOKEN=" + string(sessionToken),
		"AWS_ENDPOINT=" + endpoint,
	}
	if region != "" {
		parameters = append(parameters, "AWS_REGION="+region)
	}

	return exec.CommandContext(
		ctx,
		"cockroach",
		"sql",
		"--certs-dir="+certsFolder,
		"--host="+host,
		"--port="+port,
		"-e",
		`RESTORE FROM 's3://`+bucketName+`/`+filePath+`?`+strings.Join(parameters, "&")+`';`,
	)
}
