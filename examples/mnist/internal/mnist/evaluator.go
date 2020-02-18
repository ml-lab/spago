// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mnist

import (
	"github.com/gosuri/uiprogress"
	"saientist.dev/spago/pkg/mat"
	"saientist.dev/spago/pkg/mat/f64utils"
	"saientist.dev/spago/pkg/ml/ag"
	"saientist.dev/spago/pkg/ml/nn"
	"saientist.dev/spago/pkg/ml/stats"
)

type Evaluator struct {
	model nn.Model
}

func NewEvaluator(model nn.Model) *Evaluator {
	e := &Evaluator{
		model: model,
	}
	return e
}

// Predict performs the forward pass and returns the predict label
func (t *Evaluator) Predict(image *mat.Dense) int {
	g := ag.NewGraph()
	x := g.NewVariable(image, false)
	y := t.model.NewProc(g).Forward(x)[0]
	return f64utils.ArgMax(y.Value().Data())
}

func (t *Evaluator) Evaluate(dataset Dataset) *stats.ClassMetrics {
	uip := uiprogress.New()
	bar := newTestBar(uip, dataset)
	uip.Start()
	defer uip.Stop()

	counter := stats.NewMetricCounter()
	for i := 0; i < dataset.Count(); i++ {
		image, label := dataset.GetNormalized(i)
		if t.Predict(image) == int(label) {
			counter.IncTruePos()
		} else {
			counter.IncFalsePos()
		}
		bar.Incr()
	}
	return counter
}

func newTestBar(p *uiprogress.Progress, dataset Dataset) *uiprogress.Bar {
	bar := p.AddBar(dataset.Count())
	bar.AppendCompleted().PrependElapsed()
	return bar
}