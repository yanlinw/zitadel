package core

import (
	"github.com/caos/zitadel/operator/common"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestBackup_JobSpec1(t *testing.T) {
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{
		{Key: "testKey", Operator: "testOp"}}
	backupName := "testName"
	version := "testVersion"
	command := "test"
	sourceAKIDName := "testAKIDN"
	sourceAKIDKey := "testAKIDK"
	sourceSAKName := "testSAKN"
	sourceSAKKey := "testSAK"
	destAKIDName := "testAKIDN"
	destAKIDKey := "testAKID"
	destSAKName := "testSAKN"
	destSAKKey := "testSAK"

	equals := batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				NodeSelector:  nodeselector,
				Tolerations:   tolerations,
				Containers: []corev1.Container{{
					Name:  backupName,
					Image: common.BackupImage.Reference("", version),
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

	assert.Equal(t, equals, GetJobSpecDef(
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
		common.BackupImage.Reference("", version),
		command))
}

func TestBackup_JobSpec2(t *testing.T) {
	nodeselector := map[string]string{"test2": "test2"}
	tolerations := []corev1.Toleration{
		{Key: "testKey2", Operator: "testOp2"}}
	backupName := "testName2"
	version := "testVersion2"
	command := "test2"
	sourceAKIDName := "testAKIDN2"
	sourceAKIDKey := "testAKIDK2"
	sourceSAKName := "testSAKN2"
	sourceSAKKey := "testSAK2"
	destAKIDName := "testAKIDN2"
	destAKIDKey := "testAKID2"
	destSAKName := "testSAKN2"
	destSAKKey := "testSAK2"

	equals := batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				NodeSelector:  nodeselector,
				Tolerations:   tolerations,
				Containers: []corev1.Container{{
					Name:  backupName,
					Image: common.BackupImage.Reference("", version),
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

	assert.Equal(t, equals, GetJobSpecDef(
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
		common.BackupImage.Reference("", version),
		command,
	))
}
