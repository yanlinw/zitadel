package core

import (
	"github.com/caos/orbos/pkg/labels"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	volumeMounts []corev1.VolumeMount,
	volumes []corev1.Volume,
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
					VolumeMounts: volumeMounts,

					ImagePullPolicy: corev1.PullAlways,
				}},
				Volumes: volumes,
			},
		},
	}
}
