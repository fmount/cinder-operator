package cindervolume

import (
	"github.com/openstack-k8s-operators/lib-common/modules/storage"
	"strings"

	"github.com/openstack-k8s-operators/cinder-operator/pkg/cinder"
	corev1 "k8s.io/api/core/v1"
)

// GetVolumes -
func GetVolumes(parentName string, name string, extraVol []storage.CinderExtraVolMounts) []corev1.Volume {
	var config0640AccessMode int32 = 0640
	var dirOrCreate = corev1.HostPathDirectoryOrCreate

	volumeVolumes := []corev1.Volume{
		{
			Name: "etc-iscsi",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/etc/iscsi",
				},
			},
		},
		{
			Name: "dev",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/dev",
				},
			},
		},
		{
			Name: "lib-modules",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/lib/modules",
				},
			},
		},
		{
			Name: "run",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/run",
				},
			},
		},
		{
			Name: "sys",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/sys",
				},
			},
		},
		{
			Name: "var-lib-cinder",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/lib/cinder",
					Type: &dirOrCreate,
				},
			},
		},
		{
			Name: "var-lib-iscsi",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/lib/iscsi",
					Type: &dirOrCreate,
				},
			},
		},
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

	// Set the propagation levels for CinderVolume, including the backend name
	propagation := []storage.ServiceType{storage.Cinder, storage.CinderVolume, storage.ServiceType(strings.TrimPrefix(name, "cinder-volume-"))}
	return append(cinder.GetVolumes(parentName, extraVol, propagation), volumeVolumes...)
}

// GetInitVolumeMounts - Cinder Volume init task VolumeMounts
func GetInitVolumeMounts(name string, extraVol []storage.CinderExtraVolMounts) []corev1.VolumeMount {

	customConfVolumeMount := corev1.VolumeMount{
		Name:      "config-data-custom",
		MountPath: "/var/lib/config-data/custom",
		ReadOnly:  true,
	}

	// Set the propagation levels for CinderVolume, including the backend name
	propagation := []storage.ServiceType{storage.Cinder, storage.CinderVolume, storage.ServiceType(strings.TrimPrefix(name, "cinder-volume-"))}
	return append(cinder.GetInitVolumeMounts(extraVol, propagation), customConfVolumeMount)
}

// GetVolumeMounts - Cinder Volume VolumeMounts
func GetVolumeMounts(name string, extraVol []storage.CinderExtraVolMounts) []corev1.VolumeMount {
	volumeVolumeMounts := []corev1.VolumeMount{
		{
			Name:      "etc-iscsi",
			MountPath: "/etc/iscsi",
			ReadOnly:  true,
		},
		{
			Name:      "dev",
			MountPath: "/dev",
		},
		{
			Name:      "lib-modules",
			MountPath: "/lib/modules",
			ReadOnly:  true,
		},
		{
			Name:      "run",
			MountPath: "/run",
		},
		{
			Name:      "sys",
			MountPath: "/sys",
		},
		{
			Name:      "var-lib-cinder",
			MountPath: "/var/lib/cinder",
		},
		{
			Name:      "var-lib-iscsi",
			MountPath: "/var/lib/iscsi",
		},
	}
	// Set the propagation levels for CinderVolume, including the backend name
	propagation := []storage.ServiceType{storage.Cinder, storage.CinderVolume, storage.ServiceType(strings.TrimPrefix(name, "cinder-volume-"))}
	return append(cinder.GetVolumeMounts(extraVol, propagation), volumeVolumeMounts...)
}
