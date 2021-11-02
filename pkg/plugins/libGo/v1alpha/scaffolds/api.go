/*
Copyright 2021 The Kubernetes Authors.

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

package scaffolds

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/varshaprasad96/lib-go-plugins/pkg/plugins/libGo/v1alpha/scaffolds/internal/templates/api"
	"github.com/varshaprasad96/lib-go-plugins/pkg/plugins/libGo/v1alpha/scaffolds/internal/templates/controllers"
	"github.com/varshaprasad96/lib-go-plugins/pkg/plugins/libGo/v1alpha/scaffolds/internal/templates/hack"
	"sigs.k8s.io/kubebuilder/v3/pkg/config"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v3/pkg/model/resource"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugins"
)

var _ plugins.Scaffolder = &apiScaffolder{}

type apiScaffolder struct {
	config   config.Config
	resource resource.Resource

	fs machinery.Filesystem

	force bool
}

func NewAPIScaffolder(config config.Config, res resource.Resource, force bool) plugins.Scaffolder {
	return &apiScaffolder{
		config:   config,
		resource: res,
		force:    force,
	}
}

func (s *apiScaffolder) InjectFS(fs machinery.Filesystem) {
	s.fs = fs
}

func (s *apiScaffolder) Scaffold() error {
	fmt.Println("Writing scaffold for you to edit")

	// Load the boilerplate
	boilerplate, err := afero.ReadFile(s.fs.FS, hack.DefaultBoilerplatePath)
	if err != nil {
		return fmt.Errorf("error scaffolding API/controller: unable to load boilerplate: %w", err)
	}

	// Initialize the machinery.Scaffold that will write files to disk
	scaffold := machinery.NewScaffold(s.fs,
		machinery.WithConfig(s.config),
		machinery.WithBoilerplate(string(boilerplate)),
		machinery.WithResource(&s.resource))

	// Keep track of the values before the update
	doAPI := s.resource.HasAPI()
	doController := s.resource.HasController()

	if err := s.config.UpdateResource(s.resource); err != nil {
		return fmt.Errorf("error updating the resource: %w", err)
	}

	if doAPI {
		if err := scaffold.Execute(
			&api.Types{Force: s.force},
			&api.Group{},
		); err != nil {
			return fmt.Errorf("error scaffolding apis: %v", err)
		}
	}

	if doController {
		if err := scaffold.Execute(
			&controllers.Controller{Force: s.force},
		); err != nil {
			return fmt.Errorf("error scaffolding controller: %v", err)
		}
	}
	return nil
}
