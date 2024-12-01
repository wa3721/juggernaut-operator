/*
Copyright 2024.

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
	"fmt"
	k8s "juggernaut/pkg"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1 "juggernaut/api/v1"
)

// JuggernautReconciler reconciles a Juggernaut object
type JuggernautReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=operator.oceanoperator.com,resources=juggernauts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operator.oceanoperator.com,resources=juggernauts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=operator.oceanoperator.com,resources=juggernauts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Juggernaut object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *JuggernautReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling Juggernaut")
	// 创建一个实例
	juggernaut := &operatorv1.Juggernaut{}
	err := r.Get(ctx, req.NamespacedName, juggernaut)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			logger.Info("juggernaut resource not found. Ignoring")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Juggernaut")
		return ctrl.Result{}, err
	}
	//写逻辑
	//我需要一个deployment对象，一个service对象，一个configmap对象
	//根据juggernaut对象，如果这个对象不存在，则skip，如果存在，则根据字段生成以上三个对象
	//如果对象的metadata.annotations.存在skip-sync: "true"则不做调谐
	return ctrl.Result{}, nil
}

func (r *JuggernautReconciler) reconcileJuggernaut(ctx context.Context, juggernaut *operatorv1.Juggernaut) error {
	if err := r.reconcileDeployment(ctx, juggernaut); err != nil {
		return err
	}

	if err := r.reconcileService(ctx, juggernaut); err != nil {
		return err
	}

	if err := r.reconcileConfigmap(ctx, juggernaut); err != nil {
		return err
	}

	return nil
}

// 调谐deployment
func (r *JuggernautReconciler) reconcileDeployment(ctx context.Context, juggernaut *operatorv1.Juggernaut) error {
	newDeploy, err := k8s.NewDeployment(juggernaut)
	if err != nil {
		return fmt.Errorf("failed to build Deployment from Nginx: %w", err)
	}

	var currentDeploy appv1.Deployment
	err = r.Client.Get(ctx, types.NamespacedName{Name: newDeploy.Name, Namespace: newDeploy.Namespace}, &currentDeploy)
	if errors.IsNotFound(err) {
		return r.Client.Create(ctx, newDeploy)
	}

	if err != nil {
		return fmt.Errorf("failed to retrieve Deployment: %w", err)
	}

	currentDeploy.Spec = newDeploy.Spec
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JuggernautReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1.Juggernaut{}).
		Named("juggernaut").
		Complete(r)
}
