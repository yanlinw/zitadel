package backup

import (
	core2 "github.com/caos/zitadel/operator/zitadel/kinds/backups/bucket/core"
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/core"
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
	image := "testImage"
	secretName := "testSecretName"
	saSecretKey := "testKey"
	akidKey := "testAKID"
	sakKey := "testSAK"
	endpoint := "testEndpoint"
	prefix := "testPrefix"
	region := "testRegion"
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
		core2.GetJobSpecDef(
			nodeselector,
			tolerations,
			secretName,
			saSecretKey,
			akidKey,
			sakKey,
			backupName,
			getBackupCommand(
				timestamp,
				bucketName,
				backupName,
				core2.CertPath,
				core2.SaSecretPath,
				core2.AkidSecretPath,
				core2.SakSecretPath,
				dbURL,
				dbPort,
				endpoint,
				prefix,
				region,
			),
			image,
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
		secretName,
		saSecretKey,
		akidKey,
		sakKey,
		timestamp,
		nodeselector,
		tolerations,
		dbURL,
		dbPort,
		features,
		image,
		endpoint,
		prefix,
		region,
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
	image := "testImage2"
	saSecretKey := "testKey2"
	akidKey := "testAKID2"
	sakKey := "testSAK2"
	endpoint := "testEndpoint2"
	prefix := "testPrefix2"
	secretName := "testSecretName2"
	region := "testRegion"
	jobName := GetJobName(backupName)
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testVersion2"), "testKind2", "testVersion2"), "testComponent")
	nameLabels := labels.MustForName(componentLabels, jobName)

	checkDBReady := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}

	jobDef := core.GetJob(
		namespace,
		nameLabels,
		core2.GetJobSpecDef(
			nodeselector,
			tolerations,
			secretName,
			saSecretKey,
			akidKey,
			sakKey,
			backupName,
			getBackupCommand(
				timestamp,
				bucketName,
				backupName,
				core2.CertPath,
				core2.SaSecretPath,
				core2.AkidSecretPath,
				core2.SakSecretPath,
				dbURL,
				dbPort,
				endpoint,
				prefix,
				region,
			),
			image,
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
		secretName,
		saSecretKey,
		akidKey,
		sakKey,
		timestamp,
		nodeselector,
		tolerations,
		dbURL,
		dbPort,
		features,
		image,
		endpoint,
		prefix,
		region,
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
	image := "testImage"
	saSecretKey := "testKey"
	akidKey := "testAKID"
	sakKey := "testSAK"
	endpoint := "testEndpoint"
	prefix := "testPrefix"
	secretName := "testSecretName"
	region := "testRegion"
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
		core2.GetJobSpecDef(
			nodeselector,
			tolerations,
			secretName,
			saSecretKey,
			akidKey,
			sakKey,
			backupName,
			getBackupCommand(
				timestamp,
				bucketName,
				backupName,
				core2.CertPath,
				core2.SaSecretPath,
				core2.AkidSecretPath,
				core2.SakSecretPath,
				dbURL,
				dbPort,
				endpoint,
				prefix,
				region,
			),
			image,
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
		secretName,
		saSecretKey,
		akidKey,
		sakKey,
		timestamp,
		nodeselector,
		tolerations,
		dbURL,
		dbPort,
		features,
		image,
		endpoint,
		prefix,
		region,
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
	image := "testImage2"
	saSecretKey := "testKey2"
	akidKey := "testAKID2"
	sakKey := "testSAK2"
	endpoint := "testEndpoint2"
	prefix := "testPrefix2"
	secretName := "testSecretName2"
	region := "testRegion"
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
		core2.GetJobSpecDef(
			nodeselector,
			tolerations,
			secretName,
			saSecretKey,
			akidKey,
			sakKey,
			backupName,
			getBackupCommand(
				timestamp,
				bucketName,
				backupName,
				core2.CertPath,
				core2.SaSecretPath,
				core2.AkidSecretPath,
				core2.SakSecretPath,
				dbURL,
				dbPort,
				endpoint,
				prefix,
				region,
			),
			image,
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
		secretName,
		saSecretKey,
		akidKey,
		sakKey,
		timestamp,
		nodeselector,
		tolerations,
		dbURL,
		dbPort,
		features,
		image,
		endpoint,
		prefix,
		region,
	)

	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))
}
