package controller

import (
	"github.com/leg100/bucket-operator/pkg/controller/bucket"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, bucket.Add)
}
