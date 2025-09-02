# OASnake 🐍

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Generate a complete, standalone CLI client for any REST API from an OpenAPI specification.

OASnake reads your OpenAPI (v2 or v3) specification and
generates a Cobra-based CLI client, ready to be compiled and used.

## ✨ Features

- **CLI Generation:** Transforms any OpenAPI specification into a functional CLI application.
- **Code Generation:** Will create ready to use (and extends) go code.
- **Compilation:** Can create a binary file from go command line, or docker.
- **Customization:** Allows customization of the command name, binary name, and more.
- **Sub-command Handling:** Generates sub-commands for each API path, with flags for parameters.
- **No External Dependencies:** The generated binary is self-contained and does not require `oasnake` to run.

## 🚀 Installation

There are two ways to install `oasnake`.

### Install with Docker

If you don't have a Go environment set up, you can build a Docker image.

```bash
docker build -t oasnake:0.0.1 .
```

### Install with Go

To install `oasnake` directly, use `go install`:

```bash
go install github.com/louislouislouislouis/oasnake@latest
```

Make sure your `$GOPATH/bin` environment variable is in your `PATH`.

## 🛠️ Usage

The main command is `generate`.
It takes an OpenAPI specification file and generates the source code for a CLI client.

### Run with Docker

```bash
docker run -it --rm \
-v "$(pwd):/app" oasnake:0.0.1 generate \
--input /app/oas-spec.yaml \
--output /app/out \
--module github.com/myusername/myrepo
```

### Run with Go

```bash
oasnake generate \
--input <path/to/openapi.yaml> \
--module <your/go/module> [flags]
```

For a full list of available flags and their descriptions, please refer to the [documentation](doc/oasnake.md).

## 📝 Example: Generating a GitHub CLI

This example demonstrates how to generate a command-line interface
for the GitHub API.
The resulting CLI will be named `github-cli` and will be compiled and ready to use.

> You need Go installed on your machine to compile the generated code.

### 1. Download the OpenAPI Specification

First, download the OpenAPI specification for the GitHub API.
This example uses the specification for GitHub Enterprise Server 3.14.

```bash
curl -L -o ghes-3.14.yaml \
  https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/ghes-3.14/ghes-3.14.yaml
```

### 2. Generate and Compile the CLI

Next, use `oasnake` to generate the Go source code and compile it into a binary.
This example targets macOS on an ARM64 architecture,
but you can adjust the `--target-os` and `--target-arch` flags as needed.

The `--server-url` is set to `https://api.github.com` to ensure
the generated CLI communicates with the public GitHub API.

```bash
oasnake generate \
  --input ghes-3.14.yaml \
  --output github-cli \
  --module github.com/example/github-cli \
  --name github-cli \
  --compile-with-go \
  --target-os darwin \
  --target-arch arm64 \
  --server-url https://api.github.com
```

This command creates a new directory named `github-cli` containing
the generated Go project and the compiled binary.

### 3. Using Your New CLI

You can now use your newly generated CLI to interact with the GitHub API.

First, navigate into the output directory:

```bash
cd github-cli
```

To enable command-line completion for `zsh` (or your preferred shell), run:

```bash
source <(./github-cli completion zsh)
```

You can now explore the available commands:

```bash
./github-cli --help
```

To see the subcommands for a specific resource, like `repos`, you can also use `--help`:

```bash
./github-cli repos --help
```

#### Example: Listing Repository Information

To fetch information about a specific repository, you can use the `repos get` command.
This example retrieves details for the `oasnake` repository, owned by `louislouislouislouis`.

A GitHub Personal Access Token (PAT) is required for authentication.

```bash
./github-cli repos get \
  --owner louislouislouislouis \
  --repo oasnake \
  --tokenBearer "YOUR_GITHUB_PAT"
```

For a more readable output, you can pipe the JSON response to `jq`:

```bash
./github-cli repos get \
  --owner louislouislouislouis \
  --repo oasnake \
  --tokenBearer "YOUR_GITHUB_PAT" | jq .
```

## 🔧 How It Works

OASnake parses the provided OpenAPI specification.
Using a series of Go templates (`.gotmpl`), it generates:

- A `Cobra` command structure.
- `http.Client` calls for each API method.
- Logic to handle parameters (query, header, body).
- A `main.go` and `go.mod` to create a complete Go project. (optional)

The result is a standalone Go project in the output directory, ready to be compiled.

## 🗺️ Roadmap

- [ ] Add usage for reserved keyword
- [ ] Add default value for compilation (go).
- [ ] Better management of stdOut and stdErr for output of compilation.
- [ ] Autodetection of the host OS as the default value for `target-os` and `target-arch`.
- [ ] Docker compilation : create custom image. Search for this image. Create if not. Use after.
- [ ] Add installation option.

## 🤝 Contributing

Contributions are welcome! Feel free to open an issue or a pull request.

## 📄 License

This project is licensed under the MIT License. See the `LICENSE` file for details.
