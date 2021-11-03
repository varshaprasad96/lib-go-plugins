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

package scaffolds

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/varshaprasad96/lib-go-plugins/pkg/plugins/libGo/v1alpha/scaffolds/internal/templates"
	"github.com/varshaprasad96/lib-go-plugins/pkg/plugins/libGo/v1alpha/scaffolds/internal/templates/hack"
	"sigs.k8s.io/kubebuilder/v3/pkg/config"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugins"
)

const (
	// ControllerRuntimeVersion is the kubernetes-sigs/controller-runtime version to be used in the project.
	ControllerRuntimeVersion = "v0.10.0"
	// LibraryGoVersion is the openshift/library-go version used in the project
	LibraryGoVersion = "v0.0.0-20210914071953-94a0fd1d5849"
	// ControllerToolsVersion is the kubernetes-sigs/controller-tools version to be used in the project
	ControllerToolsVersion = "v0.7.0"
	// KustomizeVersion is the kubernetes-sigs/kustomize version to be used in the project
	KustomizeVersion = "v3.8.7"

	imageName = "controller:latest"
)

var _ plugins.Scaffolder = &initScaffolder{}

type initScaffolder struct {
	config          config.Config
	boilerplatePath string
	license         string
	owner           string

	// fs is the filesystem that will be used by the scaffolder
	fs machinery.Filesystem
}

// NewInitScaffolder returns a new scaffolder for project initialization operations
func NewInitScaffolder(config config.Config, license, owner string) plugins.Scaffolder {
	return &initScaffolder{
		config:          config,
		boilerplatePath: hack.DefaultBoilerplatePath,
		license:         license,
		owner:           owner,
	}
}

// InjectFS implements cmdutil.Scaffolder
func (s *initScaffolder) InjectFS(fs machinery.Filesystem) {
	s.fs = fs
}

// Scaffold implements cmdutil.Scaffolder
func (s *initScaffolder) Scaffold() error {
	fmt.Println("writing scaffold for you to edit...")

	// Initialize machinery.Scaffold that writes boilerplate to disk.
	scaffold := machinery.NewScaffold(s.fs,
		machinery.WithConfig(s.config),
	)

	bpFile := &hack.Boilerplate{
		License: s.license,
		Owner:   s.owner,
	}

	bpFile.Path = s.boilerplatePath
	if err := scaffold.Execute(bpFile); err != nil {
		return err
	}

	boilerplate, err := afero.ReadFile(s.fs.FS, s.boilerplatePath)
	if err != nil {
		return err
	}

	// Initialize machinery.Scaffold to write files to disk
	scaffold = machinery.NewScaffold(s.fs,
		machinery.WithConfig(s.config),
		machinery.WithBoilerplate(string(boilerplate)))

	return scaffold.Execute(
		&templates.GoMod{
			ControllerRuntimeVersion: ControllerRuntimeVersion,
			LibraryGoVersion:         LibraryGoVersion,
		},
		&templates.GitIgnore{},
		&templates.Main{},
		&templates.Makefile{
			Image:                  imageName,
			BoilerplatePath:        s.boilerplatePath,
			ControllerToolsVersion: ControllerToolsVersion,
			KustomizeVersion:       KustomizeVersion,
		},
		&templates.Dockerfile{},
		&templates.DockerIgnore{},
	)
}
