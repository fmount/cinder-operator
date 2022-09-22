package cinder

import (
	cinderv1beta1 "github.com/openstack-k8s-operators/cinder-operator/api/v1beta1"
	common "github.com/openstack-k8s-operators/lib-common/modules/common"
	"github.com/openstack-k8s-operators/lib-common/modules/common/env"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// DBSyncCommand -
	// FIXME?: The old CN-OSP use of bootstrap.sh does not work here, but not using it might be
	// a problem as it has a few conditionals that should perhaps be considered (and they're not here)
	DBSyncCommand = "/usr/local/bin/kolla_set_configs && su -s /bin/sh -c \"cinder-manage db sync\""
)

// DbSyncJob func
func DbSyncJob(instance *cinderv1beta1.Cinder, labels map[string]string) *batchv1.Job {

	args := []string{"-c"}
	if instance.Spec.Debug.DBSync {
		args = append(args, common.DebugCommand)
	} else {
		args = append(args, DBSyncCommand)
	}

	runAsUser := int64(0)
	envVars := map[string]env.Setter{}
	envVars["KOLLA_CONFIG_FILE"] = env.SetValue(KollaConfigDbSync)
	envVars["KOLLA_CONFIG_STRATEGY"] = env.SetValue("COPY_ALWAYS")
	envVars["KOLLA_BOOTSTRAP"] = env.SetValue("TRUE")

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-db-sync",
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:      "OnFailure",
					ServiceAccountName: ServiceAccount,
					Containers: []corev1.Container{
						{
							Name: instance.Name + "-db-sync",
							Command: []string{
								"/bin/bash",
							},
							Args:  args,
							Image: instance.Spec.CinderAPI.ContainerImage,
							SecurityContext: &corev1.SecurityContext{
								RunAsUser: &runAsUser,
							},
							Env:          env.MergeEnvs([]corev1.EnvVar{}, envVars),
							VolumeMounts: GetVolumeMounts(instance.Spec.CephSecret),
						},
					},
					Volumes: GetVolumes(instance.Name, instance.Spec.CephSecret),
				},
			},
		},
	}

	job.Spec.Template.Spec.Volumes = GetVolumes(ServiceName, instance.Spec.CephSecret)

	initContainerDetails := APIDetails{
		ContainerImage:       instance.Spec.CinderAPI.ContainerImage,
		DatabaseHost:         instance.Status.DatabaseHostname,
		DatabaseUser:         instance.Spec.DatabaseUser,
		DatabaseName:         DatabaseName,
		OSPSecret:            instance.Spec.Secret,
		DBPasswordSelector:   instance.Spec.PasswordSelectors.Database,
		UserPasswordSelector: instance.Spec.PasswordSelectors.Service,
		VolumeMounts:         GetInitVolumeMounts(instance.Spec.CephSecret),
	}
	job.Spec.Template.Spec.InitContainers = InitContainer(initContainerDetails)

	return job
}
