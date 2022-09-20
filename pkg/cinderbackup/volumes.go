package cinderbackup

import (
	"github.com/openstack-k8s-operators/cinder-operator/pkg/cinder"
	corev1 "k8s.io/api/core/v1"
)

// GetVolumes -
func GetVolumes(parentName string, name string, cephsecret []string) []corev1.Volume {
	var config0640AccessMode int32 = 0640

	backupVolumes := []corev1.Volume{
		{
			Name: "config-data-custom",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: &config0640AccessMode,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: name + "-config-data",
					},
				},
			},
		},
	}

	// Get all the (ceph) secrets passed as argument
	var p []corev1.VolumeProjection

	for _, v := range cephsecret {
		curr := corev1.VolumeProjection{
			Secret: &corev1.SecretProjection{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: v,
				},
			},
		}
		p = append(p, curr)
	}

	if len(cephsecret) > 0 {
		curr := corev1.Volume{
			Name: "ceph-client-conf",
			VolumeSource: corev1.VolumeSource{
				Projected: &corev1.ProjectedVolumeSource{
					Sources: p},
			},
		}
		backupVolumes = append(backupVolumes, curr)
	}

	return append(cinder.GetVolumes(parentName, cephsecret), backupVolumes...)
}

// GetInitVolumeMounts - Cinder Backup init task VolumeMounts
func GetInitVolumeMounts(cephsecret []string) []corev1.VolumeMount {

	customConfVolumeMount := corev1.VolumeMount{
		Name:      "config-data-custom",
		MountPath: "/var/lib/config-data/custom",
		ReadOnly:  true,
	}
	return append(cinder.GetInitVolumeMounts(cephsecret), customConfVolumeMount)
}

// GetVolumeMounts - Cinder Backup VolumeMounts
func GetVolumeMounts(cephsecret []string) []corev1.VolumeMount {
	return cinder.GetVolumeMounts(cephsecret)
}
