package app

import "github.com/louislouislouislouis/oasnake/app/cmd"

func Run() error {
	cmd := cmd.NewRootCmd()
	return cmd.Execute()
}
