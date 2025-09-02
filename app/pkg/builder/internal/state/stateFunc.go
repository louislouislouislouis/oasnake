package state

import "github.com/louislouislouislouis/oasnake/app/pkg/builder/internal/state/events"

type StateFunc func(event events.Event) events.Event
