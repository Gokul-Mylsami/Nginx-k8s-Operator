/*
Copyright 2025.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	nginxv1alpha1 "github.com/gokul-mylsami/nginx-operator/api/v1alpha1"
	"github.com/gokul-mylsami/nginx-operator/internal/utils"
)

// NginxUpstreamReconciler reconciles a NginxUpstream object
type NginxUpstreamReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=nginx.gokul-mylsami.com,resources=nginxupstreams,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nginx.gokul-mylsami.com,resources=nginxupstreams/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=nginx.gokul-mylsami.com,resources=nginxupstreams/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NginxUpstream object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *NginxUpstreamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	upstreamRoutes := &nginxv1alpha1.NginxUpstream{}
	err := r.Get(ctx, req.NamespacedName, upstreamRoutes)

	if err != nil {
		log.Log.Error(err, "unable to fetch NginxRoutes")
		return ctrl.Result{}, err
	}

	utils.UpstreamTemplateGenerator(*upstreamRoutes, upstreamRoutes.Spec.TemplateFile)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NginxUpstreamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nginxv1alpha1.NginxUpstream{}).
		Complete(r)
}
