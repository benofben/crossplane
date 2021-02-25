/*
Copyright 2021 The Crossplane Authors.

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

package rbac

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/logging"

	v1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
	pkgv1 "github.com/crossplane/crossplane/apis/pkg/v1"
	"github.com/crossplane/crossplane/internal/initializer"
)

// InitCommand configuration for the initialization of RBAC controllers.
type InitCommand struct {
	Name string
}

// Run starts the initialization process.
func (c *InitCommand) Run(s *runtime.Scheme, log logging.Logger) error {
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return errors.Wrap(err, "Cannot get config")
	}

	cl, err := client.New(cfg, client.Options{Scheme: s})
	if err != nil {
		return errors.Wrap(err, "cannot create new kubernetes client")
	}
	// NOTE(muvaf): The plural form of the kind name is not available in Go code.
	i := initializer.New(cl,
		initializer.NewCRDWaiter([]string{
			fmt.Sprintf("%s.%s", "compositeresourcedefinitions", v1.Group),
			fmt.Sprintf("%s.%s", "providerrevisions", pkgv1.Group),
		}, time.Minute, log),
	)
	if err := i.Init(context.TODO()); err != nil {
		return errors.Wrap(err, "cannot initialize core")
	}
	log.Info("Initialization has been completed")
	return nil
}
