---
title: Author Bundles
description: Create a Bundle with Porter
aliases:
- /authoring-bundles/
---

Porter generates a bundle from its manifest, porter.yaml. The manifest is made up of a few components:

* [Bundle Metadata](#bundle-metadata)
* [Mixins](#mixins)
* [Parameters](#parameters)
* [Outputs](#outputs)
* [Parameter and Output Schema](#parameter-and-output-schema)
* [Credentials](#credentials)
* [Bundle Actions](#bundle-actions)
* [Dependencies](#dependencies)
* [Images](#images)
* [Custom](#custom)
* [Required](#required)
* [Generated Files](#generated-files)

We have full [examples](https://github.com/deislabs/porter/tree/main/examples) of Porter manifests in the Porter repository.

## Bundle Metadata

A lot of the metadata is defined by the [CNAB Spec](https://github.com/cnabio/cnab-spec/blob/master/101-bundle-json.md)
although Porter does have extra fields that are specific to making Porter bundles.

```yaml
name: azure-wordpress
description: Install Wordpress on Azure
version: 0.1.0
tag: getporter/azure-wordpress
invocationImage: getporter/azure-wordpress-installer
dockerfile: dockerfile.tmpl
```

* `name`: The name of the bundle
* `description`: A description of the bundle
* `version`: The version of the bundle, uses [semver](https://semver.org). Should not have a 'v' prefix.
* `tag`: The tag to use when the bundle is published to a registry. The format is `REGISTRY/IMAGE`, or optionally `REGISTRY/IMAGE:TAG` 
   if you do not want the default behavior of defaulting the docker image TAG portion to the version of the bundle.
* `invocationImage`: OPTIONAL. The tag to use when the invocation image is built and published to a registry. The format is
    `REGISTRY/IMAGE` or optionally `REGISTRY/IMAGE:TAG` if you do not want the default behavior of defaulting the docker image TAG 
    portion to the version of the bundle. When omitted the `invocationImage` defaults to the value of `tag` with `-installer` suffixed. 
    For example if the bundle `tag` is `getporter/porter-hello` and the `version` is `0.1.0`, then the `invocationImage` 
    will default to `getporter/porter-hello-installer:v0.1.0`
* `dockerfile`: OPTIONAL. The relative path to a Dockerfile to use as a template during `porter build`. 
    See [Custom Dockerfile](/custom-dockerfile/) for details on how to use a custom Dockerfile.
* `custom`: OPTIONAL. A map of [custom bundle metadata](https://github.com/cnabio/cnab-spec/blob/master/101-bundle-json.md#custom-extensions).

## Mixins

Mixins are adapters between the Porter and an existing tool or system. They know how to talk to Porter to include everything
they need to run, such as a CLI or config files, and how to execute their steps in the Porter manifest.

Anyone can [create a mixin](/mixin-dev-guide/), here's a list of the mixins that are installed with Porter by default:

* exec - run shell scripts and commands
* helm - use the helm cli
* azure - provision services on the Azure cloud

Declare the mixins that your bundle uses with the `mixins` section of the manifest:

```yaml
mixins:
- exec
- helm
```

Some mixins allow you to specify configuration data that is provided to the mixin during `porter build`. Each mixin
has its own format for the configuration. For example, the az mixin allows you to specify extensions to install:

```yaml
mixins:
- az:
    extensions:
    - azure-cli-iot-ext
```

See [Using Mixins](/use-mixins) to learn more about how mixins work.

## Parameters

Parameters are part of the [CNAB Spec](https://github.com/cnabio/cnab-spec/blob/master/101-bundle-json.md#parameters) and
allow you to pass in configuration values when you execute the bundle.

Learn more about [how parameters work in Porter](/parameters/).

```yaml
parameters:
- name: mysql_user
  type: string
  default: azureuser
- name: mysql_password
  type: string
  sensitive: true
  applyTo:
    - install
    - upgrade
- name: database_name
  type: string
  default: "wordpress"
  env: MYSQL_DATABASE
```

* `name`: The name of the parameter.
* `type`: The data type of the parameter: string, integer, number, boolean or file.
* `default`: (Optional) The default value for the parameter, which will be used if not supplied elsewhere.
* `env`: (Optional) The name for the destination environment variable in the bundle. Defaults to the name of the parameter in upper case, if path is not specified.
* `path`: (Optional) The destination file path in the bundle.
* `sensitive`: (Optional) Designate this parameter's value as sensitive, for masking in console output.
* `applyTo`: (Optional) Designate which actions this parameter is applicable. When not supplied, it is assumed the parameter applies to all actions.
 
## Outputs

Outputs are part of the [CNAB Spec](https://github.com/cnabio/cnab-spec/blob/master/101-bundle-json.md#outputs) to
allow access to outputs generated during the course of executing a bundle.  These are global/bundle-wide outputs,
as opposed to step outputs described in [Parameters, Credentials and Outputs](/wiring/).

```yaml
outputs:
- name: mysql_user
  type: string
  description: "MySQL user name"
- name: mysql_password
  type: string
  applyTo:
    - install
    - upgrade
- name: kubeconfig
  type: file
  path: /root/.kube/config
```

* `name`: The name of the output.
* `type`: The data type of the output: string, integer, number, boolean.
* `applyTo`: (Optional) Restrict this output to a given list of actions. If empty or missing, applies to all actions.
* `description`: (Optional) A brief description of the given output.
* `sensitive`: (Optional) Designate an output as sensitive. Defaults to false.
* `path`: (Optional) Path where the output file should be retrieved.

Outputs must either have the same name as an output from a step, meaning that the output is generated by a step, or
it must define a `path` where the output file can be located on the filesystem.

### Parameter and Output Schema

The [CNAB Spec for definitions](https://github.com/cnabio/cnab-spec/blob/master/101-bundle-json.md#definitions)
applies to both parameters and outputs.  Parameters and outputs can use [json schema 7](https://json-schema.org) 
properties to describe acceptable values. Porter uses a slightly [modified schema][json-schema] because CNAB disallows 
non-integer values.

Below are a few examples of common json schema properties and how they are used by Porter:

* `default`: The default value for the parameter. When a default is not provided, the parameter is required.
* `enum`: A list of allowed values.
* [Numeric Range](https://json-schema.org/understanding-json-schema/reference/numeric.html?highlight=minimum#range) 
    using `minimum`, `maximum`, `exclusiveMinimum` and `exclusiveMaximum`.
* String length using `minLength` and `maxLength`.

```yaml
parameters:
- name: color
  type: string
  default: blue
  enum:
  - red
  - green
  - blue
- name: size
  type: integer
  minimum: 1
  maximum: 11
- name: label
  type: string
  minLength: 3
```

[json-schema]: https://github.com/cnabio/cnab-spec/blob/master/schema/definitions.schema.json

## Credentials

Credentials are part of the [CNAB Spec](https://github.com/cnabio/cnab-spec/blob/master/802-credential-sets.md) and allow
you to pass in sensitive data when you execute the bundle, such as passwords or configuration files. 

Learn more about [how credentials work in Porter](/credentials/).

By default, all credential values are considered sensitive and will be masked in console output.

```yaml
credentials:
- name: SUBSCRIPTION_ID
  env: AZURE_SUBSCRIPTION_ID
- name: CLIENT_ID
  env: AZURE_CLIENT_ID
- name: TENANT_ID
  env: AZURE_TENANT_ID
- name: CLIENT_SECRET
  env: AZURE_CLIENT_SECRET
- name: kubeconfig
  path: /root/.kube/config
```

* `name`: The name of the credential on the local system.
* `env`: The name of the environment variable to create with the value from the credential.
* `path`: The file path to create with the file from the credential.

### See Also
* [porter credentials generate](/cli/porter_credentials_generate/)
* [How Credentials Work](/how-credentials-work/)
* [Wiring Credentials](/wiring/)

## Bundle Actions

Porter and its mixins supports the three CNAB actions: install, upgrade, and
uninstall. Within each action, you define an ordered list of steps to execute.
Each step has a mixin, a `description`, and optionally `outputs`.

```yaml
install:
- helm:
    description: "Install MySQL"
    name: mydb
    chart: stable/mysql
    version: 0.10.2
    set:
      mysqlDatabase: "{{ bundle.parameters.database-name }}"
      mysqlUser: "{{ bundle.parameters.mysql-user }}"
  outputs:
  - name: mysql-root-password
    secret: mydb-creds
    key: mysql-root-password
  - name: mysql-password
    secret: mydb-creds
    key: mysql-password
```

* `MIXIN`: The name of the mixin that will handle this step. In the example above, `helm` is the mixin.
* `description`: A description of the step, used for logging.
* `outputs`: Any outputs provided by the steps. The `name` is required but the rest of the the schema for the 
output is specific to the mixin. In the example above, the mixin will make the Kubernetes secret data available as outputs.
By default, all output values are considered sensitive and will be masked in console output.

### Custom Actions
You can also define custom actions, such as `status` or `dry-run`, and define steps for them just as you would for
the main actions (install/upgrade/uninstall). Most of the mixins support custom actions but not all do.

You have the option of declaring your custom action, though it is not required. Custom actions are defaulted 
to `stateless: false` and `modifies: true`. [Well-known actions][well-known-actions] defined in the CNAB specification 
are automatically defaulted, such as `dry-run`, `help`, `log`, and `status`. You do not need to declare custom 
actions unless you want to change the defaults.

You may want to declare your custom action when the action does not make any changes, and its execution should not 
be recorded (`stateless: true` and `modifies: false`). The `help` action is supported out-of-the-box by Porter
and is automatically defaulted to this definition so you do not need to declare it. If you have an action that is 
similar to `help`, but has a different name, you should declare it in the `customActions` section.

```
customActions:
  myhelp:
    description: "Print a special help message"
    stateless: true
    modifies: false
```

* `description`: Description of the action.
* `stateless`: Indicates that the action is purely informational, that credentials are not required, 
   and that Porter should not keep track of when this action was last executed.
* `modifies`: Indicates whether this action modifies resources managed by the bundle.

[well-known-actions]: https://github.com/cnabio/cnab-spec/blob/master/804-well-known-custom-actions.md

## Dependencies

Dependencies are an extension of the [CNAB Spec](https://github.com/cnabio/cnab-spec/blob/master/500-CNAB-dependencies.md).
See [dependencies](/dependencies/) for more details on how Porter handles dependencies.

```yaml
dependencies:
  mysql:
    tag: getporter/wordpress:v0.1.0
    parameters:
      database_name: wordpress
      mysql_user: wordpress
```

* `tag`: The tag where the bundle can be found in an OCI registry. The format should be `REGISTRY/NAME:TAG` where TAG is 
    the semantic version of the bundle.
* `parameters`: Optionally set default values for parameters in the bundle.

## Images

The `images` section of the Porter manifest corresponds to the [CNAB Spec](https://github.com/cnabio/cnab-spec/blob/master/103-bundle-runtime.md#image-maps).
These are images used in the bundle and declaring them enables Porter to manage them for you:

* publishing the bundle [copies referenced images into the published bundle](/distribute-bundles/#image-references-after-publishing).
* [archiving the bundle](/archive-bundles/) includes the referenced images in the archive.

```yaml
images:
  websvc:
      description: "A simple web service"
      imageType: "docker"
      repository: "jeremyrickard/devops-days-msp"
      digest: "sha256:85b1a9b4b60a4cf73a23517dad677e64edf467107fa7d58fce9c50e6a3e4c914"
```

This information is used to generate the corresponding section of the `bundle.json` and can be
used to in [template expressions](/wiring), much like `parameters`, `credentials` and `outputs`, allowing you to build image references using 
the `repository` and `digest` attributes. For example:

```
image: "{{bundle.images.websvc.repository}}@{{bundle.images.websvc.digest}}"
```

At runtime, these will be updated appropriately if a bundle has been [copied](/copy-bundles). Note that while `tag` is available, you should prefer the use of `digest`.

## Custom

The Custom section of a Porter manifest is intended for bundle authors to capture custom data used by the bundle.
Porter passes all custom data through to the resulting `bundle.json` as-is, without attempting to parse or otherwise
understand the data.

```yaml
custom:
  - some-custom-config:
      item: value
  - more-custom-config:
      enabled: true
      succeed: please!
```

See the [Custom Extensions](https://github.com/cnabio/cnab-spec/blob/master/101-bundle-json.md#custom-extensions)
section of the CNAB Specification for more details.

## Required

The `required` section of a Porter manifest is intended for bundle authors to declare which
[Required Extensions](https://github.com/cnabio/cnab-spec/blob/master/101-bundle-json.md#required-extensions)
known and supported by Porter are needed to run the bundle.  Hence, all extension configuration data in this section
is processed by Porter at runtime; if unsupported extension configuration exists, Porter will error out accordingly.

Currently, Porter supports the following required extensions and configuration:

### Docker

Access to the host Docker daemon is necessary to run this bundle.

When the bundle is executed, this elevated privilege must be explicitly granted to the bundle using the
[Allow Docker Host Access configuration](/configuration/#allow-docker-host-access) setting.

**Name:** 

`docker`

**Configuration:**

  * `privileged: BOOLEAN` - OPTIONAL. Whether or not the `--privileged` flag should be set when the bundle's invocation image runs. Defaults to false.

Example:

```yaml
required:
  - docker:
      privileged: true
```

## Generated Files

In addition to the porter manifest, Porter generates a few files for you to create a compliant CNAB Spec bundle.

* `.cnab/`: This directory structure is created during `porter init`. You can add files to this directory and they will
be copied into the final bundle so that you can access them at runtime. The path to this directory at runtime is `/cnab`.
* `.cnab/app/run`: This file is created during `porter init` and should not be modified.
* `Dockerfile`: This file is generated during `porter build` and cannot be modified.

## See Also

* [Using Mixins](/use-mixins/)
* [Bundle Dependencies](/dependencies/)
* [Parameters, Credentials, Outputs, and Images in Porter](/wiring/)
