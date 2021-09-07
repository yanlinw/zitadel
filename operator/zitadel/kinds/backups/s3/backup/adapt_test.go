package backup

import (
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/s3/core"
	"testing"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	macherrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestBackup_AdaptInstantBackup1(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	features := []string{Instant}
	monitor := mntr.Monitor{}
	namespace := "testNs"

	bucketName := "testBucket"
	cron := "testCron"
	timestamp := "test"
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{
		{Key: "testKey", Operator: "testOp"}}
	backupName := "testName"
	version := "testVersion"
	sourceAKIDName := "testAKIDN"
	sourceAKIDKey := "testAKIDK"
	sourceSAKName := "testSAKN"
	sourceSAKKey := "testSAK"
	destAKIDName := "testAKIDN"
	destAKIDKey := "testAKID"
	destSAKName := "testSAKN"
	destSAKKey := "testSAK"
	sourcePrefix := "prefix"
	sourceEndpoint := "endpoint"
	destEndpoint := "endpoint"
	dbURL := "testDB"
	dbPort := int32(80)
	jobName := GetJobName(backupName)
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testVersion2"), "testKind2", "testVersion2"), "testComponent")
	nameLabels := labels.MustForName(componentLabels, jobName)

	checkDBReady := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}

	jobDef := core.GetJob(
		namespace,
		nameLabels,
		core.GetJobSpecDef(
			nodeselector,
			tolerations,
			sourceAKIDName,
			sourceAKIDKey,
			sourceSAKName,
			sourceSAKKey,
			destAKIDName,
			destAKIDKey,
			destSAKName,
			destSAKKey,
			backupName,
			version,
			getBackupCommand(
				timestamp,
				bucketName,
				backupName,
				core.CertPath,
				sourceEndpoint,
				core.SourceAkidSecretPath,
				core.SourceSakSecretPath,
				sourcePrefix,
				core.DestAkidSecretPath,
				core.DestSakSecretPath,
				destEndpoint,
				dbURL,
				dbPort,
			),
		),
	)

	client.EXPECT().ApplyJob(jobDef).Times(1).Return(nil)
	client.EXPECT().GetJob(jobDef.Namespace, jobDef.Name).Times(1).Return(nil, macherrs.NewNotFound(schema.GroupResource{"batch", "jobs"}, jobName))

	query, _, err := AdaptFunc(
		monitor,
		backupName,
		namespace,
		componentLabels,
		checkDBReady,
		bucketName,
		cron,
		sourceAKIDName,
		sourceAKIDKey,
		sourceSAKName,
		sourceSAKKey,
		destAKIDName,
		destAKIDKey,
		destSAKName,
		destSAKKey,
		timestamp,
		nodeselector,
		tolerations,
		dbURL,
		dbPort,
		features,
		version,
		sourceEndpoint,
		sourcePrefix,
		destEndpoint,
	)

	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))
}

func TestBackup_AdaptInstantBackup2(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	features := []string{Instant}
	monitor := mntr.Monitor{}
	namespace := "testNs2"
	dbURL := "testDB"
	dbPort := int32(80)
	bucketName := "testBucket2"
	cron := "testCron2"
	timestamp := "test2"
	nodeselector := map[string]string{"test2": "test2"}
	tolerations := []corev1.Toleration{
		{Key: "testKey2", Operator: "testOp2"}}
	backupName := "testName2"
	version := "testVersion2"
	sourceAKIDName := "testAKIDN2"
	sourceAKIDKey := "testAKIDK2"
	sourceSAKName := "testSAKN2"
	sourceSAKKey := "testSAK2"
	destAKIDName := "testAKIDN2"
	destAKIDKey := "testAKID2"
	destSAKName := "testSAKN2"
	destSAKKey := "testSAK2"
	sourcePrefix := "prefix2"
	sourceEndpoint := "endpoint2"
	destEndpoint := "endpoint2"
	jobName := GetJobName(backupName)
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testVersion2"), "testKind2", "testVersion2"), "testComponent")
	nameLabels := labels.MustForName(componentLabels, jobName)

	checkDBReady := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}

	jobDef := core.GetJob(
		namespace,
		nameLabels,
		core.GetJobSpecDef(
			nodeselector,
			tolerations,
			sourceAKIDName,
			sourceAKIDKey,
			sourceSAKName,
			sourceSAKKey,
			destAKIDName,
			destAKIDKey,
			destSAKName,
			destSAKKey,
			backupName,
			version,
			getBackupCommand(
				timestamp,
				bucketName,
				backupName,
				core.CertPath,
				sourceEndpoint,
				core.SourceAkidSecretPath,
				core.SourceSakSecretPath,
				sourcePrefix,
				core.DestAkidSecretPath,
				core.DestSakSecretPath,
				destEndpoint,
				dbURL,
				dbPort,
			),
		),
	)

	client.EXPECT().ApplyJob(jobDef).Times(1).Return(nil)
	client.EXPECT().GetJob(jobDef.Namespace, jobDef.Name).Times(1).Return(nil, macherrs.NewNotFound(schema.GroupResource{"batch", "jobs"}, jobName))

	query, _, err := AdaptFunc(
		monitor,
		backupName,
		namespace,
		componentLabels,
		checkDBReady,
		bucketName,
		cron,
		sourceAKIDName,
		sourceAKIDKey,
		sourceSAKName,
		sourceSAKKey,
		destAKIDName,
		destAKIDKey,
		destSAKName,
		destSAKKey,
		timestamp,
		nodeselector,
		tolerations,
		dbURL,
		dbPort,
		features,
		version,
		sourceEndpoint,
		sourcePrefix,
		destEndpoint,
	)

	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))
}

func TestBackup_AdaptBackup1(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	features := []string{Normal}
	monitor := mntr.Monitor{}
	namespace := "testNs"
	bucketName := "testBucket"
	cron := "testCron"
	timestamp := "test"
	dbURL := "testDB"
	dbPort := int32(80)
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{
		{Key: "testKey", Operator: "testOp"}}
	backupName := "testName"
	version := "testVersion"
	sourceAKIDName := "testAKIDN"
	sourceAKIDKey := "testAKIDK"
	sourceSAKName := "testSAKN"
	sourceSAKKey := "testSAK"
	destAKIDName := "testAKIDN"
	destAKIDKey := "testAKID"
	destSAKName := "testSAKN"
	destSAKKey := "testSAK"
	sourcePrefix := "prefix"
	sourceEndpoint := "endpoint"
	destEndpoint := "endpoint"
	jobName := GetJobName(backupName)
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testVersion2"), "testKind2", "testVersion2"), "testComponent")
	nameLabels := labels.MustForName(componentLabels, jobName)

	checkDBReady := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}

	jobDef := core.GetCronJob(
		namespace,
		nameLabels,
		cron,
		core.GetJobSpecDef(
			nodeselector,
			tolerations,
			sourceAKIDName,
			sourceAKIDKey,
			sourceSAKName,
			sourceSAKKey,
			destAKIDName,
			destAKIDKey,
			destSAKName,
			destSAKKey,
			backupName,
			version,
			getBackupCommand(
				timestamp,
				bucketName,
				backupName,
				core.CertPath,
				sourceEndpoint,
				core.SourceAkidSecretPath,
				core.SourceSakSecretPath,
				sourcePrefix,
				core.DestAkidSecretPath,
				core.DestSakSecretPath,
				destEndpoint,
				dbURL,
				dbPort,
			),
		),
	)

	client.EXPECT().ApplyCronJob(jobDef).Times(1).Return(nil)

	query, _, err := AdaptFunc(
		monitor,
		backupName,
		namespace,
		componentLabels,
		checkDBReady,
		bucketName,
		cron,
		sourceAKIDName,
		sourceAKIDKey,
		sourceSAKName,
		sourceSAKKey,
		destAKIDName,
		destAKIDKey,
		destSAKName,
		destSAKKey,
		timestamp,
		nodeselector,
		tolerations,
		dbURL,
		dbPort,
		features,
		version,
		sourceEndpoint,
		sourcePrefix,
		destEndpoint,
	)

	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))
}

func TestBackup_AdaptBackup2(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	features := []string{Normal}
	monitor := mntr.Monitor{}
	namespace := "testNs2"
	dbURL := "testDB"
	dbPort := int32(80)
	bucketName := "testBucket2"
	cron := "testCron2"
	timestamp := "test2"
	nodeselector := map[string]string{"test2": "test2"}
	tolerations := []corev1.Toleration{
		{Key: "testKey2", Operator: "testOp2"}}
	backupName := "testName2"
	version := "testVersion2"
	sourceAKIDName := "testAKIDN2"
	sourceAKIDKey := "testAKIDK2"
	sourceSAKName := "testSAKN2"
	sourceSAKKey := "testSAK2"
	destAKIDName := "testAKIDN2"
	destAKIDKey := "testAKID2"
	destSAKName := "testSAKN2"
	destSAKKey := "testSAK2"
	sourcePrefix := "prefix2"
	sourceEndpoint := "endpoint2"
	destEndpoint := "endpoint2"
	jobName := GetJobName(backupName)
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testVersion2"), "testKind2", "testVersion2"), "testComponent")
	nameLabels := labels.MustForName(componentLabels, jobName)

	checkDBReady := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}

	jobDef := core.GetCronJob(
		namespace,
		nameLabels,
		cron,
		core.GetJobSpecDef(
			nodeselector,
			tolerations,
			sourceAKIDName,
			sourceAKIDKey,
			sourceSAKName,
			sourceSAKKey,
			destAKIDName,
			destAKIDKey,
			destSAKName,
			destSAKKey,
			backupName,
			version,
			getBackupCommand(
				timestamp,
				bucketName,
				backupName,
				core.CertPath,
				sourceEndpoint,
				core.SourceAkidSecretPath,
				core.SourceSakSecretPath,
				sourcePrefix,
				core.DestAkidSecretPath,
				core.DestSakSecretPath,
				destEndpoint,
				dbURL,
				dbPort,
			),
		),
	)

	client.EXPECT().ApplyCronJob(jobDef).Times(1).Return(nil)

	query, _, err := AdaptFunc(
		monitor,
		backupName,
		namespace,
		componentLabels,
		checkDBReady,
		bucketName,
		cron,
		sourceAKIDName,
		sourceAKIDKey,
		sourceSAKName,
		sourceSAKKey,
		destAKIDName,
		destAKIDKey,
		destSAKName,
		destSAKKey,
		timestamp,
		nodeselector,
		tolerations,
		dbURL,
		dbPort,
		features,
		version,
		sourceEndpoint,
		sourcePrefix,
		destEndpoint,
	)

	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))
}
