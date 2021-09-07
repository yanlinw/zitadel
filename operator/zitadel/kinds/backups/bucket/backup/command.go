package backup

import (
	"strconv"
	"strings"
)

func getBackupCommand(
	timestamp string,
	bucketName string,
	backupName string,
	certsFolder string,
	serviceAccountPath string,
	dbURL string,
	dbPort int32,
	assetEndpoint string,
	assetPrefix string,
) string {

	backupCommands := make([]string, 0)
	if timestamp != "" {
		backupCommands = append(backupCommands, "export "+backupNameEnv+"="+timestamp)
	} else {
		backupCommands = append(backupCommands, "export "+backupNameEnv+"=$(date +%Y-%m-%dT%H:%M:%SZ)")
	}

	backupCommands = append(backupCommands,
		strings.Join([]string{
			"/backupctl",
			"backup",
			"gcs",
			"--backupname=" + backupName,
			"--backupnameenv=" + backupNameEnv,
			"--asset-endpoint=" + assetEndpoint,
			"--asset-akid=" + akidSecretPath,
			"--asset-sak=" + sakSecretPath,
			"--asset-prefix=" + assetPrefix,
			"--host=" + dbURL,
			"--port=" + strconv.Itoa(int(dbPort)),
			"--destination-sajsonpath=" + serviceAccountPath,
			"--destination-bucket=" + bucketName,
			"--certs-dir=" + certsFolder,
			"--configpath=/rsync.conf",
		}, " ",
		),
	)

	return strings.Join(backupCommands, " && ")
}
