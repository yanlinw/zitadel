package backup

import (
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/core"
	core2 "github.com/caos/zitadel/operator/zitadel/kinds/backups/s3/core"
	"time"

	"github.com/caos/zitadel/operator"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/cronjob"
	"github.com/caos/orbos/pkg/kubernetes/resources/job"
	"github.com/caos/orbos/pkg/labels"
	corev1 "k8s.io/api/core/v1"
)

const (
	backupNameEnv     = "BACKUP_NAME"
	cronJobNamePrefix = "backup-"
	timeout           = 15 * time.Minute
	Normal            = "backup"
	Instant           = "instantbackup"
)

func AdaptFunc(
	monitor mntr.Monitor,
	backupName string,
	namespace string,
	componentLabels *labels.Component,
	checkDBReady operator.EnsureFunc,
	bucketName string,
	cron string,
	sourceAKIDName string,
	sourceAKIDKey string,
	sourceSAKName string,
	sourceSAKKey string,
	destAKIDName string,
	destAKIDKey string,
	destSAKName string,
	destSAKKey string,
	timestamp string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	dbURL string,
	dbPort int32,
	features []string,
	image string,
	sourceEndpoint string,
	sourcePrefix string,
	destEndpoint string,
	sourceRegion string,
	destRegion string,
) (
	queryFunc operator.QueryFunc,
	destroyFunc operator.DestroyFunc,
	err error,
) {

	command := getBackupCommand(
		timestamp,
		bucketName,
		backupName,
		core2.CertPath,
		sourceEndpoint,
		core2.SourceAkidSecretPath,
		core2.SourceSakSecretPath,
		sourcePrefix,
		sourceRegion,
		core2.DestAkidSecretPath,
		core2.DestSakSecretPath,
		destEndpoint,
		destRegion,
		dbURL,
		dbPort,
	)

	jobSpecDef := core2.GetJobSpecDef(
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
	)

	destroyers := []operator.DestroyFunc{}
	queriers := []operator.QueryFunc{}

	cronJobDef := core.GetCronJob(
		namespace,
		labels.MustForName(componentLabels, GetJobName(backupName)),
		cron,
		jobSpecDef,
	)

	destroyCJ, err := cronjob.AdaptFuncToDestroy(cronJobDef.Namespace, cronJobDef.Name)
	if err != nil {
		return nil, nil, err
	}

	queryCJ, err := cronjob.AdaptFuncToEnsure(cronJobDef)
	if err != nil {
		return nil, nil, err
	}

	jobDef := core.GetJob(
		namespace,
		labels.MustForName(componentLabels, cronJobNamePrefix+backupName),
		jobSpecDef,
	)

	destroyJ, err := job.AdaptFuncToDestroy(jobDef.Namespace, jobDef.Name)
	if err != nil {
		return nil, nil, err
	}

	queryJ, err := job.AdaptFuncToEnsure(jobDef)
	if err != nil {
		return nil, nil, err
	}

	for _, feature := range features {
		switch feature {
		case Normal:
			destroyers = append(destroyers,
				operator.ResourceDestroyToZitadelDestroy(destroyCJ),
			)
			queriers = append(queriers,
				operator.EnsureFuncToQueryFunc(checkDBReady),
				operator.ResourceQueryToZitadelQuery(queryCJ),
			)
		case Instant:
			destroyers = append(destroyers,
				operator.ResourceDestroyToZitadelDestroy(destroyJ),
			)
			queriers = append(queriers,
				operator.EnsureFuncToQueryFunc(checkDBReady),
				operator.ResourceQueryToZitadelQuery(queryJ),
			)
		}
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(monitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(monitor, destroyers),
		nil
}

func GetJobName(backupName string) string {
	return cronJobNamePrefix + backupName
}
