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

package cel

import (
	"fmt"
	"sync"

	"github.com/google/cel-go/cel"
	"github.com/wzshiming/easycel"
)

// EnvironmentConfig holds configuration for a cel program
type EnvironmentConfig struct {
	EnableCompileCache bool

	Conversions []any
	Types       []any
	Vars        map[string]any
	Funcs       map[string][]any
	Methods     map[string][]any
}

// NewEnvironment returns a Environment with the given configuration
func NewEnvironment(conf EnvironmentConfig) (*Environment, error) {
	registry := easycel.NewRegistry("kwok.ext.node",
		easycel.WithTagName("json"),
	)

	e := &Environment{
		registry: registry,

		conversions: conf.Conversions,
		types:       conf.Types,
		vars:        conf.Vars,
		funcs:       conf.Funcs,
		methods:     conf.Methods,
	}

	if conf.EnableCompileCache {
		e.cacheProgram = map[string]cel.Program{}
	}

	err := e.init()
	if err != nil {
		return nil, err
	}

	env, err := easycel.NewEnvironment(cel.Lib(registry))
	if err != nil {
		return nil, fmt.Errorf("failed to create CEL environment: %w", err)
	}
	e.env = env

	return e, nil
}

// Environment is environment in which cel programs are executed
type Environment struct {
	registry     *easycel.Registry
	env          *easycel.Environment
	cacheProgram map[string]cel.Program
	cacheMut     sync.Mutex

	conversions []any
	types       []any
	vars        map[string]any
	funcs       map[string][]any
	methods     map[string][]any
}

func (e *Environment) init() error {
	for _, convert := range e.conversions {
		err := e.registry.RegisterConversion(convert)
		if err != nil {
			return fmt.Errorf("failed to register convert %T: %w", convert, err)
		}
	}
	for _, typ := range e.types {
		err := e.registry.RegisterType(typ)
		if err != nil {
			return fmt.Errorf("failed to register type %T: %w", typ, err)
		}
	}
	for name, val := range e.vars {
		err := e.registry.RegisterVariable(name, val)
		if err != nil {
			return fmt.Errorf("failed to register variable %s: %w", name, err)
		}
	}
	for name, list := range e.funcs {
		for _, fun := range list {
			err := e.registry.RegisterFunction(name, fun)
			if err != nil {
				return fmt.Errorf("failed to register function %s: %w", name, err)
			}
		}
	}
	for name, list := range e.methods {
		for _, fun := range list {
			err := e.registry.RegisterMethod(name, fun)
			if err != nil {
				return fmt.Errorf("failed to register method %s: %w", name, err)
			}
		}
	}

	return nil
}

// Program is a cel program
type Program = cel.Program

// Compile is responsible for compiling a cel program
func (e *Environment) Compile(src string) (cel.Program, error) {
	if e.cacheProgram != nil {
		e.cacheMut.Lock()
		defer e.cacheMut.Unlock()

		if program, ok := e.cacheProgram[src]; ok {
			return program, nil
		}
	}
	program, err := e.env.Program(src)
	if err != nil {
		return nil, fmt.Errorf("failed to compile expression: %w", err)
	}

	if e.cacheProgram != nil {
		e.cacheProgram[src] = program
	}

	return program, nil
}
