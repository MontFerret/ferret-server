package dto

import "github.com/MontFerret/ferret-server/server/http/api/models"

func ExecutionParamsFrom(params map[string]interface{}) map[string]models.Any {
	res := map[string]models.Any{}

	if params == nil {
		return res
	}

	for k, v := range params {
		res[k] = v
	}

	return res
}

func ExecutionParamsTo(params map[string]models.Any) map[string]interface{} {
	res := map[string]interface{}{}

	if params == nil {
		return res
	}

	for k, v := range params {
		res[k] = v
	}

	return res
}
