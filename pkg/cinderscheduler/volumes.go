package cinderscheduler

import (
	"github.com/openstack-k8s-operators/cinder-operator/pkg/cinder"
	corev1 "k8s.io/api/core/v1"
	storage "lib-common/modules/storage"
)

// GetVolumes -
func GetVolumes(parentName string, name string, extraVol []storage.CinderExtraVolMounts) []corev1.Volume {
	var config0640AccessMode int32 = 0640

	schedulerVolumes := []corev1.Volume{
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

	return append(cinder.GetVolumes(parentName, extraVol, []storage.ServiceType{storage.Cinder, storage.CinderScheduler}), schedulerVolumes...)
}

// GetInitVolumeMounts - Cinder Scheduler init task VolumeMounts
func GetInitVolumeMounts(extraVol []storage.CinderExtraVolMounts) []corev1.VolumeMount {

	customConfVolumeMount := corev1.VolumeMount{
		Name:      "config-data-custom",
		MountPath: "/var/lib/config-data/custom",
		ReadOnly:  true,
	}

	return append(cinder.GetInitVolumeMounts(extraVol, []storage.ServiceType{storage.Cinder, storage.CinderScheduler}), customConfVolumeMount)
}

// GetVolumeMounts - Cinder Scheduler VolumeMounts
func GetVolumeMounts(extraVol []storage.CinderExtraVolMounts) []corev1.VolumeMount {
	return cinder.GetVolumeMounts(extraVol, []storage.ServiceType{storage.Cinder, storage.CinderScheduler})
}
