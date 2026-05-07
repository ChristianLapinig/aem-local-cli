# AEM Local CLI

AEM Local CLI is a cli that helps manage local AEM environments.

## Technologies

### Programming Language

- Go

### CLI Framework

- [Cobra](https://cobra.dev/)

## Installation

```sh
curl -fsSL https://raw.githubusercontent.com/ChristianLapinig/aem-local-cli/main/install.sh | sh
```

This installs the `aemlocal` binary to `/usr/local/bin`. A specific version can be installed by passing it as an argument:

```sh
curl -fsSL https://raw.githubusercontent.com/ChristianLapinig/aem-local-cli/main/install.sh | sh -s v0.1.0
```

## Updating

Run the built-in update command to upgrade to the latest version:

```sh
aemlocal update
```

If the binary is installed in a system directory (e.g. `/usr/local/bin`), you may need to run it with `sudo`:

```sh
sudo aemlocal update
```

Alternatively, re-run the install script directly:

```sh
curl -fsSL https://raw.githubusercontent.com/ChristianLapinig/aem-local-cli/main/install.sh | sh
```

To install a specific version, pass the version as an argument:

```sh
curl -fsSL https://raw.githubusercontent.com/ChristianLapinig/aem-local-cli/main/install.sh | sh -s v0.2.0
```

### Update notifications

Once per day, `aemlocal` checks for a new release in the background. If a newer version is available, a notice is printed after any command:

```
A new version of aemlocal is available: v0.3.0 (you have v0.2.0)
Run 'aemlocal update' to upgrade.
```

## Commands/Usage

### `init`

The `init` command generates a configuration directory `.aemlocal` with the following structure

```
.aemlocal
├── temp // create jobs are initially stored here 
└─- config.json
```

**Options**

- `-p, --path` - Path where to create the configuration directory.

**Usage**

```bash
$ aemlocal init

# With options
$ aemlocal init -p /Users/me/Documents
```

### `create`

The `create` command allows you to generate a local AEM environment with the necessary JARs, and [`license.properties`](http://license.properties) files in the respective author and publish folders. This command assumes that you have a valid AEM Quickstart JAR, and [`license.properties`](http://license.properties) file required to run AEM locally.

The environment is created as a named subdirectory within the base path (e.g. `-p /Users/me/envs -n cloud-service` creates the environment at `/Users/me/envs/cloud-service`). If `-n` is not provided, you will be prompted to enter a name. The command will error if an environment with the same name already exists in the config, or if the destination directory already exists on disk.

**Options**

- `-n/--name` - Name of the environment. If omitted, you will be prompted to enter one.
- `-p/--path` - Base directory where the environment subdirectory should be created. Must already exist. Defaults to the current working directory.
- `--author-port` - Specifies the port the author instance should run on (default = 4502).
- `--publish-port` - Specifies the port the publish instance should run on (default = 4503).

**Usage**

```bash
# Prompts for environment name, creates within the current working directory
$ aemlocal create /path/to/license.properties /path/to/aem-quickstart.jar

# With options — creates environment at /Users/me/envs/cloud-service
$ aemlocal create /path/to/license.properties /path/to/aem-quickstart.jar -n cloud-service -p /Users/me/envs --author-port 8080 --publish-port 8081
```

### `add`

The `add` command allows you to add an existing environment.

**Arguments**

- `<name>` - Name to assign to the environment.
- `<path-to-environment>` - Path to the existing environment directory.

**Usage**

```bash
$ aemlocal add my-env /path/to/existing/environment
```

### `delete`

The `delete` command removes a local AEM environment from the config. If `-n` is not provided, you will be prompted to select an environment from a list. You will always be asked to confirm before deletion proceeds.

**Options**

- `-n, --name` - Name of the environment to delete.
- `--purge` - Also deletes the environment directory from the filesystem. If omitted, you will be prompted after confirming deletion.

**Usage**

```bash
# Prompts to select an environment, then confirms before deleting
$ aemlocal delete

# With options
$ aemlocal delete -n my-env

# Also remove the environment folder from disk
$ aemlocal delete -n my-env --purge
```

### `update`

The `update` command checks for a newer version of `aemlocal` and, if one is available, downloads it and replaces the running binary in place.

```bash
$ aemlocal update
Checking for updates...
Updating v0.2.0 → v0.3.0
Downloading... done
Successfully updated to v0.3.0

# Already on the latest version
$ aemlocal update
Checking for updates...
Already on the latest version (v0.3.0)
```

### `list`

The `list` command list all local AEM environments.

```bash
$ aemlocal list

┌──────────────┬─────────────────────────────────────────────────────┐
│     NAME     │                        PATH                         │
├──────────────┼─────────────────────────────────────────────────────┤
│ my-env       │ /Users/me/aem/cloud-service/my-env                  │
└──────────────┴─────────────────────────────────────────────────────┘
```

## Building and Testing

### Build

Run `go build` to build the binary

```bash
$ go build
$ ./aem-local-cli

# Custom binary name
$ go build -o aemlocal
$ ./aemlocal
```

### Testing

Run `go test -v ./...`  to run all tests.

```bash
go test -v ./...
```

## Contributing

### Issues

If you encounter any bugs or would like to propose a feature enhancement, please open an issue in the “Issues” tab of the repository.

### Pull Requests

If you want to contribute code, we recommend checking the “Issues” tab first. Please follow these steps to submit a pull request:

1. For the repository and `git clone` it to your local machine.
2. Create a new branch similar to `git checkout -b name-of-branch`
    1. If you are working on an issue, we recommend using the issue number as your branch name. For example, `issue-2`.
3. Commit and push changes to your fork and the current branch you are working on.
4. Submit a pull request.
