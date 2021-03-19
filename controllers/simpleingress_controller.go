/*


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
	"net/http/httputil"
	"net/url"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	networkingv1 "simpleingress/api/v1"
)

// SimpleIngressReconciler reconciles a SimpleIngress object
type SimpleIngressReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=networking.gravityloop.io,resources=simpleingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.gravityloop.io,resources=simpleingresses/status,verbs=get;update;patch

func (r *SimpleIngressReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("simpleingress", req.NamespacedName)

	inst := &networkingv1.SimpleIngress{}
	err := r.Get(ctx, req.NamespacedName, inst)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "error reading the object")
		return ctrl.Result{}, err
	}

	// if the status does not match the spec, then set status fields according to spec
	if inst.Status.Host == "" || inst.Status.Host != inst.Spec.Host {
		inst.Status.Host = inst.Spec.Host
	} else if inst.Status.ServiceName == "" || inst.Status.ServiceName != inst.Spec.ServiceName {
		inst.Status.ServiceName = inst.Spec.ServiceName
	}

	//convert inst.Spec.ServiceName to *url.URL
	backendTarget := &url.URL{
		Host: inst.Spec.ServiceName,
	}
	//config reverse proxy here by passing in backendTarget to NewSingleHostReverseProxy
	inst.Proxy = httputil.NewSingleHostReverseProxy(backendTarget)
	





	return ctrl.Result{}, nil
}

func (r *SimpleIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1.SimpleIngress{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
