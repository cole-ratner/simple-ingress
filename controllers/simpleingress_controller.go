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

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	networkingv1 "simpleingress/api/v1"
	"simpleingress/internal/pod"
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

	if inst.Status.Phase == "" {
		inst.Status.Phase = networkingv1.PhasePending
	}

	switch inst.Status.Phase {

	case networkingv1.PhasePending:
		log.Info("State: PENDING")

		pod := pod.New(inst)
		err := ctrl.SetControllerReference(inst, pod, r.Scheme)
		if err != nil {
			// requeue with error
			return ctrl.Result{}, err
		}

		query := &corev1.Pod{}
		err = r.Get(context.TODO(), req.NamespacedName, query)
		if err != nil && errors.IsNotFound(err) {
			err = r.Create(context.TODO(), pod)
			if err != nil {
				return ctrl.Result{}, err
			}
			log.Info("pod successfully created", "name", pod.Name)
			return ctrl.Result{}, nil
		} else if err != nil {
			// requeue with err
			log.Error(err, "cannot create pod")
			return ctrl.Result{}, err
		} else if query.Status.Phase == corev1.PodFailed || query.Status.Phase == corev1.PodSucceeded {
			// pod already finished or errored out`
			log.Info("container terminated", "reason", query.Status.Reason, "message", query.Status.Message)
			inst.Status.Phase = networkingv1.PhaseReady
		} else {
			// don't requeue, it will happen automatically when pod status changes
			return ctrl.Result{}, nil
		}
		// POST request to pod /routeMap {"backend": inst.Spec.ServiceName}
		// if 20

	}

	return ctrl.Result{}, nil
}

func (r *SimpleIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1.SimpleIngress{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
