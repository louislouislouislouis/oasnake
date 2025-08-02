# OASnake 🐍

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Generate a complete, standalone CLI client for any REST API from an OpenAPI specification.

OASnake reads your OpenAPI (v2 or v3) specification and generates a Cobra-based CLI client, ready to be compiled and used.

## ✨ Features

- **CLI Generation:** Transforms any OpenAPI specification into a functional CLI application.
- **Code Generation:** Will create ready to use (and extends) go code.
- **Compilation:** Can create a binary file from go command line, or docker.
- **Customization:** Allows customization of the command name, binary name, and more.
- **Sub-command Handling:** Generates sub-commands for each API path, with flags for parameters.
- **No External Dependencies:** The generated binary is self-contained and does not require `oasnake` to run.

## 🚀 Installation

There are two ways to install `oasnake`.

### With Docker

If you don't have a Go environment set up, you can build a Docker image.

```bash
docker build -t oasnake:0.0.1 .
```

### With Go

To install `oasnake` directly, use `go install`:

```bash
go install github.com/louislouislouislouis/oasnake@latest
```

Make sure your `$GOPATH/bin` environment variable is in your `PATH`.

## 🛠️ Usage

The main command is `generate`. It takes an OpenAPI specification file and generates the source code for a CLI client.

If you are using the Docker image, the command will be slightly different. This command mounts your current working directory to `/app` inside the container, so you should adjust the `--input` and `--output` paths accordingly.

**Docker:**

```bash
docker run -it --rm -v "$(pwd):/app" oasnake:0.0.1 generate --input /app/oas-spec.yaml --output /app/out --module github.com/myusername/myrepo
```

**With Go:**

```bash
oasnake generate -i <path/to/openapi.yaml> -m <your/go/module> [flags]
```

For a full list of available flags and their descriptions, please refer to the [documentation](doc/oasnake.md).

## 📝 Full Example

This example generates a CLI for a test API, names it `petstore-cli`, and installs it automatically.

1. **Generate the CLI:**

    ```bash
    oasnake generate \
      -i petstore-api.yaml \
      -m github.com/my-org/petstore-cli \
      -b petstore-cli
    ```

2. **Verify Installation:**

    The `petstore-cli` binary is now in your folder.

3. **Use Your New CLI:**

    ```bash
    petstore-cli --help
    petstore-cli list-pets
    petstore-cli create-pet --name "My new pet"
    ```

## 🔧 How It Works

OASnake parses the provided OpenAPI specification. Using a series of Go templates (`.gotmpl`), it generates:

- A `Cobra` command structure.
- `http.Client` calls for each API method.
- Logic to handle parameters (query, header, body).
- A `main.go` and `go.mod` to create a complete Go project.

The result is a standalone Go project in the output directory, ready to be compiled.

## 🗺️ Roadmap

- [ ] Add default value for compilation (go).
- [ ] Make parsing of spec an independent spec as the code generation.
- [ ] Better management of stdOut and stdErr for output of compilation.
- [ ] Autodetection of the host OS as the default value for `target-os` and `target-arch`.
- [ ] Add offline installation. Create custom image. Search for this image.
- [ ] Propose an installation path.

## 🤝 Contributing

Contributions are welcome! Feel free to open an issue or a pull request.

## 📄 License

This project is licensed under the MIT License. See the `LICENSE` file for details.
