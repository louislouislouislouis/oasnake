# OASnake 🐍

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Generate a complete, standalone CLI client for any REST API from an OpenAPI specification.

OASnake reads your OpenAPI (v2 or v3) specification and generates a Go-based CLI client, ready to be compiled and used. Say goodbye to manually scripting interactions with your APIs!

```
      _             _
     | |           | |
  ___| |__   __ _  | |__   ___  _ __
 / __| '_ \ / _` | | '_ \ / _ \| '_ \
 \__ \ | | | (_| | | | | | (_) | | | |
 |___/_| |_|\__,_| |_| |_|\___/|_| |_|
```

## ✨ Features

- **CLI Generation:** Transforms an OpenAPI specification into a functional CLI application.
- **Automatic Installation:** Can automatically compile and install the generated CLI to your `GOPATH` or via Docker.
- **Customization:** Allows customization of the command name, binary name, and more.
- **Sub-command Handling:** Generates sub-commands for each API path, with flags for parameters.
- **No External Dependencies:** The generated binary is self-contained and does not require `oasnake` to run.

## 🚀 Installation

To install `oasnake`, you can use `go install`:

```bash
go install github.com/louislouislouislouis/oasnake@latest
```

Make sure your `$GOPATH/bin` environment variable is in your `PATH`.

## 🛠️ Usage

The main command is `generate`. It takes an OpenAPI specification file and generates the source code for a CLI client.

```bash
oasnake generate -i <path/to/openapi.yaml> -m <your/go/module> [flags]
```

### Required Flags

| Flag      | Alias | Description                               | Example                               |
| :-------- | :---- | :---------------------------------------- | :------------------------------------ |
| `--input` | `-i`  | Path to the OpenAPI file (local).         | `-i ./specs/my-api.yaml`              |
| `--module`| `-m`  | Go module name for the generated code.    | `-m github.com/user/my-cli`           |

### Generation Flags

| Flag                | Alias | Description                                                                                             |
| :------------------ | :---- | :------------------------------------------------------------------------------------------------------ |
| `--output`          | `-o`  | Output directory for the generated code (defaults to `out`).                                            |
| `--name`            | `-n`  | Root command name for the CLI (defaults to the API title in the spec).                                  |
| `--server-url`      |       | Server URL to use (defaults to the first server URL in the spec).                                       |
| `--with-model`      |       | Generates data models in a separate file.                                                               |

### Binary Installation Flags

| Flag                    | Alias | Description                                                                                             |
| :---------------------- | :---- | :------------------------------------------------------------------------------------------------------ |
| `--install`             |       | Enables binary installation after generation.                                                           |
| `--binary`              | `-b`  | Name of the binary to install (defaults to the command name). Implies `--install`.                      |
| `--install-with-go`     |       | Use `go install` for installation (default).                                                            |
| `--install-with-docker` |       | Use Docker to compile and install the binary.                                                           |

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

    The `petstore-cli` binary is now in your `$GOPATH/bin`.

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

- [ ] Implementation of the compilation with `go`.
- [ ] Autodetection of the host OS as the default value for `target-os` and `target-arch`.
- [ ] Add online installation. Create custom image. Search for this image.
- [ ] Propose an installation path.

## 🤝 Contributing

Contributions are welcome! Feel free to open an issue or a pull request.

## 📄 License

This project is licensed under the MIT License. See the `LICENSE` file for details.
