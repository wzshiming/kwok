/*
Copyright 2024 The Kubernetes Authors.

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

package metrics

import (
	"fmt"
	"maps"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	corev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/kwok/pkg/utils/cel"
	"sigs.k8s.io/kwok/pkg/utils/slices"
)

// EnvironmentConfig holds configuration for a cel program
type EnvironmentConfig struct {
	EnableCompileCache bool
	EnableResultCache  bool

	Now                    func() time.Time
	StartedContainersTotal func(nodeName string) int64

	ContainerResourceUsage func(resourceName, podNamespace, podName, containerName string) float64
	PodResourceUsage       func(resourceName, podNamespace, podName string) float64
	NodeResourceUsage      func(resourceName, nodeName string) float64

	ContainerResourceCumulativeUsage func(resourceName, podNamespace, podName, containerName string) float64
	PodResourceCumulativeUsage       func(resourceName, podNamespace, podName string) float64
	NodeResourceCumulativeUsage      func(resourceName, nodeName string) float64
}

// NewEnvironment returns a Environment that is able to evaluate node metrics
func NewEnvironment(conf EnvironmentConfig) (*Environment, error) {
	const (
		nowOldName                    = "now"                    // deprecated
		startedContainersTotalOldName = "startedContainersTotal" // deprecated

		nowName                    = "Now"
		startedContainersTotalName = "StartedContainersTotal"

		usageName           = "Usage"
		cumulativeUsageName = "CumulativeUsage"
	)
	types := slices.Clone(cel.DefaultTypes)
	conversions := slices.Clone(cel.DefaultConversions)
	funcs := maps.Clone(cel.DefaultFuncs)
	methods := maps.Clone(cel.FuncsToMethods(cel.DefaultFuncs))

	if conf.Now != nil {
		funcs[nowOldName] = []any{conf.Now}
		funcs[nowName] = []any{conf.Now}
	}

	if conf.ContainerResourceUsage != nil {
		methods[usageName] = append(methods[usageName], func(pod corev1.Pod, resourceName string, containerName string) float64 {
			return conf.ContainerResourceUsage(resourceName, pod.Namespace, pod.Name, containerName)
		})
	}

	if conf.PodResourceUsage != nil {
		methods[usageName] = append(methods[usageName], func(pod corev1.Pod, resourceName string) float64 {
			return conf.PodResourceUsage(resourceName, pod.Namespace, pod.Name)
		})
	}

	if conf.NodeResourceUsage != nil {
		methods[usageName] = append(methods[usageName], func(node corev1.Node, resourceName string) float64 {
			return conf.NodeResourceUsage(resourceName, node.Name)
		})
	}

	if conf.ContainerResourceCumulativeUsage != nil {
		methods[cumulativeUsageName] = append(methods[cumulativeUsageName], func(pod corev1.Pod, resourceName string, containerName string) float64 {
			return conf.ContainerResourceCumulativeUsage(resourceName, pod.Namespace, pod.Name, containerName)
		})
	}

	if conf.PodResourceCumulativeUsage != nil {
		methods[cumulativeUsageName] = append(methods[cumulativeUsageName], func(pod corev1.Pod, resourceName string) float64 {
			return conf.PodResourceCumulativeUsage(resourceName, pod.Namespace, pod.Name)
		})
	}

	if conf.NodeResourceCumulativeUsage != nil {
		methods[cumulativeUsageName] = append(methods[cumulativeUsageName], func(node corev1.Node, resourceName string) float64 {
			return conf.NodeResourceCumulativeUsage(resourceName, node.Name)
		})
	}

	if conf.StartedContainersTotal != nil {
		startedContainersTotal := conf.StartedContainersTotal
		startedContainersTotalByNode := func(node corev1.Node) float64 {
			return float64(conf.StartedContainersTotal(node.Name))
		}
		methods[startedContainersTotalOldName] = append(methods[startedContainersTotalOldName], startedContainersTotal, startedContainersTotalByNode)
		funcs[startedContainersTotalOldName] = append(funcs[startedContainersTotalOldName], startedContainersTotal, startedContainersTotalByNode)

		methods[startedContainersTotalName] = append(methods[startedContainersTotalName], startedContainersTotal, startedContainersTotalByNode)
		funcs[startedContainersTotalName] = append(funcs[startedContainersTotalName], startedContainersTotal, startedContainersTotalByNode)
	}

	env, err := cel.NewEnvironment(cel.EnvironmentConfig{
		EnableCompileCache: conf.EnableCompileCache,

		Types:       types,
		Conversions: conversions,
		Methods:     methods,
		Funcs:       funcs,
		Vars: map[string]any{
			"node":      corev1.Node{},
			"pod":       corev1.Pod{},
			"container": corev1.Container{},
		},
	})
	if err != nil {
		return nil, err
	}
	e := &Environment{
		env:  env,
		conf: conf,
	}

	if conf.EnableResultCache {
		e.resultCacheVer = new(int64)
	}

	return e, nil
}

// Environment is environment in which cel programs are executed
type Environment struct {
	env *cel.Environment

	conf           EnvironmentConfig
	resultCacheVer *int64
}

// Compile is responsible for compiling a cel program
func (e *Environment) Compile(src string) (*Evaluator, error) {
	eval, err := e.env.Compile(src)
	if err != nil {
		return nil, fmt.Errorf("failed to compile metric expression: %w", err)
	}

	evaluator := &Evaluator{
		evaluator:      eval,
		latestCacheVer: e.resultCacheVer,
	}
	return evaluator, nil
}

// ClearResultCache clears the result cache
func (e *Environment) ClearResultCache() {
	if e.resultCacheVer == nil {
		return
	}
	atomic.AddInt64(e.resultCacheVer, 1)
}

// Evaluator evaluates a cel program
type Evaluator struct {
	evaluator cel.Program

	latestCacheVer *int64
	cacheVer       int64

	cache    map[string]ref.Val
	cacheMut sync.Mutex
}

func resultUniqueKey(node *corev1.Node, pod *corev1.Pod, container *corev1.Container) string {
	tmp := make([]string, 0, 5)
	if node != nil {
		tmp = append(tmp, string(node.UID), node.ResourceVersion)
	}
	if pod != nil {
		tmp = append(tmp, string(pod.UID), pod.ResourceVersion)
	}
	if container != nil {
		tmp = append(tmp, container.Name)
	}
	return strings.Join(tmp, "/")
}

func (e *Evaluator) evaluate(data Data) (ref.Val, error) {
	var key string
	if e.latestCacheVer != nil {
		e.cacheMut.Lock()
		defer e.cacheMut.Unlock()

		if e.cache == nil || *e.latestCacheVer != e.cacheVer {
			e.cache = map[string]ref.Val{}
			e.cacheVer = *e.latestCacheVer
		}

		key = resultUniqueKey(data.Node, data.Pod, data.Container)
		if val, ok := e.cache[key]; ok {
			return val, nil
		}
	}

	refVal, _, err := e.evaluator.Eval(map[string]any{
		"node":      data.Node,
		"pod":       data.Pod,
		"container": data.Container,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate metric expression: %w", err)
	}
	if key != "" {
		e.cache[key] = refVal
	}
	return refVal, nil
}

// EvaluateFloat64 evaluates a cel program and returns a float64.
func (e *Evaluator) EvaluateFloat64(data Data) (float64, error) {
	refVal, err := e.evaluate(data)
	if err != nil {
		return 0, err
	}

	switch v := refVal.(type) {
	case types.Duration:
		return float64(v.Duration), nil
	case types.Int:
		return float64(v), nil
	case types.Double:
		return float64(v), nil
	case types.Uint:
		return float64(v), nil
	case types.Bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case cel.Quantity:
		return v.Quantity.AsApproximateFloat64(), nil
	default:
		return 0, fmt.Errorf("unsupported metric value type: %T", v)
	}
}

// EvaluateString evaluates a cel program and returns a string
func (e *Evaluator) EvaluateString(data Data) (string, error) {
	refVal, err := e.evaluate(data)
	if err != nil {
		return "", err
	}

	v, ok := refVal.(types.String)
	if !ok {
		return "", fmt.Errorf("unsupported metric type: %T", v)
	}
	return string(v), nil
}

// Data is a data structure that is passed to the cel program
type Data struct {
	Node      *corev1.Node
	Pod       *corev1.Pod
	Container *corev1.Container
}
