# Instance configuration

Monocular's backend can be configured by providing a config file in the following location `home/monocular/config.yaml`.

You can find an example file [here](config.example.yaml).

#### Change Chart repositories to pull from

Monocular, by default, uses the Stable and Incubator official chart repositories as Charts sources. This behavior can be changed by providing a repositories list as shown in [this example](config.example.yaml).

> **NOTE:** If you are deploying Monocular using the [provided Helm chart](deployment.yaml), refer to the `values.yaml` file to make any configuration modifications.
