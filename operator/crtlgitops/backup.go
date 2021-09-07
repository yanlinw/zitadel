package crtlgitops

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/zitadel/pkg/databases"
	"github.com/caos/zitadel/pkg/zitadel"
)

func Restore(
	monitor mntr.Monitor,
	gitClient *git.Client,
	k8sClient *kubernetes.Client,
	orbconfig *orb.Orb,
	backup string,
) error {
	if err := databases.GitOpsClear(monitor, k8sClient, gitClient); err != nil {
		return err
	}

	if err := zitadel.GitOpsRestore(
		monitor,
		k8sClient,
		gitClient,
		orbconfig,
		backup,
	); err != nil {
		return err
	}
	return nil
}
