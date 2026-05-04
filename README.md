# AEM Local CLI

AEM Local CLI is a cli that helps manage local AEM environments.

## Technologies

### Programming Language

- Go

### CLI Framework

- [Cobra](https://cobra.dev/)

## Commands/Usage

### `init`

The `init` command generates a configuration directory `.aemlocal` with the following structure

```
.aemlocal
├── temp // create jobs are initially stored here 
└─- config.json
```

**Options**

- `-e, --envsPath` - Path where local AEM environments are stored.
- `-p, --path` - Path where to create the configuration directory.

**Usage**

```bash
$ aemlocal init

# With options
$ aemlocal init -e /Users/me/aem -p /Users/me/Documents
```

### `create`

The `create` command allows you to generate a local AEM environment with the necessary JARs, and [`license.properties`](http://license.properties) files in the respective author and publish folders.This command assumes that you have a valid AEM Quickstart JAR, and [`license.properties`](http://license.properties) file required to run AEM locally.

**Options**

- `-n/--name` - Name of the environment (default = `aem`)
- `-p/--path` - Relative path under `envsPath` property in `.aemlocal/config.json` where the environment should be created.
- `--author-port` - Specifies the port the author instance should run on (default = 4502).
- `--publish-port` - Specifies the port the publish instance should run on. (default = 4503).

**Usage**

```bash
$ aemlocal create /path/to/license.properties /path/to/aem-quickstart.jar

# With options
# In this example, if .aemlocal/config.json/envsPath is /Users/me/envs, the
# generated environment will be at /Users/me/envs/cloud-service/test.
$ aemlocal create /path/to/license.properties /path/to/aem-quickstart.jar -n test -p cloud-service --author-port 8080 --publish-port 8081 

```

### `add`

The `add` command allows you to add an existing environment.

```bash
aemlocal /path/to/existing/environment

```

### `delete` - COMING SOON

The delete command allows you to delete a local AEM environment.

```bash
# Default command will list environments where you can select which environment to delete
$ aemlocal delete

# With options
$ aemlocal delete -n name-of-env
```

**Options**

- `-n, --name` - Name of the environment to delete

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
