package core

import (
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator/helpers"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func GetCronJob(
	namespace string,
	nameLabels *labels.Name,
	cron string,
	jobSpecDef batchv1.JobSpec,
) *v1beta1.CronJob {
	return &v1beta1.CronJob{
		ObjectMeta: v1.ObjectMeta{
			Name:      nameLabels.Name(),
			Namespace: namespace,
			Labels:    labels.MustK8sMap(nameLabels),
		},
		Spec: v1beta1.CronJobSpec{
			Schedule:          cron,
			ConcurrencyPolicy: v1beta1.ForbidConcurrent,
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: jobSpecDef,
			},
		},
	}
}

func GetJob(
	namespace string,
	nameLabels *labels.Name,
	jobSpecDef batchv1.JobSpec,
) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: v1.ObjectMeta{
			Name:      nameLabels.Name(),
			Namespace: namespace,
			Labels:    labels.MustK8sMap(nameLabels),
		},
		Spec: jobSpecDef,
	}
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
	return batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				NodeSelector:  nodeselector,
				Tolerations:   tolerations,
				Containers: []corev1.Container{{
					Name:  backupName,
					Image: image,
					Command: []string{
						"/bin/bash",
						"-c",
						command,
					},
					VolumeMounts: []corev1.VolumeMount{{
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
					}},
					ImagePullPolicy: corev1.PullAlways,
				}},
				Volumes: []corev1.Volume{{
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
				}},
			},
		},
	}
}
