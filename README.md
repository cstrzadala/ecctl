# Elastic Cloud Control

`ecctl` is the CLI on steroids for the Elastic Cloud API. It wraps common operations commonly needed by Elastic Cloud operators within a single command line tool.

The goal of this project is to provide an easier way to interact with the Elastic Cloud API, ensuring each one of the provided commands is thoroughly tested.

## Installation

The latest stable binaries can always be found on the [release page](https://github.com/elastic/ecctl/releases) or compiled from the latest on the master branch to leverage the most recently merged features. To find out more information about building from source view the steps from our [contributing page](./CONTRIBUTING.md).

### macOS

The simplest installation for macOS users is to install `ecctl` with `brew`.

```console
$ brew tap elastic/ecctl
$ brew tap-pin elastic/ecctl
$ brew install elastic/ecctl/ecctl
Updating Homebrew...
==> Installing ecctl from elastic/ecctl
...
==> Summary
🍺  /usr/local/Cellar/ecctl/1.0.0: 6 files, 13MB, built in 7 seconds
```

#### Note: to get the autocompletions working they must be sourced from your shell profile

## Linux based OS

<!-- TODO: Update this with the approparate repo when applicable -->

## Configuration Settings and Precedence

In order for `ecctl` to be able to communicate with the API, it needs to have a set of configuration parameters defined.  
These parameters can be set in a configuration file, through environment variables, or at runtime using the CLI's global flags.

The hierarchy is as follows listed from higher precedence to lower:

1. Command line flags `--region`, `--host`, `--user`, `--pass`.
2. Environment variables.
3. Shared configuration file `$HOME/.ecctl/config.<json|toml|yaml|hcl>`.

## Generating a configuration file

If it's your first time using `ecctl`, you might want to use the `init` command to assist you
in generating a configuration file. The resulting configuration file will be saved under `~/.ecctl/config.json`:

```console
$ ecctl init
Welcome to the Elastic Cloud CLI! This command will guide you through authenticating and setting some default values.

Missing configuration file, woud you like to initialise it? [y/n]: y
Enter the URL of the Elastic Cloud API or your ECE installation: https://api.elastic-cloud.com

What default output format would you like?
  [1] text - Human-readable output format, commands with no output templates defined will fall back to JSON.
  [2] json - JSON formatted output API responses.

Please enter a choice: 1


Which authentication mechanism would you like to use?
  [1] API Keys (Recommended).
  [2] Username and Password login.

Please enter your choice: 1

Paste your API Key and press enter: xxxxx

Your credentials seem to be valid, and show you're authenticated as "user".

You're all set! Here are some commands to try:
  $ ecctl auth user key list
  $ ecctl deployment elasticsearch list
```

## Elastic Cloud authentication

Elastic Cloud uses API keys to authenticate users against its API. Additionally, it supports
the usage of [JWT](https://jwt.io/) to validate authenticated clients. The preferred authentication method is API keys.

There's two ways to authenticate against the Elastic Cloud API using `ecctl`:

* **Specifying an API key using the `--apikey` flag**.
* Specifying a `--user` and `--pass` flags.

The first method requires the user to already have an API key, if this is the case, all the outgoing API requests will use an Authentication API key header.

The second method uses the `user` and `pass` values to obtain a valid JWT token, that token is then used as the Authentication Bearer header for every API call. A goroutine that refreshes the token every minute is started, so that the token doesn't expire while we're performing actions.

### Shared configuration file

Below is an example `YAML` configuration file `$HOME/.ecctl/config.yaml` that will effectively point and configure the binary for Elastic Cloud:

```yaml
host: https://api.elastic-cloud.com # URL of Elastic Cloud or your ECE installation.
region: us-east-1 # For ECE installations you can just omit this setting.

# Credentials
## apikey is the preferred authentication mechanism.
apikey: bWFyYzo4ZTJmNmZkNjY5ZmQ0MDBkOTQ3ZjI3MTg3ZWI5MWZhYjpOQktHY05jclE0cTBzcUlnTXg3QTd3

## username and password can be used, specially when no API key is available.
user: username
pass: password
```

### Environment variables

The same settings can be defined as environment variables instead of a configuration file or to override certain settings of the `YAML` file.  If setting environment variables, you'll need to prefix the configuration parameter with `EC_` and capitalize the setting, i.e. `EC_HOST` or `EC_USER`.

```sh
export EC_APIKEY=bWFyYzo4ZTJmNmZkNjY5ZmQ0MDBkOTQ3ZjI3MTg3ZWI5MWZhYjpOQktHY05jclE0cTBzcUlnTXg3QTd3
```

#### Special Environment Variables

```sh
export EC_CONFIG=$HOME/.ecctl/cloud.yaml
```

`EC_REGIONS` will define autocomplete entries for the `--region` flag, though there is no validation and it is not environment aware (e.g. no different list is provided for staging vs production).

```sh
export EC_REGIONS=="ap-northeast-1 ap-southeast-1 ap-southeast-2 aws-eu-central-1 eu-west-1 gcp-europe-west1 gcp-europe-west3 gcp-us-central1 gcp-us-west1 sa-east-1 us-east-1 us-west-1 us-west-2"
```

## Multiple configuration support

`ecctl` supports having multiple configuration files out of the box.  This allows for easy management of multiple environments or specialized targets.  By default it will use `$HOME/.ecctl/config.<json|toml|yaml|hcl>`, but when the `--config` flag is specified, it will append the `--config` name to the file:

```console
# Default behaviour
$ ecctl version
# will use ~/.ecctl/staging.yaml

# When an environment is specified, the configuration file used will change
$ ecctl version --config ece
# will use ~/.ecctl/ece.yaml
```

## Output format

The `--output` flag allows for the response to be presented in a particular way (see `ecctl help` for an updated list of allowed formats).  The default formatter behavior is to fallback to `json` when there's no _text_ format template or if the formatting fails.

## Custom formatting

`ecctl` supports a global `--format` flag which can be passed to any existing command or subcommand.
Using the `--format` flag allows you to obtain a specific part of a command response that might not have been shown before with the default `--output=text`. The `--format` internally uses Go templates which means that you can use the power of the Go built-in [`text/templates`](https://golang.org/pkg/text/template/) on demmand.

### Examples

Obtaining the ID, Version and health status

```console
$ ecctl elasticsearch list --format '{{.ClusterID}} {{.PlanInfo.Current.Plan.Elasticsearch.Version}} {{.Healthy}}'
a2c4f423c1014941b75a48292264dd25 6.7.0 true
a4f29ff3ba554e69a1e1b40c3ee1b6e3 6.7.0 true
5e29960763ef496ea8cf6a5371328a6a 6.7.0 true
53023f28d68b4b329d9d913f110709d2 6.7.0 true
```

Since the template is executed we can also embed logic inside of the template to filter the results.

```console
$ export EC_FORMAT='{{range .Elasticsearch.DefaultPlugins}}{{if eq . "discovery-file" }}{{$.Version}}{{end}}{{end}}'
# Since the template is executed on every item of the list, filter the empty lines to have a cleaner output.
$ ecctl stack list --format "${EC_FORMAT}" | sed '/^\s*$/d'


6.2.3
$ unset EC_FORMAT
```

## Contributing

See the contribution guidelines specified in the [CONTRIBUTING](./CONTRIBUTING.md) doc.

## Release Process

See the [release guide](./developer_docs/RELEASE.md).

## Full command help

* [ecctl command help](./docs/ecctl.md)
