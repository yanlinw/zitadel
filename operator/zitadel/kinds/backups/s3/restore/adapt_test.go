package restore

import (
	core2 "github.com/caos/zitadel/operator/zitadel/kinds/backups/core"
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

func TestBackup_Adapt1(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	monitor := mntr.Monitor{}
	namespace := "testNs"
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{
		{Key: "testKey", Operator: "testOp"}}
	timestamp := "testTs"
	backupName := "testName2"
	bucketName := "testBucket2"
	sourceAKIDName := "testAKIDN"
	sourceAKIDKey := "testAKIDK"
	sourceSAKName := "testSAKN"
	sourceSAKKey := "testSAK"
	destAKIDName := "testAKIDN"
	destAKIDKey := "testAKID"
	destSAKName := "testSAKN"
	destSAKKey := "testSAK"
	sourceEndpoint := "endpoint"
	destEndpoint := "endpoint"
	sourcePrefix := "prefix"
	sourceRegion := "region"
	destRegion := "region"
	dbURL := "testDB"
	dbPort := int32(80)
	image := "testImage"
	jobName := GetJobName(backupName)
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "testKind", "testVersion"), "testComponent")
	nameLabels := labels.MustForName(componentLabels, jobName)

	checkDBReady := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}

	jobDef := core2.GetJob(
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
			image,
			getCommand(
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
		bucketName,
		timestamp,
		sourceAKIDName,
		sourceAKIDKey,
		sourceSAKName,
		sourceSAKKey,
		destAKIDName,
		destAKIDKey,
		destSAKName,
		destSAKKey,
		nodeselector,
		tolerations,
		checkDBReady,
		dbURL,
		dbPort,
		image,
		sourceEndpoint,
		destEndpoint,
		sourcePrefix,
		sourceRegion,
		destRegion,
	)

	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))
}

func TestBackup_Adapt2(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	monitor := mntr.Monitor{}
	namespace := "testNs2"
	nodeselector := map[string]string{"test2": "test2"}
	tolerations := []corev1.Toleration{
		{Key: "testKey2", Operator: "testOp2"}}
	timestamp := "testTs"
	backupName := "testName2"
	bucketName := "testBucket2"
	sourceAKIDName := "testAKIDN2"
	sourceAKIDKey := "testAKIDK2"
	sourceSAKName := "testSAKN2"
	sourceSAKKey := "testSAK2"
	destAKIDName := "testAKIDN2"
	destAKIDKey := "testAKID2"
	destSAKName := "testSAKN2"
	destSAKKey := "testSAK2"
	sourceEndpoint := "endpoint"
	destEndpoint := "endpoint"
	sourcePrefix := "prefix"
	sourceRegion := "region"
	destRegion := "region"
	dbURL := "testDB"
	dbPort := int32(80)
	image := "testImage"
	jobName := GetJobName(backupName)
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testVersion2"), "testKind2", "testVersion2"), "testComponent2")
	nameLabels := labels.MustForName(componentLabels, jobName)

	checkDBReady := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}

	jobDef := core2.GetJob(
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
			image,
			getCommand(
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
		bucketName,
		timestamp,
		sourceAKIDName,
		sourceAKIDKey,
		sourceSAKName,
		sourceSAKKey,
		destAKIDName,
		destAKIDKey,
		destSAKName,
		destSAKKey,
		nodeselector,
		tolerations,
		checkDBReady,
		dbURL,
		dbPort,
		image,
		sourceEndpoint,
		destEndpoint,
		sourcePrefix,
		sourceRegion,
		destRegion,
	)

	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))
}
