package restore

import (
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/s3/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackup_Command1(t *testing.T) {
	timestamp := "test1"
	bucketName := "testBucket"
	backupName := "testBackup"
	dbURL := "testDB"
	dbPort := int32(80)
	sourceEndpoint := "endpoint"
	destEndpoint := "endpoint"
	sourcePrefix := "prefix"
	sourceRegion := "region"
	destRegion := "region"

	cmd := getCommand(
		timestamp,
		bucketName,
		backupName,
		core.CertPath,
		destEndpoint,
		core.DestAkidSecretPath,
		core.DestSakSecretPath,
		destRegion,
		core.SourceAkidSecretPath,
		core.SourceSakSecretPath,
		sourceEndpoint,
		sourcePrefix,
		sourceRegion,
		dbURL,
		dbPort,
	)

	equals := "export BACKUP_NAME=test1 && /backupctl restore s3 --backupname=testBackup --backupnameenv=BACKUP_NAME --asset-endpoint=endpoint --asset-akid=/secrets/dakid --asset-sak=/secrets/dsak --asset-region=region --source-endpoint=endpoint --source-akid=/secrets/sakid --source-sak=/secrets/ssak --source-bucket=testBucket --source-prefix=prefix --source-region=region --host=testDB --port=80 --certs-dir=/cockroach/cockroach-certs --configpath=/rsync.conf"
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command2(t *testing.T) {
	timestamp := "test2"
	bucketName := "testBucket"
	backupName := "testBackup"
	dbURL := "testDB2"
	dbPort := int32(81)
	sourceEndpoint := "endpoint2"
	destEndpoint := "endpoint2"
	sourcePrefix := "prefix2"
	sourceRegion := "region2"
	destRegion := "region2"

	cmd := getCommand(
		timestamp,
		bucketName,
		backupName,
		core.CertPath,
		destEndpoint,
		core.DestAkidSecretPath,
		core.DestSakSecretPath,
		destRegion,
		core.SourceAkidSecretPath,
		core.SourceSakSecretPath,
		sourceEndpoint,
		sourcePrefix,
		sourceRegion,
		dbURL,
		dbPort,
	)
	equals := "export BACKUP_NAME=test2 && /backupctl restore s3 --backupname=testBackup --backupnameenv=BACKUP_NAME --asset-endpoint=endpoint2 --asset-akid=/secrets/dakid --asset-sak=/secrets/dsak --asset-region=region2 --source-endpoint=endpoint2 --source-akid=/secrets/sakid --source-sak=/secrets/ssak --source-bucket=testBucket --source-prefix=prefix2 --source-region=region2 --host=testDB2 --port=81 --certs-dir=/cockroach/cockroach-certs --configpath=/rsync.conf"
	assert.Equal(t, equals, cmd)
}
