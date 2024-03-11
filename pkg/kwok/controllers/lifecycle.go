/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/cel-go/common/types"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/google/cel-go/common/types"
	"sigs.k8s.io/kwok/pkg/apis/internalversion"
	"sigs.k8s.io/kwok/pkg/utils/cel"
	"sigs.k8s.io/kwok/pkg/utils/expression"
	"sigs.k8s.io/kwok/pkg/utils/format"
)

// NewLifecycle returns a new Lifecycle.
func NewLifecycle(stages []*internalversion.Stage, env *cel.Environment) (Lifecycle, error) {
	lcs := Lifecycle{}
	for _, stage := range stages {
		lc, err := NewLifecycleStage(stage, env)
		if err != nil {
			return nil, fmt.Errorf("lifecycle stage: %w", err)
		}
		if lc == nil {
			continue
		}
		lcs = append(lcs, lc)
	}
	return lcs, nil
}

// Lifecycle is a list of lifecycle stage.
type Lifecycle []*LifecycleStage

func (s Lifecycle) match(label, annotation labels.Set, data, jsonStandard any) ([]*LifecycleStage, error) {
	out := []*LifecycleStage{}
	for _, stage := range s {
		ok, err := stage.match(label, annotation, data, jsonStandard)
		if err != nil {
			return nil, err
		}
		if ok {
			out = append(out, stage)
		}
	}
	return out, nil
}

// Match returns matched stage.
func (s Lifecycle) Match(label, annotation labels.Set, data, jsonStandard any) (*LifecycleStage, error) {
	stages, err := s.match(label, annotation, data, jsonStandard)
	if err != nil {
		return nil, err
	}
	if len(stages) == 0 {
		return nil, nil
	}
	if len(stages) == 1 {
		return stages[0], nil
	}
	totalWeights := 0
	for _, stage := range stages {
		totalWeights += stage.weight
	}
	if totalWeights == 0 {
		//nolint:gosec
		return stages[rand.Intn(len(stages))], nil
	}

	//nolint:gosec
	off := rand.Intn(totalWeights)
	for _, stage := range stages {
		if stage.weight == 0 {
			continue
		}
		off -= stage.weight
		if off < 0 {
			return stage, nil
		}
	}
	return stages[len(stages)-1], nil
}

// NewLifecycleStage returns a new LifecycleStage.
func NewLifecycleStage(s *internalversion.Stage, env *cel.Environment) (*LifecycleStage, error) {
	stage := &LifecycleStage{
		name: s.Name,
	}
	selector := s.Spec.Selector
	if selector == nil {
		return nil, nil
	}

	if selector.MatchLabels != nil {
		stage.matchLabels = labels.SelectorFromSet(selector.MatchLabels)
	}
	if selector.MatchAnnotations != nil {
		stage.matchAnnotations = labels.SelectorFromSet(selector.MatchAnnotations)
	}
	if selector.MatchExpressions != nil {
		for _, express := range selector.MatchExpressions {
			requirement, err := expression.NewRequirement(express.Key, express.Operator, express.Values)
			if err != nil {
				return nil, err
			}
			stage.matchExpressions = append(stage.matchExpressions, requirement)
		}
	}
	if selector.MatchConditions != nil {
		for _, express := range selector.MatchConditions {
			program, err := env.Compile(express.Expression)
			if err != nil {
				return nil, err
			}
			stage.matchConditions = append(stage.matchConditions, program)
		}
	}

	stage.next = &s.Spec.Next
	if delay := s.Spec.Delay; delay != nil {
		var durationFrom *string
		if delay.DurationFrom != nil {
			durationFrom = &delay.DurationFrom.ExpressionFrom
		}
		var delayDuration time.Duration
		if delay.DurationMilliseconds != nil {
			delayDuration = time.Duration(*delay.DurationMilliseconds) * time.Millisecond
		}
		duration, err := expression.NewDurationFrom(&delayDuration, durationFrom)
		if err != nil {
			return nil, err
		}
		stage.duration = duration

		if delay.JitterDurationMilliseconds != nil || delay.JitterDurationFrom != nil {
			var jitterDurationFrom *string
			if delay.JitterDurationFrom != nil {
				jitterDurationFrom = &delay.JitterDurationFrom.ExpressionFrom
			}
			var jitterDuration *time.Duration
			if delay.JitterDurationMilliseconds != nil {
				jitterDuration = format.Ptr(time.Duration(*delay.JitterDurationMilliseconds) * time.Millisecond)
			}
			jitterDurationGetter, err := expression.NewDurationFrom(jitterDuration, jitterDurationFrom)
			if err != nil {
				return nil, err
			}
			stage.jitterDuration = jitterDurationGetter
		}
	}

	if weight := s.Spec.Weight; weight > 1 {
		stage.weight = weight
	} else {
		stage.weight = 0
	}

	stage.immediateNextStage = s.Spec.ImmediateNextStage

	return stage, nil
}

// LifecycleStage is a resource lifecycle stage manager
type LifecycleStage struct {
	name             string
	matchLabels      labels.Selector
	matchAnnotations labels.Selector
	matchExpressions []*expression.Requirement
	matchConditions  []cel.Program

	weight int
	next   *internalversion.StageNext

	duration       expression.DurationGetter
	jitterDuration expression.DurationGetter

	immediateNextStage bool
}

func (s *LifecycleStage) match(label, annotation labels.Set, data, jsonStandard any) (bool, error) {
	if s.matchLabels != nil {
		if !s.matchLabels.Matches(label) {
			return false, nil
		}
	}
	if s.matchAnnotations != nil {
		if !s.matchAnnotations.Matches(annotation) {
			return false, nil
		}
	}

	if s.matchExpressions != nil {
		for _, requirement := range s.matchExpressions {
			ok, err := requirement.Matches(context.Background(), jsonStandard)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
		}
	}

	if s.matchConditions != nil {
		for _, program := range s.matchConditions {
			val, _, err := program.Eval(data)
			if err != nil {
				return false, err
			}
			if !val.(types.Bool) {
				return false, nil
			}
		}
	}

	return true, nil
}

// Delay returns the delay duration of the stage.
// It's not a constant value, it can be a random value.
func (s *LifecycleStage) Delay(ctx context.Context, v interface{}, now time.Time) (time.Duration, bool) {
	if s.duration == nil {
		return 0, false
	}

	duration, ok := s.duration.Get(ctx, v, now)
	if !ok {
		return 0, false
	}

	if s.jitterDuration == nil {
		return duration, true
	}

	jitterDuration, ok := s.jitterDuration.Get(ctx, v, now)
	if !ok {
		return duration, true
	}

	if jitterDuration < duration {
		return jitterDuration, true
	}

	//nolint:gosec
	return duration + time.Duration(rand.Int63n(int64(jitterDuration-duration))), true
}

// Next returns the next of the stage.
func (s *LifecycleStage) Next() *internalversion.StageNext {
	return s.next
}

// Name returns the name of the stage
func (s *LifecycleStage) Name() string {
	return s.name
}

// ImmediateNextStage returns whether the stage is immediate next stage.
func (s *LifecycleStage) ImmediateNextStage() bool {
	return s.immediateNextStage
}
