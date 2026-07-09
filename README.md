# dis-json-template-linter

This is a custom linter written in Go for parsing Go template files that are constructing JSON snippets.

There is not a great product out there for this need and meets our needs for writing ElasticSearch queries.

## Dependencies

Requires Go version as specified in the go mod to install presently, before we introduce a distributable binary.

### Tooling

### Audit

We use `dis-vulncheck` to do auditing, which you will [need to install](https://github.com/ONSdigital/dis-vulncheck).

### Linting

We use v2 of golangci-lint, which you will [need to install](https://golangci-lint.run/docs/welcome/install).

We use [check-jsonschema](https://github.com/python-jsonschema/check-jsonschema#installing-and-running-as-a-cli-tool) for the config file schema validation.

## Installation

```sh
  go install github.com/ONSdigital/dis-json-template-linter/cmd/dis-json-template-linter@latest
```

We may provide a brew tap in future.

## Configuration

dis-json-template-linter looks for a configuration file, named '.dis-json-template-linter.yaml'.

You can view the [specification for the configuration file](config-spec.yaml) in this repository.

You can also view [the configuration file](.dis-json-template-linter.yaml) in this repository as an example.

## Running

To run dis-json-template-linter you have to provide a glob of files to look at:

```sh
  dis-json-template-linter ./**/*.tmpl
```

To see this in action, you can test it against the testdata directory:

```sh
  make install
  dis-json-template-linter ./testdata/**/*.tmpl
```

### Flags

| Flag       | Short | Default             | Description                                                                                                                                            |
|------------|-------|---------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--config` | `-c`  | _(auto-discovered)_ | Path to a `.dis-json-template-linter.yaml` config file. When omitted, the linter walks up the directory tree from each linted file until it finds one. |

Example - point to an explicit config file:

```sh
dis-json-template-linter --config /path/to/.dis-json-template-linter.yaml ./**/*.tmpl
```

```sh
dis-json-template-linter -c .dis-json-template-linter.yaml ./**/*.tmpl
```
