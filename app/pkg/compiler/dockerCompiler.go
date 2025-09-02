package compiler

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
)

type DockerCompiler struct {
	Config *CompilerConfig
}

func (c *DockerCompiler) Compile() error {
	log.Debug().Msg("Compiling binary...")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("Cannot get current working directory")
		return err
	}
	imageValue := "golang:1.23.0"
	projectPath := filepath.Join(cwd, c.Config.OutputDirectory)
	containerWorkdir := "/go/src/app"

	// Pull l'image si elle n'existe pas
	reader, err := cli.ImagePull(ctx, imageValue, image.PullOptions{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to pull Docker image")
		return err
	}
	io.Copy(os.Stderr, reader)
	script := `
go mod tidy && \
GOOS=%s GOARCH=%s go build -o %s
`
	script = fmt.Sprintf(script, c.Config.TargetOs, c.Config.TargetArch, c.Config.BinaryName)

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageValue,
			Cmd: []string{
				"sh", "-c",
				script,
			},
			WorkingDir: containerWorkdir,
			Tty:        false,
		}, &container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: projectPath,
					Target: containerWorkdir,
				},
			},
		}, nil, nil, "")
	log.Debug().Msgf("Container created with ID: %s", resp.ID)
	if err != nil {
		panic(err)
	}

	defer cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	// Attendre que le conteneur termine
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	// Lire les logs
	logs, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}

	defer logs.Close()
	io.Copy(os.Stdout, logs)

	fmt.Println("✅ Compilation terminée. Binaire dans :", projectPath)
	return nil
}

func (c *DockerCompiler) GetConfig() *CompilerConfig {
	return c.Config
}

func NewDockerCompiler(cfg *CompilerConfig) *DockerCompiler {
	return &DockerCompiler{
		Config: cfg,
	}
}
