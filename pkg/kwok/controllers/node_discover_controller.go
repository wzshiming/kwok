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
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	"sigs.k8s.io/kwok/pkg/log"
)

// NodeDiscoverController is the controller for node discovery
type NodeDiscoverController struct {
	clientSet                    kubernetes.Interface
	nodeSelectorFunc             func(node *corev1.Node) bool
	manageNodesWithLabelSelector string
	addFunc                      func(node *corev1.Node) error
	deleteFunc                   func(node *corev1.Node) error
}

// NodeDiscoverControllerConfig is the config for NodeDiscoverController
type NodeDiscoverControllerConfig struct {
	ClientSet                    kubernetes.Interface
	NodeSelectorFunc             func(node *corev1.Node) bool
	ManageNodesWithLabelSelector string
	AddFunc                      func(node *corev1.Node) error
	DeleteFunc                   func(node *corev1.Node) error
}

// NewNodeDiscoverController creates a new NodeDiscoverController
func NewNodeDiscoverController(conf NodeDiscoverControllerConfig) (*NodeDiscoverController, error) {
	c := &NodeDiscoverController{
		clientSet:                    conf.ClientSet,
		nodeSelectorFunc:             conf.NodeSelectorFunc,
		manageNodesWithLabelSelector: conf.ManageNodesWithLabelSelector,
	}
	return c, nil
}

// Start starts the controller
func (c *NodeDiscoverController) Start(ctx context.Context) error {
	opt := metav1.ListOptions{}
	opt.LabelSelector = c.manageNodesWithLabelSelector
	opt.ResourceVersion = "0"

	err := c.WatchNodes(ctx, opt)
	if err != nil {
		return fmt.Errorf("failed watch node: %w", err)
	}
	return nil
}

// WatchNodes watches nodes in the cluster
func (c *NodeDiscoverController) WatchNodes(ctx context.Context, opt metav1.ListOptions) error {
	// Watch nodes in the cluster
	watcher, err := c.clientSet.CoreV1().Nodes().Watch(ctx, opt)
	if err != nil {
		return err
	}

	logger := log.FromContext(ctx)
	go func() {
		rc := watcher.ResultChan()
	loop:
		for {
			select {
			case event, ok := <-rc:
				if !ok {
					for {
						watcher, err := c.clientSet.CoreV1().Nodes().Watch(ctx, opt)
						if err == nil {
							rc = watcher.ResultChan()
							continue loop
						}

						logger.Error("Failed to watch nodes", err)
						select {
						case <-ctx.Done():
							break loop
						case <-time.After(time.Second * 5):
						}
					}
				}
				switch event.Type {
				case watch.Added:
					node := event.Object.(*corev1.Node)
					if c.needLockNode(node) {
						err = c.addFunc(node)
						if err != nil {
							logger.Error("Failed to add node", err)
						}
					}
				case watch.Deleted:
					node := event.Object.(*corev1.Node)
					if c.needLockNode(node) {
						err = c.deleteFunc(node)
						if err != nil {
							logger.Error("Failed to delete node", err)
						}
					}
				}
			case <-ctx.Done():
				watcher.Stop()
				break loop
			}
		}
		logger.Info("Stop watch nodes")
	}()
	return nil
}

func (c *NodeDiscoverController) needLockNode(node *corev1.Node) bool {
	if !c.nodeSelectorFunc(node) {
		return false
	}
	return true
}
