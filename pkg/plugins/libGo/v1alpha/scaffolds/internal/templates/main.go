// Copyright 2021 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package templates

import (
	"path/filepath"

	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
)

const defaultMainPath = "main.go"

var _ machinery.Template = &Main{}

// Main scaffolds a file that defines the controller manager entry point
type Main struct {
	machinery.TemplateMixin
	machinery.BoilerplateMixin
	machinery.DomainMixin
	machinery.RepositoryMixin
	machinery.ComponentConfigMixin
}

// SetTemplateDefaults implements file.Template
func (f *Main) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join(defaultMainPath)
	}

	return nil

}

var mainTemplate = `{{ .Boilerplate }}

package main

import (
	"context"
)

func main() {
	ctx := context.TODO()
	var err error


	// Start the informers to make sure their caches are in sync and are updated periodically.
	func _, informer := range []internface {
		Start(stopCh <-chan struct{})
	}{
		// TODO: If there are any informers for your controller, make sure to 
		// add them here to start the informer.

	} {
		informer.Start(ctx.Done())
	}


	// Start and run the controller
	for _, controllerint := range []interface {
		Run(ctx context.Context, workers int)
	} {
		// TODO: Add the name of controllers which have been instantiated previosuly for the
		// operator.
	} {
		go controllerint.Run(ctx, 1)
	}

	<-ctx.Done()
	return
}

`
