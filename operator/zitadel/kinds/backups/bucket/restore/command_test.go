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
	prefix := "testPrefix"
	region := "testRegion"

	cmd := getCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		saSecretPath,
		dbURL,
		dbPort,
		enpoint,
		prefix,
		region,
	)

	equals := "export BACKUP_NAME=test1 && /backupctl restore gcs --backupname=testBackup --backupnameenv=BACKUP_NAME --asset-endpoint=testEndpoint --asset-akid=/secrets/akid --asset-sak=/secrets/sak --asset-region=testRegion --host=testDB --port=80 --source-sajsonpath=/secrets/sa.json --source-bucket=testBucket --source-prefix=testPrefix --certs-dir=/cockroach/cockroach-certs --configpath=/rsync.conf"
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command2(t *testing.T) {
	timestamp := "test2"
	bucketName := "testBucket2"
	backupName := "testBackup2"
	dbURL := "testDB2"
	dbPort := int32(81)
	enpoint := "testEndpoint2"
	prefix := "testPrefix2"
	region := "testRegion2"

	cmd := getCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		saSecretPath,
		dbURL,
		dbPort,
		enpoint,
		prefix,
		region,
	)
	equals := "export BACKUP_NAME=test2 && /backupctl restore gcs --backupname=testBackup2 --backupnameenv=BACKUP_NAME --asset-endpoint=testEndpoint2 --asset-akid=/secrets/akid --asset-sak=/secrets/sak --asset-region=testRegion2 --host=testDB2 --port=81 --source-sajsonpath=/secrets/sa.json --source-bucket=testBucket2 --source-prefix=testPrefix2 --certs-dir=/cockroach/cockroach-certs --configpath=/rsync.conf"
	assert.Equal(t, equals, cmd)
}
