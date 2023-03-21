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

package clientsets

import (
	"net"
	"sync"
)

func newConnExitHook(conn net.Conn, exitFunc func(c *connExitHook)) *connExitHook {
	return &connExitHook{
		Conn:     conn,
		exitFunc: exitFunc,
	}
}

type connExitHook struct {
	once     sync.Once
	exitFunc func(c *connExitHook)
	net.Conn
}

func (c *connExitHook) callHook() {
	c.once.Do(func() {
		c.exitFunc(c)
	})
}

func (c *connExitHook) Close() error {
	c.callHook()
	return c.Conn.Close()
}

func (c *connExitHook) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	if err != nil {
		c.callHook()
	}
	return n, err
}

func (c *connExitHook) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	if err != nil {
		c.callHook()
	}
	return n, err
}
