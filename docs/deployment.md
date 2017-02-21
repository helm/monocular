# Deploying Monocular in Kubernetes

## Quickstart

Monocular is packaged as a Helm Chart that can be found in the `/deployment/monocular` directory.

> **Prerequisites**
>
> The chart is configured and tested to be used alongside an NGINX Ingress controller. Please be sure that you have a running instance in your cluster before proceeding. More information  [here](https://github.com/kubernetes/ingress/tree/master/controllers/nginx).

Install the chart:

```
helm install deployment/monocular
```

Visit [Using Helm](https://github.com/kubernetes/helm/blob/master/docs/using_helm.md) to learn more about how to use Helm.

Once deployed, your application should be available at `https://[nginx-ingress-controller-service-ip]`.

### Chart details

The chart contains 3 tiers and one ingress resource.

#### Components

  * UI: AngularJS web frontend.
  * API: Golang based backend API.
  * Prerederer: PhantomJS prerenderer for SEO purposes. More information [here](https://github.com/prerender/prerender).

#### Ingress resource

The chart includes an ingress resource that is configured to route the backend API via `[your-domain]/api` so it can be easily consumed by the UI avoiding any CORS issue or configuration.

You can configure many different settings from hosts, ingress-class to tls options using the `values.yaml` file.

# CI/CD scripts

In the `/deployment` directory you can find a set of convenience scripts helpful in CI/CD setups.

## Deploy for the first time

You can deploy a new release of the operations service
using the `deploy_install.sh` script.

### Arguments

#### Required

* `API_TAG`: Image tag to be used on the API tier, it will be appended to API_IMAGE
* `UI_TAG`: Image tag to be used on the API tier. It will be appended to UI_IMAGE

#### Optional

* `API_NAME` (default: `gcr.io/helm-ui/monocular-api`)
* `UI_NAME` (default: `gcr.io/helm-ui/monocular-ui`)
* `RELEASE_NAME` (default: `Helm's provided random name`) Helm release
  name.
* `VALUES_OVERRIDE` (default: `api.image.repository=${API_IMAGE},api.image.tag=${API_TAG},ui.image.repository=${UI_IMAGE},ui.image.tag=${UI_TAG}`) Helm values to be overridden (`helm install --set $VALUES_OVERRIDE`)
* `HELM_OPTS` (default: no value) Extra options to be passed to the helm
  command i.e --dry-run --debug

### Examples

```
# Deploy the release
API_TAG=1 UI_TAG=1 ./deploy_install.sh

# Deploy setting the release name
API_TAG=1 UI_TAG=1 RELEASE_NAME='my_release' ./deploy_install.sh

# Deploy with helm debug and dry-run
API_TAG=1 UI_TAG=1 HELM_OPTS="--debug --dry-run" ./deploy_install.sh

# Deploy overridding default values
API_TAG=1 UI_TAG=1\
 VALUES_OVERRIDE="backend.image.repository=my-image,backend.service.type=LoadBalancer"\
 ./deploy_install.sh
```

## Upgrade an existing release

You can upgrade an existing release of the operations service using the `deploy_upgrade.sh` script.

### Arguments

#### Required

* `API_TAG`: Image tag to be used on the API tier, it will be appended to API_IMAGE
* `UI_TAG`: Image tag to be used on the API tier. It will be appended to UI_IMAGE
* `RELEASE_NAME`: Helm release name.

#### Optional

* `API_NAME` (default: `gcr.io/helm-ui/monocular-api`)
* `UI_NAME` (default: `gcr.io/helm-ui/monocular-ui`)
* `SKIP_UPGRADE_CONFIRMATION` (default: false) Skip the confirmation
  prompt if an upgrade needs to be performed. Useful for unattended
  name.
* `VALUES_OVERRIDE` (default: `api.image.repository=${API_IMAGE},api.image.tag=${API_TAG},ui.image.repository=${UI_IMAGE},ui.image.tag=${UI_TAG}`) Helm values to be overridden (`helm install --set $VALUES_OVERRIDE`)
* `HELM_OPTS` (default: no value) Extra options to be passed to the helm
  command i.e --dry-run --debug

### Examples

```
# Upgrade the release in attended mode
API_TAG=1 UI_TAG=1 RELEASE_NAME='my_release' ./deploy_upgrade.sh

# Upgrade in unattended mode
API_TAG=1 UI_TAG=1 RELEASE_NAME='my-release' SKIP_UPGRADE_CONFIRMATION=true ./deploy_upgrade.sh
```
