package zitadel

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/zitadel"
	orbz "github.com/caos/zitadel/operator/zitadel/kinds/orb"
)

func GitOpsInstantBackup(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
	orbconfig *orb.Orb,
	name string,
) error {
	desired, err := gitClient.ReadTree(git.ZitadelFile)
	if err != nil {
		monitor.Error(err)
		return err
	}
	return instantBackup(monitor, k8sClient, desired, orbconfig, true, name)
}

func CrdInstantBackup(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	name string,
) error {
	desired, err := zitadel.ReadCrd(k8sClient)
	if err != nil {
		monitor.Error(err)
		return err
	}
	return instantBackup(monitor, k8sClient, desired, nil, false, name)
}

func instantBackup(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	desired *tree.Tree,
	orbconfig *orb.Orb,
	gitops bool,
	name string,
) error {
	current := &tree.Tree{}

	query, _, _, _, _, _, err := orbz.AdaptFunc(orbconfig, "instantbackup", nil, gitops, []string{"instantbackup"}, name)(monitor, desired, current)
	if err != nil {
		monitor.Error(err)
		return err
	}

	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	if err != nil {
		monitor.Error(err)
		return err
	}

	if err := ensure(k8sClient); err != nil {
		monitor.Error(err)
		return err
	}
	return nil
}
