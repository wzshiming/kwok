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
	"strings"
	"text/template"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/yaml"

	"sigs.k8s.io/kwok/pkg/apis/internalversion"
)

var (
	startTime = time.Now().Format(time.RFC3339Nano)

	defaultFuncMap = template.FuncMap{
		"Now": func() string {
			return time.Now().Format(time.RFC3339Nano)
		},
		"StartTime": func() string {
			return startTime
		},
		"YAML": func(s interface{}, indent ...int) (string, error) {
			d, err := yaml.Marshal(s)
			if err != nil {
				return "", err
			}

			data := string(d)
			if len(indent) == 1 && indent[0] > 0 {
				pad := strings.Repeat(" ", indent[0]*2)
				data = strings.ReplaceAll("\n"+data, "\n", "\n"+pad)
			}
			return data, nil
		},
	}
)

// HostController is a controller that manages the fake hosts
type HostController struct {
	nodes       *NodeController
	pods        *PodController
	broadcaster record.EventBroadcaster
	clientSet   kubernetes.Interface
}

// HostControllerConfig is the configuration for the HostController
type HostControllerConfig struct {
	EnableCNI                             bool
	ClientSet                             kubernetes.Interface
	ManageNodesWithAnnotationSelector     string
	ManageNodesWithLabelSelector          string
	ManageNodeWithName                    string
	DisregardStatusWithAnnotationSelector string
	DisregardStatusWithLabelSelector      string
	CIDR                                  string
	NodeIP                                string
	NodeName                              string
	NodePort                              int
	Parallelism                           int
	PodStages                             []*internalversion.Stage
	NodeStages                            []*internalversion.Stage
	FuncMap                               template.FuncMap
}

// NewHostController creates a new HostController
func NewHostController(conf HostControllerConfig) (*HostController, error) {
	var nodeSelectorFunc func(node *corev1.Node) bool
	switch {
	case conf.ManageNodesWithAnnotationSelector != "":
		selector, err := labels.Parse(conf.ManageNodesWithAnnotationSelector)
		if err != nil {
			return nil, err
		}
		nodeSelectorFunc = func(node *corev1.Node) bool {
			return selector.Matches(labels.Set(node.Annotations))
		}
	default:
		nodeSelectorFunc = func(node *corev1.Node) bool {
			return true
		}
	}

	eventBroadcaster := record.NewBroadcaster()
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "kwok_controller"})

	var lockPodsOnNodeFunc func(ctx context.Context, nodeName string) error

	nodes, err := NewNodeController(NodeControllerConfig{
		ClientSet:                             conf.ClientSet,
		NodeIP:                                conf.NodeIP,
		NodeName:                              conf.NodeName,
		NodePort:                              conf.NodePort,
		DisregardStatusWithAnnotationSelector: conf.DisregardStatusWithAnnotationSelector,
		DisregardStatusWithLabelSelector:      conf.DisregardStatusWithLabelSelector,
		ManageNodesWithLabelSelector:          conf.ManageNodesWithLabelSelector,
		ManageNodeWithName:                    conf.ManageNodeWithName,
		NodeSelectorFunc:                      nodeSelectorFunc,
		LockPodsOnNodeFunc: func(ctx context.Context, nodeName string) error {
			return lockPodsOnNodeFunc(ctx, nodeName)
		},
		Stages:              conf.NodeStages,
		LockNodeParallelism: conf.Parallelism,
		FuncMap:             conf.FuncMap,
		Recorder:            recorder,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create nodes controller: %w", err)
	}

	pods, err := NewPodController(PodControllerConfig{
		EnableCNI:                             conf.EnableCNI,
		ClientSet:                             conf.ClientSet,
		NodeIP:                                conf.NodeIP,
		CIDR:                                  conf.CIDR,
		DisregardStatusWithAnnotationSelector: conf.DisregardStatusWithAnnotationSelector,
		DisregardStatusWithLabelSelector:      conf.DisregardStatusWithLabelSelector,
		Stages:                                conf.PodStages,
		LockPodParallelism:                    conf.Parallelism,
		NodeGetFunc:                           nodes.Get,
		FuncMap:                               conf.FuncMap,
		Recorder:                              recorder,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create pods controller: %w", err)
	}

	lockPodsOnNodeFunc = pods.LockPodsOnNode

	n := &HostController{
		pods:        pods,
		nodes:       nodes,
		broadcaster: eventBroadcaster,
		clientSet:   conf.ClientSet,
	}

	return n, nil
}

// Start starts the controller
func (c *HostController) Start(ctx context.Context) error {
	c.broadcaster.StartRecordingToSink(&clientcorev1.EventSinkImpl{Interface: c.clientSet.CoreV1().Events("")})

	err := c.pods.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start pods controller: %w", err)
	}
	err = c.nodes.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start nodes controller: %w", err)
	}
	return nil
}
