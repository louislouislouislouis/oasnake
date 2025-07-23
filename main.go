/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"os"

	"github.com/louislouislouislouis/oasnake/app"
)

func main() {
	log.SetOutput(os.Stdout)
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
