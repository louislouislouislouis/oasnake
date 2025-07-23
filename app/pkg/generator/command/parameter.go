package command

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/louislouislouislouis/oasnake/app/pkg/utils"
)

type Parameter struct {
	openapi3.Parameter
}

func NewParameter(param openapi3.Parameter) *Parameter {
	return &Parameter{
		param,
	}
}

func (p Parameter) GetSafeDescription() string {
	return utils.RemoveBackTicks(p.Description)
}
