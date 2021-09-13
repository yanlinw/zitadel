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
	assetEndpoint string,
	assetAKIDPath string,
	assetSAKPath string,
	assetPrefix string,
	assetRegion string,
	destinationAKIDPath string,
	destinationSAKPath string,
	destinationEndpoint string,
	destinationRegion string,
	dbURL string,
	dbPort int32,
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
			"s3",
			"--backupname=" + backupName,
			"--backupnameenv=" + backupNameEnv,
			"--asset-endpoint=" + assetEndpoint,
			"--asset-akid=" + assetAKIDPath,
			"--asset-sak=" + assetSAKPath,
			"--asset-prefix=" + assetPrefix,
			"--asset-region=" + assetRegion,
			"--destination-endpoint=" + destinationEndpoint,
			"--destination-akid=" + destinationAKIDPath,
			"--destination-sak=" + destinationSAKPath,
			"--destination-bucket=" + bucketName,
			"--destination-region=" + destinationRegion,
			"--host=" + dbURL,
			"--port=" + strconv.Itoa(int(dbPort)),
			"--certs-dir=" + certsFolder,
			"--configpath=/rsync.conf",
		}, " ",
		),
	)

	return strings.Join(backupCommands, " && ")
}
