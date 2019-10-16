package controller

import (
	"github.com/example-inc/zuul-operator/pkg/controller/zuul"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, zuul.Add)
}
