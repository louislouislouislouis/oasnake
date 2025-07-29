package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type GoCompiler struct {
	Config *CompilerConfig
}

func (c *GoCompiler) Compile() error {
	log.Debug().Msg("Compiling binary using local go compiler...")

	projectPath := filepath.Join(c.Config.OutputDirectory)

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error().Err(err).Msg("Failed to run go mod tidy")
		return err
	}

	cmd = exec.Command("go", "build", "-o", c.Config.BinaryName)
	cmd.Dir = projectPath
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", c.Config.TargetOs),
		fmt.Sprintf("GOARCH=%s", c.Config.TargetArch),
	)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Error().Err(err).Msg("Failed to compile binary")
		return err
	}

	log.Debug().Msgf("✅ Compilation succeeded : binary %s/%s created ", c.Config.OutputDirectory, c.Config.BinaryName)
	return nil
}

func (c *GoCompiler) GetConfig() *CompilerConfig {
	return c.Config
}

func NewGoCompiler(cfg *CompilerConfig) *GoCompiler {
	return &GoCompiler{
		Config: cfg,
	}
}
