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

package controllers

import (
	"fmt"
	"path/filepath"

	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
)

var _ machinery.Template = &Controller{}

// Controller scaffolds the file that defines the controller for a CRD or a builtin resource
// nolint:maligned
type Controller struct {
	machinery.TemplateMixin
	machinery.MultiGroupMixin
	machinery.BoilerplateMixin
	machinery.ResourceMixin

	Force bool
}

// SetTemplateDefaults implements file.Template
func (f *Controller) SetTemplateDefaults() error {
	if f.Path == "" {
		if f.MultiGroup && f.Resource.Group != "" {
			f.Path = filepath.Join("controllers", "%[group]", "%[kind]_controller.go")
		} else {
			f.Path = filepath.Join("controllers", "%[kind]_controller.go")
		}
	}
	f.Path = f.Resource.Replacer().Replace(f.Path)
	fmt.Println(f.Path)

	f.TemplateBody = controllerTemplate

	if f.Force {
		f.IfExistsAction = machinery.OverwriteFile
	} else {
		f.IfExistsAction = machinery.Error
	}

	return nil
}

//nolint:lll
const controllerTemplate = `{{ .Boilerplate}}

package {{ if and .MultiGroup .Resource.Group }}{{ .Resource.PackageName }}{{ else }}controllers{{ end }}

import (
	"context"
	"github.com/openshift/library-go/pkg/controller/factory"
)

type {{ .Resource.Kind }}Controller struct {
	// TODO: fill in the controller you would like to configure
	// for the resource.
}

func New{{ .Resource.Kind }}Controller() factory.Controller {
	// TODO: use this to create a new instance of your controller.
	
	return nil
}

// sync contains the logic of the reconciler
func (c *{{ .Resource.Kind }}Controller) sync(ctx context.Context, syncContext factory.SyncContext) error {

	// TODO: implement your reconciler logic here.

	return nil
}
`
