/* Package compiler handle generation of a binary file*/
package compiler

type Compiler interface {
	Compile() error
	GetConfig() *CompilerConfig
}
