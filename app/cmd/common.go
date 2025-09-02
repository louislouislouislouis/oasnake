package cmd

import (
	"github.com/rs/zerolog/log"
)

func handleError(err error) {
	log.Error().Err(err).Msg("An error occurred")
}
