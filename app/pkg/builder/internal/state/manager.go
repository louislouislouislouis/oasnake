package state

import (
	"fmt"

	"github.com/louislouislouislouis/oasnake/app/pkg/builder/internal/state/events"
	"github.com/rs/zerolog/log"
)

type StateManager struct {
	current        State
	transitions    map[State]map[events.EventType]State
	stateFunctions map[State]StateFunc
}

func NewStateManager(stateFuncs map[State]StateFunc) *StateManager {
	return &StateManager{
		current:        Ready,
		stateFunctions: stateFuncs,
		transitions: map[State]map[events.EventType]State{
			Ready: {
				events.StartParsing: Parsing,
			},
			Parsing: {
				events.FinishParsing: Generating,
			},
			Generating: {
				events.FinishGenerateCode: WithCode,
			},
			WithCode: {
				events.End:          Success,
				events.StartCompile: Compiling,
			},
			Compiling: {
				events.FinishCompile: WithBinary,
			},
			WithBinary: {
				events.End:          Success,
				events.StartInstall: Installing,
			},
			Installing: {
				events.FinishInstall: Installed,
			},
			Installed: {
				events.End: Success,
			},
		},
	}
}

func (sm *StateManager) Accept(event events.Event) error {
	log.Debug().Msgf("received event: %s in state %s", event.Type().String(), sm.current.String())

	if event.Type() == events.Error {
		sm.current = Failure
		return event.(events.ErrorEvent).Error
	}
	next, ok := sm.transitions[sm.current][event.Type()]
	if !ok {
		err := fmt.Errorf("event %s not allowed in state %s", event.Type().String(), sm.current.String())
		log.Error().Err(err).Msg("state transition error")
		sm.current = Failure
		return err
	}

	log.Debug().Msgf("transitioning to %s", next.String())
	sm.current = next
	log.Info().Msgf("%s", sm.current.String())
	handler, exists := sm.stateFunctions[sm.current]
	if !exists {
		log.Debug().Msgf("No state function for state %s", sm.current.String())
	}
	if exists {
		log.Debug().Msgf("State function of state %s will be launched after receiving event %s", sm.current.String(), event.Type().String())
		return sm.Accept(handler(event))
	}
	return nil
}

func (sm *StateManager) On(state State, fn StateFunc) {
	if sm.stateFunctions == nil {
		sm.stateFunctions = make(map[State]StateFunc)
	}
	sm.stateFunctions[state] = fn
}
