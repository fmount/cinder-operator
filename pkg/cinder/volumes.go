package cinder

import (
	cinderv1 "github.com/openstack-k8s-operators/cinder-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

// GetVolumes -
func GetVolumes(name string, extraVol []cinderv1.CinderVolMounts) []corev1.Volume {
	var scriptsVolumeDefaultMode int32 = 0755
	var config0640AccessMode int32 = 0640

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
	if len(extraVol) > 0 {
		for _, ev := range extraVol {
			vms = append(vms, ev.Volumes...)
		}
	}

	return vms
}

// GetInitVolumeMounts - Nova Control Plane init task VolumeMounts
func GetInitVolumeMounts(extraVol []cinderv1.CinderVolMounts) []corev1.VolumeMount {
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

	if len(extraVol) > 0 {
		for _, ev := range extraVol {
			vm = append(vm, ev.Mounts...)
		}
	}
	return vm
}

// GetVolumeMounts - Nova Control Plane VolumeMounts
func GetVolumeMounts(extraVol []cinderv1.CinderVolMounts) []corev1.VolumeMount {
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

	if len(extraVol) > 0 {
		for _, ev := range extraVol {
			vm = append(vm, ev.Mounts...)
		}
	}
	return vm
}
