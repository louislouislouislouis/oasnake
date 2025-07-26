/* Package compiler handle generation of a binary file*/
package compiler

import "github.com/rs/zerolog/log"

type Compiler struct {
	Config *CompilerConfig
}

func (c *Compiler) Compile() error {
	log.Debug().Msg("Compiling binary...")
	return nil
}

func NewCompiler(cfg *CompilerConfig) *Compiler {
	return &Compiler{
		Config: cfg,
	}
}
