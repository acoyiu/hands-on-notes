/*
Copyright 2022.

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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	acov1alpha1 "test.com/cm-deployer/api/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/builder"   // Required for Watching
	"sigs.k8s.io/controller-runtime/pkg/handler"   // Required for Watching
	"sigs.k8s.io/controller-runtime/pkg/predicate" // Required for Watching
	"sigs.k8s.io/controller-runtime/pkg/source"    // Required for Watching
)

// WebsetReconciler reconciles a Webset object
type WebsetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var websetlogger = log.Log.WithName("Webset-res")

//+kubebuilder:rbac:groups=aco.test.com,resources=websets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aco.test.com,resources=websets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aco.test.com,resources=websets/finalizers,verbs=update

// below are custom added ====================================================================================
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// below are custom added ====================================================================================

func (r *WebsetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	websetlogger.Info("Webset Reconcile started...")

	{
		// Get target CRD
		var websetDeploy acov1alpha1.Webset
		if err := r.Get(ctx, req.NamespacedName, &websetDeploy); err != nil {
			// we'll ignore not-found errors, since they can't be fixed by an immediate requeue
			// we'll need to wait for a new notification, and we can get them on deleted requests.
			websetlogger.Error(err, "unable to fetch websetDeploy")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		fmt.Println("res name: ", websetDeploy.GetName())
		fmt.Println("res namespace: ", websetDeploy.GetNamespace())

		foundConfigMap := &corev1.ConfigMap{}
		var configMapVersion string

		if websetDeploy.Spec.ConfigMap != "" {

			// get starget configmap
			configMapName := websetDeploy.Spec.ConfigMap
			err := r.Get(ctx, types.NamespacedName{Name: configMapName, Namespace: websetDeploy.Namespace}, foundConfigMap)
			if err != nil {
				// If a configMap name is provided, then it must exist
				// You will likely want to create an Event for the user to understand why their reconcile is failing.
				return ctrl.Result{}, err
			}

			// Hash the data in some way, or just use the version of the Object
			configMapVersion = foundConfigMap.ResourceVersion
			fmt.Println("configMapVersion:", configMapVersion)
		}
	}

	websetlogger.Info("Webset Reconcile Ended")

	return ctrl.Result{}, nil
}

func createNewPod(ctx context.Context, req ctrl.Request, r *WebsetReconciler) (err error) {
	newPod := v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  fmt.Sprintf("%s-container", req.Name),
					Image: "nginx:alpine",
				},
			},
		},
	}

	newPod.Name = fmt.Sprintf("%s-pod", req.Name)
	newPod.Namespace = req.Namespace

	err = r.Create(ctx, &newPod)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		return
	}

	return
}

// ========================================================================================================

// SetupWithManager sets up the controller with the Manager.
func (r *WebsetReconciler) SetupWithManager(mgr ctrl.Manager) error {

	var configMapField = ".spec.configmap"

	if err :=
		mgr.GetFieldIndexer().
			IndexField(
				context.Background(),
				&acov1alpha1.Webset{},
				configMapField, // Index 名
				func(rawObj client.Object) []string {
					// 此函數在 Controller 開始時運行，生成一個反向索引（Index）
					// 以告知“每一個”資源用什麽 “index” 可以尋找的得到
					// 例如: nameToListen := []string{"cm1", "cm2"}
					// 意思就是 “fields.OneTermEqualSelector” 找 “Index名” 這個 index 的時候
					// 可以找到是與哪一個 resource 有關聯
					configDeployment := rawObj.(*acov1alpha1.Webset)
					if configDeployment.Spec.ConfigMap == "" {
						return nil
					}
					nameToListen := []string{configDeployment.Spec.ConfigMap}
					fmt.Println(
						"IndexField: -------",
						configDeployment.GetNamespace(),
						configDeployment.GetName(),
						nameToListen,
					)
					return nameToListen
				},
			); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&acov1alpha1.Webset{}).
		Owns(&appsv1.Deployment{}).
		Watches(
			// 指定要 listen 的資源類型
			&source.Kind{Type: &corev1.ConfigMap{}},
			// configMap = currently updated target obj
			handler.EnqueueRequestsFromMapFunc(func(configMap client.Object) []reconcile.Request {
				// 用（反向）index 取回目標的 resource 字段
				fmt.Println(fields.OneTermEqualSelector(configMapField, configMap.GetName()))
				// 根據 currently updated resource (configmap), 聲明 CR 清單, 聲明筛选选项
				attachedConfigDeployments := acov1alpha1.WebsetList{}
				listOps := &client.ListOptions{
					FieldSelector: fields.OneTermEqualSelector(configMapField, configMap.GetName()),
					Namespace:     configMap.GetNamespace(),
				}
				// kubectl 獲取 CR 列表
				err := r.List(context.TODO(), &attachedConfigDeployments, listOps)
				if err != nil {
					return []reconcile.Request{}
				}
				// 根據 configmap 資料，聲明 CR slice
				requests := make([]reconcile.Request, len(attachedConfigDeployments.Items))
				// 每一個 query 了出來的 configmap 都執行一次 reconcile call（request）
				for i, item := range attachedConfigDeployments.Items {
					requests[i] = reconcile.Request{
						NamespacedName: types.NamespacedName{
							Name:      item.GetName(),
							Namespace: item.GetNamespace(),
						},
					}
				}
				fmt.Println("requests", requests)
				return requests
			}),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Complete(r)
}
