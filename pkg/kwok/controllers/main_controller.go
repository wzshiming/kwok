/*
Copyright 2023 The Kubernetes Authors.

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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"

	"sigs.k8s.io/kwok/pkg/apis/internalversion"
	"sigs.k8s.io/kwok/pkg/kwok/clientsets"
)

// MainController is a main controller
type MainController struct {
	conf MainControllerConfig
}

// MainControllerConfig is the config for MainController
type MainControllerConfig struct {
	EnableCNI                             bool
	SeparateForNodes                      bool
	ClientSet                             clientsets.Clientsets
	ManageAllNodes                        bool
	ManageNodesWithAnnotationSelector     string
	ManageNodesWithLabelSelector          string
	DisregardStatusWithAnnotationSelector string
	DisregardStatusWithLabelSelector      string
	CIDR                                  string
	NodeIP                                string
	NodeName                              string
	NodePort                              int
	PodStages                             []*internalversion.Stage
	NodeStages                            []*internalversion.Stage
}

// NewMainController creates a new MainController
func NewMainController(conf MainControllerConfig) (*MainController, error) {
	switch {
	case conf.ManageAllNodes:
		conf.ManageNodesWithAnnotationSelector = ""
		conf.ManageNodesWithLabelSelector = ""
	}

	n := &MainController{
		conf: conf,
	}
	return n, nil
}

// Start starts the controller
func (c *MainController) Start(ctx context.Context) error {
	conf := c.conf
	var nodeSelectorFunc func(node *corev1.Node) bool
	switch {
	case conf.ManageNodesWithAnnotationSelector != "":
		selector, err := labels.Parse(conf.ManageNodesWithAnnotationSelector)
		if err != nil {
			return err
		}
		nodeSelectorFunc = func(node *corev1.Node) bool {
			return selector.Matches(labels.Set(node.Annotations))
		}
	default:
		nodeSelectorFunc = func(node *corev1.Node) bool {
			return true
		}
	}

	if conf.SeparateForNodes {
		clientset, err := conf.ClientSet.Get(ctx, clientsets.InternalClientsetKey)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(ctx)

		nodeDiscoverController, err := NewNodeDiscoverController(NodeDiscoverControllerConfig{
			ClientSet:                    clientset,
			ManageNodesWithLabelSelector: conf.ManageNodesWithLabelSelector,
			NodeSelectorFunc:             nodeSelectorFunc,
			AddFunc: func(node *corev1.Node) error {
				clientset, err := conf.ClientSet.Get(ctx, node.Name)
				if err != nil {
					return err
				}
				hosts, err := NewHostController(HostControllerConfig{
					ClientSet:                             clientset,
					EnableCNI:                             conf.EnableCNI,
					ManageNodeWithName:                    node.Name,
					DisregardStatusWithAnnotationSelector: conf.DisregardStatusWithAnnotationSelector,
					DisregardStatusWithLabelSelector:      conf.DisregardStatusWithLabelSelector,
					CIDR:                                  conf.CIDR,
					NodeIP:                                conf.NodeIP,
					NodeName:                              conf.NodeName,
					NodePort:                              conf.NodePort,
					PodStages:                             conf.PodStages,
					NodeStages:                            conf.NodeStages,
					FuncMap:                               defaultFuncMap,
					Parallelism:                           1,
				})
				if err != nil {
					return err
				}

				err = hosts.Start(ctx)
				if err != nil {
					return fmt.Errorf("failed to start hosts controller: %w", err)
				}
				return err
			},
			DeleteFunc: func(node *corev1.Node) error {
				cancel()
				err := conf.ClientSet.Close(context.Background(), node.Name)
				if err != nil {
					return err
				}
				return nil
			},
		})
		if err != nil {
			return err
		}

		err = nodeDiscoverController.Start(ctx)
		if err != nil {
			return fmt.Errorf("failed to start node discover controller: %w", err)
		}

	} else {
		clientset, err := conf.ClientSet.Get(ctx, clientsets.InternalClientsetKey)
		if err != nil {
			return err
		}

		hosts, err := NewHostController(HostControllerConfig{
			ClientSet:                             clientset,
			EnableCNI:                             conf.EnableCNI,
			ManageNodesWithAnnotationSelector:     conf.ManageNodesWithAnnotationSelector,
			ManageNodesWithLabelSelector:          conf.ManageNodesWithLabelSelector,
			DisregardStatusWithAnnotationSelector: conf.DisregardStatusWithAnnotationSelector,
			DisregardStatusWithLabelSelector:      conf.DisregardStatusWithLabelSelector,
			CIDR:                                  conf.CIDR,
			NodeIP:                                conf.NodeIP,
			NodeName:                              conf.NodeName,
			NodePort:                              conf.NodePort,
			PodStages:                             conf.PodStages,
			NodeStages:                            conf.NodeStages,
			FuncMap:                               defaultFuncMap,
			Parallelism:                           16,
		})
		if err != nil {
			return err
		}

		err = hosts.Start(ctx)
		if err != nil {
			return fmt.Errorf("failed to start hosts controller: %w", err)
		}
	}

	return nil
}
