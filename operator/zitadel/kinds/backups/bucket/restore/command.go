package restore

import (
	"strconv"
	"strings"
)

func getCommand(
	timestamp string,
	bucketName string,
	backupName string,
	certsFolder string,
	serviceAccountPath string,
	dbURL string,
	dbPort int32,
	assetEndpoint string,
	assetPrefix string,
	assetRegion string,
) string {

	backupCommands := make([]string, 0)
	backupCommands = append(backupCommands, "export "+backupNameEnv+"="+timestamp)

	backupCommands = append(backupCommands,
		strings.Join([]string{
			"/backupctl",
			"restore",
			"gcs",
			"--backupname=" + backupName,
			"--backupnameenv=" + backupNameEnv,
			"--asset-endpoint=" + assetEndpoint,
			"--asset-akid=" + akidSecretPath,
			"--asset-sak=" + sakSecretPath,
			"--asset-region=" + assetRegion,
			"--host=" + dbURL,
			"--port=" + strconv.Itoa(int(dbPort)),
			"--source-sajsonpath=" + serviceAccountPath,
			"--source-bucket=" + bucketName,
			"--source-prefix=" + assetPrefix,
			"--certs-dir=" + certsFolder,
			"--configpath=/rsync.conf",
		}, " ",
		),
	)

	return strings.Join(backupCommands, " && ")
}
