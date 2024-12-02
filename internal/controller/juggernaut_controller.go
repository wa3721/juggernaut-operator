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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"

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
	if err := r.reconcileJuggernaut(ctx, juggernaut); err != nil {
		logger.Error(err, "Fail to reconcile")
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: time.Second * 5}, nil
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
	if err != nil {
		if errors.IsNotFound(err) {
			return r.Client.Create(ctx, newDeploy)
		}
		return fmt.Errorf("failed to retrieve Deployment: %w", err)
	}
	if currentDeploy.ObjectMeta.Annotations["skip-sync"] == "true" {
		klog.Infof("deployment skip sync")
		return nil
	}
	//校验currentDeploy newDeploy spec的差别 通过Client.Patch修改currentDeploy
	// 比较newDeploy和currentDeploy的spec部分
	if !equality.Semantic.DeepEqual(newDeploy.Spec, currentDeploy.Spec) {
		// 如果有差异，则通过Client.Patch更新currentDeploy
		patch := client.StrategicMergeFrom(currentDeploy.DeepCopy())
		currentDeploy.Spec = newDeploy.Spec
		if err := r.Client.Patch(ctx, &currentDeploy, patch); err != nil {
			return fmt.Errorf("failed to patch Deployment: %w", err)
		}
	}
	return nil
}

// 调谐service
func (r *JuggernautReconciler) reconcileService(ctx context.Context, juggernaut *operatorv1.Juggernaut) error {
	newService := k8s.NewService(juggernaut)
	var currentService corev1.Service
	err := r.Client.Get(ctx, types.NamespacedName{Name: newService.Name, Namespace: newService.Namespace}, &currentService)
	if errors.IsNotFound(err) {
		return r.Client.Create(ctx, newService)
	}

	if err != nil {
		return fmt.Errorf("failed to retrieve service: %w", err)
	}
	if currentService.ObjectMeta.Annotations["skip-sync"] == "true" {
		klog.Infof("service skip sync")
		return nil
	}
	if !equality.Semantic.DeepEqual(newService.Spec, currentService.Spec) {
		// 如果有差异，则通过Client.Patch更新currentDeploy
		patch := client.StrategicMergeFrom(currentService.DeepCopy())
		currentService.Spec = newService.Spec
		if err := r.Client.Patch(ctx, &currentService, patch); err != nil {
			return fmt.Errorf("failed to patch service: %w", err)
		}
	}

	return nil
}

// 调谐configmap
func (r *JuggernautReconciler) reconcileConfigmap(ctx context.Context, juggernaut *operatorv1.Juggernaut) error {
	//默认的configmap
	//根据juggernaut名字+configmap生成//挂载到deployment
	newConfigmap := k8s.NewConfigmap(juggernaut)
	var currentConfigmap corev1.ConfigMap
	err := r.Client.Get(ctx, types.NamespacedName{Name: newConfigmap.Name, Namespace: newConfigmap.Namespace}, &currentConfigmap)
	if errors.IsNotFound(err) {
		return r.Client.Create(ctx, newConfigmap)
	}
	if err != nil {
		return fmt.Errorf("failed to retrieve configmap: %w", err)
	}
	if currentConfigmap.ObjectMeta.Annotations["skip-sync"] == "true" {
		klog.Infof("configmap skip sync")
		return nil
	}
	if !equality.Semantic.DeepEqual(newConfigmap.Data, currentConfigmap.Data) {
		// 如果有差异，则通过Client.Patch更新currentDeploy
		patch := client.StrategicMergeFrom(currentConfigmap.DeepCopy())
		currentConfigmap.Data = newConfigmap.Data
		if err := r.Client.Patch(ctx, &currentConfigmap, patch); err != nil {
			return fmt.Errorf("failed to patch configmap: %w", err)
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JuggernautReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1.Juggernaut{}).
		Named("juggernaut").
		Complete(r)
}
