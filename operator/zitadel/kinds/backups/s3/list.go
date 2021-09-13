package s3

import (
	"fmt"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/secret/read"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/core"
	"github.com/caos/zitadel/pkg/backup"
)

func BackupList() core.BackupListFunc {
	return func(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, name string, desired *tree.Tree) ([]string, error) {
		desiredKind, err := ParseDesiredV0(desired)
		if err != nil {
			return nil, fmt.Errorf("parsing desired state failed: %s", err)
		}
		desired.Parsed = desiredKind

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			monitor.Verbose()
		}

		valueAKI, err := read.GetSecretValue(k8sClient, desiredKind.Spec.AccessKeyID, desiredKind.Spec.ExistingAccessKeyID)
		if err != nil {
			return nil, err
		}

		valueSAK, err := read.GetSecretValue(k8sClient, desiredKind.Spec.SecretAccessKey, desiredKind.Spec.ExistingSecretAccessKey)
		if err != nil {
			return nil, err
		}

		return backup.S3ListFilesWithFilter(
			valueAKI,
			valueSAK,
			//desiredKind.Spec.Region,
			desiredKind.Spec.Endpoint,
			desiredKind.Spec.Bucket,
			name,
			desiredKind.Spec.Region,
		)
	}
}
