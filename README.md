# Monocular
[![Build
Status](https://travis-ci.org/kubernetes-helm/monocular.svg?branch=master)](https://travis-ci.org/kubernetes-helm/monocular)

Monocular is web-based UI for managing Kubernetes applications packaged as Helm
Charts. It allows you to search and discover available charts from multiple
repositories, and install them in your cluster with one click.

![Monocular Screenshot](docs/MonocularScreenshot.gif)

See Monocular in action at [KubeApps.com](https://kubeapps.com) or click [here](docs/about.md) to learn more about Helm, Charts and Kubernetes.

##### Video links
- [Screencast](https://www.youtube.com/watch?v=YoEbvDrI5ng)
- [Helm and Monocular Webinar](https://www.youtube.com/watch?v=u8kDkHgRbWQ)

## Install

You can use the chart in this repository to install Monocular in your cluster.

### Prerequisites
- [Helm and Tiller installed](https://github.com/kubernetes/helm/blob/master/docs/quickstart.md)
- [Nginx Ingress controller](https://kubeapps.com/charts/stable/nginx-ingress)
  - Install with Helm: `helm install stable/nginx-ingress --name monocular`
  - **Minikube/Kubeadm**: `helm install stable/nginx-ingress --set controller.hostNetwork=true`

### Installation
- The monocular pods need dynamic storage provisioning based on Storage Class enabled. Run the following to list the StorageClasses in your cluster.
`$kubectl get storageclass`
The output will be similar to this
```console
NAME                 TYPE
etcd-backup-gce-pd   kubernetes.io/gce-pd
general              ceph.com/rbd
nfs-general          example.com/nfs
```

- To mark the one you want to use for dynamic provisioning, change its value to 'true'.

`kubectl patch storageclass <your-class-name> -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'`

This marks the StorageClass as default.

Note: If you are trying to mount NFS, you might run into a mount error with your pods. You need to have the 'nfs-common' binary installed which gives the '/sbin/mount.nfs' helper program required for mounting NFS.

`sudo apt-get install nfs-common`

- Add the monocular repo to helm repo list and install the monocular using the helm charts. Use the set command to specify the StorageClass claim name for the PVC. 
```console
$ helm repo add monocular https://kubernetes-helm.github.io/monocular
$ helm install monocular/monocular --name monocular --namespace <namespace> --set volumes.name.persistentVolumeClaim.claimName=<your-class-name> 
```

### Access Monocular

Use the Ingress endpoint to access your Monocular instance:

```console
# Wait for all pods to be running (this can take a few minutes)
$ kubectl get pods --watch

$ kubectl get ingress
NAME                    HOSTS            ADDRESS         PORTS      AGE
monocular-monocular       *           192.168.64.30        80        11h
```

Visit the address specified in the Ingress object in your browser, e.g. http://192.168.64.30. New repos and charts can be added from the Web UI as required for deployment.

Read more on how to deploy Monocular [here](deployment/monocular/README.md).

## Documentation

- [Configuration](deployment/monocular/README.md#configuration)
- [Deployment](deployment/monocular/README.md)
- [Development](docs/development.md)

## Roadmap

The [Monocular roadmap is currently located in the wiki](https://github.com/kubernetes-helm/monocular/wiki/Roadmap).

## Contribute

This project is still under active development, so you'll likely encounter
[issues](https://github.com/kubernetes-helm/monocular/issues).

Interested in contributing? Check out the [documentation](CONTRIBUTING.md).

Also see [developer's guide](docs/development.md) for information on how to
build and test the code.
