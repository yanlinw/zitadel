package restore

import (
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/s3/core"
	"time"

	"github.com/caos/zitadel/operator"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/job"
	"github.com/caos/orbos/pkg/labels"
	corev1 "k8s.io/api/core/v1"
)

const (
	Instant       = "restore"
	jobPrefix     = "backup-"
	jobSuffix     = "-restore"
	timeout       = 15 * time.Minute
	backupNameEnv = "BACKUP_NAME"
)

func AdaptFunc(
	monitor mntr.Monitor,
	backupName string,
	namespace string,
	componentLabels *labels.Component,
	bucketName string,
	timestamp string,
	sourceAKIDName string,
	sourceAKIDKey string,
	sourceSAKName string,
	sourceSAKKey string,
	destAKIDName string,
	destAKIDKey string,
	destSAKName string,
	destSAKKey string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	checkDBReady operator.EnsureFunc,
	dbURL string,
	dbPort int32,
	image string,
	sourceEndpoint string,
	destEndpoint string,
) (
	queryFunc operator.QueryFunc,
	destroyFunc operator.DestroyFunc,
	err error,
) {

	jobName := jobPrefix + backupName + jobSuffix
	command := getCommand(
		timestamp,
		bucketName,
		backupName,
		core.CertPath,
		destEndpoint,
		core.DestAkidSecretPath,
		core.DestSakSecretPath,
		core.SourceAkidSecretPath,
		core.SourceSakSecretPath,
		sourceEndpoint,
		dbURL,
		dbPort,
	)

	jobdef := core.GetJob(
		namespace,
		labels.MustForName(componentLabels, GetJobName(backupName)),
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
			command,
		),
	)

	destroyJ, err := job.AdaptFuncToDestroy(jobName, namespace)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyJ),
	}

	queryJ, err := job.AdaptFuncToEnsure(jobdef)
	if err != nil {
		return nil, nil, err
	}

	queriers := []operator.QueryFunc{
		operator.EnsureFuncToQueryFunc(checkDBReady),
		operator.ResourceQueryToZitadelQuery(queryJ),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(monitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(monitor, destroyers),

		nil
}

func GetJobName(backupName string) string {
	return jobPrefix + backupName + jobSuffix
}
