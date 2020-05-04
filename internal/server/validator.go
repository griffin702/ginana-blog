package server

import (
	"ginana-blog/internal/model"
	"github.com/griffin702/ginana/library/validator"
	"github.com/kataras/iris/v12"
)

func NewValidator() (valid model.ValidatorHandler, err error) {
	return func(ctx iris.Context) model.Validator {
		v := validator.NewValidator()
		return v.ValidateStruct
	}, nil
}
