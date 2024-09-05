package gofilter

import (
	"sync"

	"github.com/Knetic/govaluate"
)

type GValuate struct {
	Name    string
	lck     sync.Mutex
	Valuate *govaluate.EvaluableExpression
}

func NewGValuate(name string) *GValuate {
	return &GValuate{
		Name:    name,
		lck:     sync.Mutex{},
		Valuate: nil,
	}
}

func (gv *GValuate) Eval(exp string) error {
	gv.lck.Lock()
	defer gv.lck.Unlock()
	var err error

	if gv.Valuate == nil {
		gv.Valuate, err = govaluate.NewEvaluableExpression(exp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gv *GValuate) Exec(parameters map[string]interface{}) (bool, error) {
	result, err := gv.Valuate.Evaluate(parameters)
	if err != nil {
		return false, err
	}

	return result.(bool), nil
}
