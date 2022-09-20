# cinder-operator
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started
You’ll need a Kubernetes cluster to run against.  Our recommendation for the time being is to use
[OpenShift Local](https://access.redhat.com/documentation/en-us/red_hat_openshift_local/2.2/html/getting_started_guide/installation_gsg) (formerly known as CRC / Code Ready Containers).  
We have [companion development tools](https://github.com/openstack-k8s-operators/install_yamls/blob/master/devsetup/README.md) available that will install OpenShift Local for you.

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:
	
```sh
make docker-build docker-push IMG=<some-registry>/cinder-operator:tag
```
	
3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/cinder-operator:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

### Configure Cinder with Ceph backend

The Cinder spec API can be used to configure and customize the Ceph backend. In
particular, the `customServiceConfig` parameter should be used, for each
defined volume, to override the `enabled_backends` parameter, which must exist
in `cinder.conf` to make the `cinderVolume` pod run. The global `cephBackend`
parameter is used to specify the Ceph client-related "key/value" pairs required
to connect the service with an external Ceph cluster. Multiple external Ceph
clusters are not supported at the moment. The following represents an example
of the Cinder object that can be used to trigger the Cinder service deployment,
and enable the Cinder backend that points to an external Ceph cluster.

```
apiVersion: cinder.openstack.org/v1beta1
kind: Cinder
metadata:
  name: cinder
  namespace: openstack
spec:
  serviceUser: cinder
  databaseInstance: openstack
  databaseUser: cinder
  cinderAPI:
    replicas: 1
    containerImage: quay.io/tripleowallabycentos9/openstack-cinder-api:current-tripleo
  cinderScheduler:
    replicas: 1
    containerImage: quay.io/tripleowallabycentos9/openstack-cinder-scheduler:current-tripleo
  cinderBackup:
    replicas: 1
    containerImage: quay.io/tripleowallabycentos9/openstack-cinder-backup:current-tripleo
  secret: cinder-secret
  cinderVolumes:
    volume1:
      containerImage: quay.io/tripleowallabycentos9/openstack-cinder-volume:current-tripleo
      replicas: 1
      customServiceConfig: |
        [DEFAULT]
        enabled_backends=ceph
  cephBackend:
    cephFsid: <CephClusterFSID>
    cephMons: <CephMons>
    cephClientKey: <cephClientKey>
    cephUser: openstack
    cephPools:
      cinder:
        name: volumes
      nova:
        name: vms
      glance:
        name: images
      cinder_backup:
        name: backup
      extra_pool1:
        name: ceph_ssd_tier
      extra_pool2:
        name: ceph_nvme_tier
      extra_pool3:
        name: ceph_hdd_tier
```

When the service is up and running, it's possible to interact with the cinder
API and create the Ceph `cinder type` backend which is associated with the Ceph
tier specified in the config file.


# Example2: pass Ceph credentials as secrets

When Cinder is configured with a Ceph backend, it requires both ceph.conf and
the associated keyring.
Those files can be generated by an external entity and built as secrets in the
`openstack` namespace.

The following represents an example of secret that can be created:

```
apiVersion: v1
kind: Secret
metadata:
  name: ceph-client-conf
  namespace: openstack
stringData:
  cluster1.client.<user>.keyring: |
    [client.<user>]
        key = <client_key>
        caps mgr = "allow *"
        caps mon = "profile rbd"
        caps osd = "profile rbd pool=images"
  cluster1.conf: |
    [global]
    fsid = <fsid>
    mon_host = <list_of_mon_ip>
```

Create the secret:

`oc create -f <secret_file>.yaml`


The human operator can create several secrets according to the amount of Ceph
clusters available. The `CephSecret` key, that can be passed to the Cinder
spec, is supposed to collect the secret names associated to each cluster, then
the cinder-operator will be able to process the associated data and make them
available in /etc/ceph, which is the default location where the Ceph
credentials are stored. The following represents an example of Cinder resource
that can be used to trigger the service deployment, and enable the rbd backend
that points to multiple external Ceph clusters.

```
apiVersion: cinder.openstack.org/v1beta1
kind: Cinder
metadata:
  name: cinder
  namespace: openstack
spec:
  serviceUser: cinder
  databaseInstance: openstack
  databaseUser: cinder
  cinderAPI:
    replicas: 1
    containerImage: quay.io/tripleowallabycentos9/openstack-cinder-api:current-tripleo
  cinderScheduler:
    replicas: 1
    containerImage: quay.io/tripleowallabycentos9/openstack-cinder-scheduler:current-tripleo
  cinderBackup:
    replicas: 1
    containerImage: quay.io/tripleowallabycentos9/openstack-cinder-backup:current-tripleo
  secret: osp-secret
  cinderVolumes:
    volume1:
      containerImage: quay.io/tripleowallabycentos9/openstack-cinder-volume:current-tripleo
      replicas: 1
      customServiceConfig: |
        [DEFAULT]
        enabled_backends=ceph
        [ceph]
        volume_backend_name=ceph
        volume_driver=cinder.volume.drivers.rbd.RBDDriver
        rbd_ceph_conf=/etc/ceph/ceph.conf
        rbd_user=openstack
        rbd_pool=volumes
        rbd_flatten_volume_from_snapshot=False
        rbd_secret_uuid=4b5c8c0a-ff60-454b-a1b4-9747aa737d19
  cephSecret:
    - cluster1_secret
    - cluster2_secret
```

When the service is up and running, it's possible to interact with the exposed
API and create a volume using the configured Ceph backends.

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) 
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster 

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

