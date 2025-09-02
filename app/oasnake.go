/* Package app is the main app package */
package app

import (
	"github.com/louislouislouislouis/oasnake/app/cmd"
	"github.com/louislouislouislouis/oasnake/app/pkg/utils"
)

func Run() error {
	utils.ConfigureLogger()
	cmd := cmd.NewRootCmd()
	return cmd.Execute()
}
