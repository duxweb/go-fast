package validator

import (
	"github.com/duxweb/go-fast/response"
	"github.com/go-playground/validator/v10"
)

type ValidatorWarp struct {
	Rule    string
	Message string
}

type ValidatorRule map[string]ValidatorWarp

func ValidatorMaps(params map[string]any, rules ValidatorRule) error {
	r := map[string]any{}
	for name, warp := range rules {
		r[name] = warp.Rule
	}
	validateErr := Validator().ValidateMap(params, r)
	err := validatorMapsError(rules, validateErr)
	if err != nil {
		return err
	}
	return nil
}

func validatorMapsError(rules ValidatorRule, errs map[string]any) error {
	if len(errs) == 0 {
		return nil
	}
	data := map[string]any{}
	message := ""
	for name, err := range errs {
		x := err.(validator.ValidationErrors)

		e := ""
		if val, ok := rules[name]; ok {
			e = val.Message
		} else {
			e = x.Error()
		}
		if message == "" {
			message = e
		}
		data[name] = []string{e}
	}
	return response.ValidatorError(message, data)
}
