package backup

import (
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/s3/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackup_Command1(t *testing.T) {
	timestamp := ""
	bucketName := "test"
	backupName := "test"
	dbURL := "testDB"
	dbPort := int32(80)
	sourceEndpoint := "endpoint"
	sourcePrefix := "prefix"
	destEndpoint := "endpoint"
	sourceRegion := "region1"
	destRegion := "region2"

	cmd := getBackupCommand(
		timestamp,
		bucketName,
		backupName,
		core.CertPath,
		sourceEndpoint,
		core.SourceAkidSecretPath,
		core.SourceSakSecretPath,
		sourcePrefix,
		sourceRegion,
		core.DestAkidSecretPath,
		core.DestSakSecretPath,
		destEndpoint,
		destRegion,
		dbURL,
		dbPort,
	)
	equals := "export BACKUP_NAME=$(date +%Y-%m-%dT%H:%M:%SZ) && /backupctl backup s3 --backupname=test --backupnameenv=BACKUP_NAME --asset-endpoint=endpoint --asset-akid=/secrets/sakid --asset-sak=/secrets/ssak --asset-prefix=prefix --asset-region=region1 --destination-endpoint=endpoint --destination-akid=/secrets/dakid --destination-sak=/secrets/dsak --destination-bucket=test --destination-region=region2 --host=testDB --port=80 --certs-dir=/cockroach/cockroach-certs --configpath=/rsync.conf"
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command2(t *testing.T) {
	timestamp := "test"
	bucketName := "test"
	backupName := "test"
	dbURL := "testDB"
	dbPort := int32(80)
	sourceEndpoint := "endpoint"
	sourcePrefix := "prefix"
	destEndpoint := "endpoint"
	sourceRegion := "region3"
	destRegion := "region4"

	cmd := getBackupCommand(
		timestamp,
		bucketName,
		backupName,
		core.CertPath,
		sourceEndpoint,
		core.SourceAkidSecretPath,
		core.SourceSakSecretPath,
		sourcePrefix,
		sourceRegion,
		core.DestAkidSecretPath,
		core.DestSakSecretPath,
		destEndpoint,
		destRegion,
		dbURL,
		dbPort,
	)
	equals := "export BACKUP_NAME=test && /backupctl backup s3 --backupname=test --backupnameenv=BACKUP_NAME --asset-endpoint=endpoint --asset-akid=/secrets/sakid --asset-sak=/secrets/ssak --asset-prefix=prefix --asset-region=region3 --destination-endpoint=endpoint --destination-akid=/secrets/dakid --destination-sak=/secrets/dsak --destination-bucket=test --destination-region=region4 --host=testDB --port=80 --certs-dir=/cockroach/cockroach-certs --configpath=/rsync.conf"
	assert.Equal(t, equals, cmd)
}
