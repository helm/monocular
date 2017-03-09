# Instance configuration

Monocular's backend can be configured by providing a config file in the following location `$HOME/monocular/config/monocular.yaml`.

You can find an example file [here](config.example.yaml).

### Change Chart repositories to pull from

Monocular, by default, uses the Stable and Incubator official chart repositories as Charts sources. This behavior can be changed by providing a repositories list as shown in [this example](config.example.yaml).

> **NOTE:** If you are deploying Monocular using the [provided Helm chart](deployment.md), refer to the `values.yaml` file to make any configuration modifications.

### Enable Helm Releases integration

Monocular adds support to interact with an existing Tiller server in order to:

* List existing Helm releases. GET `/v1/releases`
* Get information about a Helm release. GET `/v1/releases/:releaseName`
* Deploy a release. POST `/v1/releases`
* Delete a release. DELETE `/v1/releases/:releaseName`

> **IMPORTANT**: Enabling this feature will allow anybody with access to the running instance to create, list and delete any Helm release existing in your cluster.
> This feature is aimed for internal, behind the firewall deployments of Monocular, please plan accordingly.

#### Requirements

* Enable the feature by setting the `releasesEnabled` flag in the configuration file (or values.yaml if using the provided chart)
* If you are running Monocular outside of the k8s cluster, you will need to enable port forwarding
by setting the following flag `tillerPortForward: true`.
