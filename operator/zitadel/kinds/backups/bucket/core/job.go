package core

import (
	"github.com/caos/zitadel/operator/helpers"
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/core"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	defaultMode             int32 = 256
	CertPath                      = "/cockroach/cockroach-certs"
	saInternalSecretName          = "sa-json"
	SaSecretPath                  = "/secrets/sa.json"
	akidInternalSecretName        = "akid"
	AkidSecretPath                = "/secrets/akid"
	sakInternalSecretName         = "sak"
	SakSecretPath                 = "/secrets/sak"
	certsInternalSecretName       = "client-certs"
	rootSecretName                = "cockroachdb.client.root"
)

func getVolumeMounts(
	saSecretKey string,
	assetAKIDKey string,
	assetSAKKey string,
) []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      certsInternalSecretName,
			MountPath: CertPath,
		}, {
			Name:      saInternalSecretName,
			SubPath:   saSecretKey,
			MountPath: SaSecretPath,
		}, {
			Name:      akidInternalSecretName,
			SubPath:   assetAKIDKey,
			MountPath: AkidSecretPath,
		}, {
			Name:      sakInternalSecretName,
			SubPath:   assetSAKKey,
			MountPath: SakSecretPath,
		},
	}
}

func getVolumes(
	backupSecretName string,
) []corev1.Volume {
	return []corev1.Volume{
		{
			Name: certsInternalSecretName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  rootSecretName,
					DefaultMode: helpers.PointerInt32(defaultMode),
				},
			},
		}, {
			Name: saInternalSecretName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  backupSecretName,
					DefaultMode: helpers.PointerInt32(defaultMode),
				},
			},
		}, {
			Name: akidInternalSecretName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  backupSecretName,
					DefaultMode: helpers.PointerInt32(defaultMode),
				},
			},
		}, {
			Name: sakInternalSecretName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  backupSecretName,
					DefaultMode: helpers.PointerInt32(defaultMode),
				},
			},
		},
	}
}

func GetJobSpecDef(
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	backupSecretName string,
	saSecretKey string,
	assetAKIDKey string,
	assetSAKKey string,
	backupName string,
	command string,
	image string,
) batchv1.JobSpec {
	return core.GetJobSpecDef(
		nodeselector,
		tolerations,
		getVolumeMounts(
			saSecretKey,
			assetAKIDKey,
			assetSAKKey,
		),
		getVolumes(
			backupSecretName,
		),
		backupName,
		image,
		command,
	)
}
