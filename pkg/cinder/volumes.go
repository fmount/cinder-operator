package cinder

import (
	corev1 "k8s.io/api/core/v1"
	"lib-common/modules/storage"
)

// GetVolumes -
func GetVolumes(name string, extraVol []storage.CinderExtraVolMounts, svc storage.ServiceType) []corev1.Volume {
	var scriptsVolumeDefaultMode int32 = 0755
	var config0640AccessMode int32 = 0640
	c := storage.CinderExtraVolMounts{}

	vms := []corev1.Volume{
		{
			Name: "etc-machine-id",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/etc/machine-id",
				},
			},
		},
		{
			Name: "etc-localtime",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/etc/localtime",
				},
			},
		},
		{
			Name: "scripts",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: &scriptsVolumeDefaultMode,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: name + "-scripts",
					},
				},
			},
		},
		{
			Name: "config-data",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: &config0640AccessMode,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: name + "-config-data",
					},
				},
			},
		},
		{
			Name: "config-data-merged",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{Medium: ""},
			},
		},
	}

	return c.AppendVolume(vms, extraVol, svc)
}

// GetInitVolumeMounts - Nova Control Plane init task VolumeMounts
func GetInitVolumeMounts(extraVol []storage.CinderExtraVolMounts, svc storage.ServiceType) []corev1.VolumeMount {
	c := storage.CinderExtraVolMounts{}
	vm := []corev1.VolumeMount{
		{
			Name:      "scripts",
			MountPath: "/usr/local/bin/container-scripts",
			ReadOnly:  true,
		},
		{
			Name:      "config-data",
			MountPath: "/var/lib/config-data/default",
			ReadOnly:  true,
		},
		{
			Name:      "config-data-merged",
			MountPath: "/var/lib/config-data/merged",
			ReadOnly:  false,
		},
	}

	return c.AppendVolumeMount(vm, extraVol, svc)
}

// GetVolumeMounts - Nova Control Plane VolumeMounts
func GetVolumeMounts(extraVol []storage.CinderExtraVolMounts, svc storage.ServiceType) []corev1.VolumeMount {
	c := storage.CinderExtraVolMounts{}
	vm := []corev1.VolumeMount{
		{
			Name:      "etc-machine-id",
			MountPath: "/etc/machine-id",
			ReadOnly:  true,
		},
		{
			Name:      "etc-localtime",
			MountPath: "/etc/localtime",
			ReadOnly:  true,
		},
		{
			Name:      "scripts",
			MountPath: "/usr/local/bin/container-scripts",
			ReadOnly:  true,
		},
		{
			Name:      "config-data-merged",
			MountPath: "/var/lib/config-data/merged",
			ReadOnly:  false,
		},
	}
	return c.AppendVolumeMount(vm, extraVol, svc)
}
