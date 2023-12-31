/*
Copyright 2023.

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
	discoveryv1beta1 "k8s.io/api/discovery/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	discovery "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// EndpointSliceReconciler reconciles a EndpointSlice object
type EndpointSliceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=networking.envoy.com,resources=EndpointSlices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.envoy.com,resources=EndpointSlices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=networking.envoy.com,resources=EndpointSlices/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EndpointSlice object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *EndpointSliceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var obj discovery.EndpointSlice
	if err := r.Get(ctx, req.NamespacedName, &obj); err != nil {
		log.Error(err, "unable to fetch EndpointSlice")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("found ***")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EndpointSliceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	_, err := mgr.GetClient().RESTMapper().RESTMapping(schema.GroupKind{
		Group: discovery.GroupName,
		Kind:  "EndpointSlice",
	}, "v1")
	if err == nil {
		return ctrl.NewControllerManagedBy(mgr).
			For(&discovery.EndpointSlice{}).
			Complete(r)
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&discoveryv1beta1.EndpointSlice{}).
		Complete(r)
}
