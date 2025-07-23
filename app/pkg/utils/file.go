package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FS struct {
	Directory string
	Filename  string
}

type WriterConfig struct {
	OutputDirectoryShouldBeEmpty bool
	Output                       FS
	Content                      string
}

func WriteFileContent(cfg WriterConfig) error {
	if cfg.OutputDirectoryShouldBeEmpty {
		_, err := os.Stat(cfg.Output.Directory)
		if err == nil {
			return fmt.Errorf("output directory already exists: %s", cfg.Output.Directory)
		}
		if !os.IsNotExist(err) {
			return fmt.Errorf("error checking output directory: %v", err)
		}

	}
	err := os.MkdirAll(cfg.Output.Directory, 0755)
	if err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	filePath := filepath.Join(cfg.Output.Directory, cfg.Output.Filename)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(cfg.Content)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}

func CopyFile(sourcePath, destPath string) error {
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("error while opening source file: %v", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error while creating destination file: %v", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("error while copying file: %v", err)
	}
	return nil
}

func WriteEmbeddedFile(outputPath string, file []byte) error {
	dir := filepath.Dir(outputPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directories: %w", err)
	}
	err := os.WriteFile(outputPath, file, 0644)
	if err != nil {
		return fmt.Errorf("error writing embedded file: %w", err)
	}
	return nil
}
