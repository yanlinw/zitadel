package restore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackup_Command1(t *testing.T) {
	timestamp := "test1"
	bucketName := "testBucket"
	backupName := "testBackup"
	dbURL := "testDB"
	dbPort := int32(80)
	enpoint := "testEndpoint"

	cmd := getCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		saSecretPath,
		dbURL,
		dbPort,
		enpoint,
	)

	equals := "export BACKUP_NAME=test1 && /backupctl restore gcs --backupname=testBackup --backupnameenv=BACKUP_NAME --asset-endpoint=testEndpoint --asset-akid=/secrets/akid --asset-sak=/secrets/sak --host=testDB --port=80 --source-sajsonpath=/secrets/sa.json --source-bucket=testBucket --certs-dir=/cockroach/cockroach-certs --configpath=/rsync.conf"
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command2(t *testing.T) {
	timestamp := "test2"
	bucketName := "testBucket"
	backupName := "testBackup"
	dbURL := "testDB2"
	dbPort := int32(81)
	enpoint := "testEndpoint"

	cmd := getCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		saSecretPath,
		dbURL,
		dbPort,
		enpoint,
	)
	equals := "export BACKUP_NAME=test2 && /backupctl restore gcs --backupname=testBackup --backupnameenv=BACKUP_NAME --asset-endpoint=testEndpoint --asset-akid=/secrets/akid --asset-sak=/secrets/sak --host=testDB2 --port=81 --source-sajsonpath=/secrets/sa.json --source-bucket=testBucket --certs-dir=/cockroach/cockroach-certs --configpath=/rsync.conf"
	assert.Equal(t, equals, cmd)
}
