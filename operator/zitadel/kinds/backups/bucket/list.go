package bucket

import (
	"fmt"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/secret/read"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/pkg/backup"

	"github.com/caos/zitadel/operator/zitadel/kinds/backups/core"
)

func BackupList() core.BackupListFunc {
	return func(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, name string, desired *tree.Tree) ([]string, error) {
		desiredKind, err := ParseDesiredV0(desired)
		if err != nil {
			return nil, fmt.Errorf("parsing desired state failed: %w", err)
		}
		desired.Parsed = desiredKind

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			monitor.Verbose()
		}

		value, err := read.GetSecretValue(k8sClient, desiredKind.Spec.ServiceAccountJSON, desiredKind.Spec.ExistingServiceAccountJSON)
		if err != nil {
			return nil, err
		}

		return backup.GCSListFilesWithFilter(
			value,
			desiredKind.Spec.Bucket,
			name,
		)
	}
}
