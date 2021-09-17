package core

import (
	"github.com/caos/zitadel/operator/helpers"
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/core"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	defaultMode                  int32 = 256
	certsInternalSecretName            = "client-certs"
	rootSecretName                     = "cockroachdb.client.root"
	destSakInternalSecretName          = "dest-sak"
	destAkidInternalSecretName         = "dest-akid"
	sourceSakInternalSecretName        = "source-sak"
	sourceAkidInternalSecretName       = "source-akid"
	CertPath                           = "/cockroach/cockroach-certs"
	SourceAkidSecretPath               = "/secrets/sakid"
	SourceSakSecretPath                = "/secrets/ssak"
	DestAkidSecretPath                 = "/secrets/dakid"
	DestSakSecretPath                  = "/secrets/dsak"
)

func getVolumeMounts(
	sourceAKIDKey string,
	sourceSAKKey string,
	destAKIDKey string,
	destSAKKey string,
) []corev1.VolumeMount {
	return []corev1.VolumeMount{{
		Name:      certsInternalSecretName,
		MountPath: CertPath,
	}, {
		Name:      sourceAkidInternalSecretName,
		SubPath:   sourceAKIDKey,
		MountPath: SourceAkidSecretPath,
	}, {
		Name:      sourceSakInternalSecretName,
		SubPath:   sourceSAKKey,
		MountPath: SourceSakSecretPath,
	}, {
		Name:      destAkidInternalSecretName,
		SubPath:   destAKIDKey,
		MountPath: DestAkidSecretPath,
	}, {
		Name:      destSakInternalSecretName,
		SubPath:   destSAKKey,
		MountPath: DestSakSecretPath,
	}}
}

func getVolumes(
	sourceAKIDName string,
	sourceSAKName string,
	destAKIDName string,
	destSAKName string,
) []corev1.Volume {
	return []corev1.Volume{{
		Name: certsInternalSecretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  rootSecretName,
				DefaultMode: helpers.PointerInt32(defaultMode),
			},
		},
	}, {
		Name: sourceAkidInternalSecretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  sourceAKIDName,
				DefaultMode: helpers.PointerInt32(defaultMode),
			},
		},
	}, {
		Name: sourceSakInternalSecretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  sourceSAKName,
				DefaultMode: helpers.PointerInt32(defaultMode),
			},
		},
	}, {
		Name: destAkidInternalSecretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  destAKIDName,
				DefaultMode: helpers.PointerInt32(defaultMode),
			},
		},
	}, {
		Name: destSakInternalSecretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  destSAKName,
				DefaultMode: helpers.PointerInt32(defaultMode),
			},
		},
	}}
}

func GetJobSpecDef(
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	sourceAKIDName string,
	sourceAKIDKey string,
	sourceSAKName string,
	sourceSAKKey string,
	destAKIDName string,
	destAKIDKey string,
	destSAKName string,
	destSAKKey string,
	backupName string,
	image string,
	command string,
) batchv1.JobSpec {
	return core.GetJobSpecDef(
		nodeselector,
		tolerations,
		getVolumeMounts(
			sourceAKIDKey,
			sourceSAKKey,
			destAKIDKey,
			destSAKKey,
		),
		getVolumes(
			sourceAKIDName,
			sourceSAKName,
			destAKIDName,
			destSAKName,
		),
		backupName,
		image,
		command,
	)

}
