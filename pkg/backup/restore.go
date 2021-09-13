package backup

import (
	"context"
	"github.com/caos/zitadel/pkg/backup/cockroachdb"
	"github.com/caos/zitadel/pkg/backup/rsync"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

func RsyncRestoreS3ToS3(
	ctx context.Context,
	backupName string,
	backupNameEnv string,
	destinationName string,
	destinationEndpoint string,
	destinationAKIDPath string,
	destinationSAKPath string,
	sourceName string,
	sourceEndpoint string,
	sourceAKIDPath string,
	sourceSAKPath string,
	sourceBucket string,
	sourcePrefix string,
	configFilePath string,
) error {

	sourceAKID, err := ioutil.ReadFile(sourceAKIDPath)
	if err != nil {
		return err
	}

	sourceSAK, err := ioutil.ReadFile(sourceSAKPath)
	if err != nil {
		return err
	}

	destinationAKID, err := ioutil.ReadFile(destinationAKIDPath)
	if err != nil {
		return err
	}

	destinationSAK, err := ioutil.ReadFile(destinationSAKPath)
	if err != nil {
		return err
	}

	assetBuckets, err := ListS3Folders(sourceEndpoint, string(sourceAKID), string(sourceSAK), sourceBucket, getAssetPath(backupName, backupNameEnv))
	if err != nil {
		return err
	}

	sourcePart, err := rsync.GetConfigPartS3(sourceName, sourceEndpoint, string(sourceAKID), string(sourceSAK))
	if err != nil {
		return err
	}

	destPart, err := rsync.GetConfigPartS3(destinationName, destinationEndpoint, string(destinationAKID), string(destinationSAK))
	if err != nil {
		return err
	}

	config := strings.Join([]string{
		sourcePart,
		destPart,
	}, "\n")

	if err := ioutil.WriteFile(configFilePath, []byte(config), fs.ModePerm); err != nil {
		return err
	}

	var wg sync.WaitGroup

	if err := CleanS3BucketsWithPrefix(destinationEndpoint, string(destinationAKID), string(destinationSAK), sourcePrefix); err != nil {
		return err
	}

	if err := EnsureS3BucketsWithPrefix(destinationEndpoint, string(destinationAKID), string(destinationSAK), assetBuckets, sourcePrefix); err != nil {
		return err
	}

	for _, assetBucket := range assetBuckets {
		wg.Add(1)
		if err := runCommand(
			rsync.GetCommand(
				ctx,
				configFilePath,
				sourceName,
				filepath.Join(getAssetFullPath(sourceBucket, backupName, backupNameEnv), assetBucket),
				destinationName,
				assetBucket,
			),
			&wg,
		); err != nil {
			return err
		}
	}
	wg.Wait()

	return nil
}

func RsyncRestoreGCSToS3(
	ctx context.Context,
	backupName string,
	backupNameEnv string,
	destinationName string,
	destinationEndpoint string,
	destinationAKIDPath string,
	destinationSAKPath string,
	sourceName string,
	sourceSaJsonPath string,
	sourceBucket string,
	sourcePrefix string,
	configFilePath string,
) error {

	destinationAKID, err := ioutil.ReadFile(destinationAKIDPath)
	if err != nil {
		return err
	}

	destinationSAK, err := ioutil.ReadFile(destinationSAKPath)
	if err != nil {
		return err
	}

	assetBuckets, err := ListGCSFolders(sourceSaJsonPath, sourceBucket, getAssetPath(backupName, backupNameEnv))
	if err != nil {
		return err
	}

	sourcePart, err := rsync.GetConfigPartGCS(sourceName, sourceSaJsonPath)
	if err != nil {
		return err
	}

	destPart, err := rsync.GetConfigPartS3(destinationName, destinationEndpoint, string(destinationAKID), string(destinationSAK))
	if err != nil {
		return err
	}

	config := strings.Join([]string{
		sourcePart,
		destPart,
	}, "\n")

	if err := ioutil.WriteFile(configFilePath, []byte(config), fs.ModePerm); err != nil {
		return err
	}

	var wg sync.WaitGroup

	if err := CleanS3BucketsWithPrefix(destinationEndpoint, string(destinationAKID), string(destinationSAK), sourcePrefix); err != nil {
		return err
	}

	if err := EnsureS3BucketsWithPrefix(destinationEndpoint, string(destinationAKID), string(destinationSAK), assetBuckets, sourcePrefix); err != nil {
		return err
	}

	for _, assetBucket := range assetBuckets {
		wg.Add(1)
		if err := runCommand(
			rsync.GetCommand(
				ctx,
				configFilePath,
				sourceName,
				filepath.Join(getAssetFullPath(sourceBucket, backupName, backupNameEnv), assetBucket),
				destinationName,
				assetBucket,
			),
			&wg,
		); err != nil {
			return err
		}
	}
	wg.Wait()

	return nil
}

func CockroachRestoreFromGCS(
	ctx context.Context,
	certsFolder string,
	bucketName string,
	backupName string,
	backupNameEnv string,
	host string,
	port string,
	serviceAccountPath string,
) error {
	var wg sync.WaitGroup

	wg.Add(1)

	data, err := ioutil.ReadFile(serviceAccountPath)
	if err != nil {
		return err
	}
	err = runCommand(
		cockroachdb.GetRestoreFromGCS(
			ctx,
			certsFolder,
			host,
			port,
			bucketName,
			getBackupPath(backupName, backupNameEnv),
			data,
		),
		&wg,
	)
	wg.Wait()
	return err
}

func CockroachRestoreFromS3(
	ctx context.Context,
	certsFolder string,
	bucketName string,
	backupName string,
	backupNameEnv string,
	host string,
	port string,
	accessKeyIDPath string,
	secretAccessKeyPath string,
	sessionTokenPath string,
	endpoint string,
	region string,
) error {
	var wg sync.WaitGroup

	wg.Add(1)

	accessKeyID, err := ioutil.ReadFile(accessKeyIDPath)
	if err != nil {
		return err
	}

	secretAccessKey, err := ioutil.ReadFile(secretAccessKeyPath)
	if err != nil {
		return err
	}

	sessionToken, err := ioutil.ReadFile(sessionTokenPath)
	if err != nil {
		return err
	}

	err = runCommand(
		cockroachdb.GetRestoreFromS3(
			ctx,
			certsFolder,
			host,
			port,
			bucketName,
			getBackupPath(backupName, backupNameEnv),
			accessKeyID,
			secretAccessKey,
			sessionToken,
			endpoint,
			region,
		),
		&wg,
	)
	wg.Wait()
	return err
}
