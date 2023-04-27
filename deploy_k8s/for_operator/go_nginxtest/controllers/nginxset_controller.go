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
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"nginx.com/nginx-deployer/api/v2beta2"
	ngv2beta2 "nginx.com/nginx-deployer/api/v2beta2"
)

// NginxsetReconciler reconciles a Nginxset object
type NginxsetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ng.nginx.com,resources=nginxsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ng.nginx.com,resources=nginxsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ng.nginx.com,resources=nginxsets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Nginxset object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *NginxsetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	fmt.Println("========================================================")
	_ = log.FromContext(ctx)
	// TODO(user): your logic here

	// 取得 Target 既 K8s 資源
	nSet := v2beta2.Nginxset{}
	err := r.Get(ctx, req.NamespacedName, &nSet)
	if err != nil {
		fmt.Println(fmt.Errorf("error: there is no target resource %v", err))
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 600}, nil
	}

	fmt.Println(req.Name)
	fmt.Println(req.Namespace)
	fmt.Println(req.NamespacedName)
	fmt.Println("---")
	fmt.Println(nSet.Spec.ReturnText)
	fmt.Println(nSet.Status.Ready)
	fmt.Println(nSet.Status.LinkedDeployment)

	// 取得 CR 資源是否有綁定有效的 deployment

	targetUid := nSet.UID

	bindedDeployment := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-deployment", req.Name),
			Namespace: req.Namespace,
			UID:       targetUid,
		},
	}

	// 如果 CR 实例不存在，则根据 CRD 创建
	// 否則 CR 实例存在，则将 Annotations 中记录的 Spec 值与当前的 Spec 比较

	existed := r.Get(ctx, req.NamespacedName, &bindedDeployment)
	if existed != nil {

		fmt.Println(fmt.Println("There is not bindedDeployment"))

		{ // Create Configmap and Deployment

			replcas := int32(1)

			labelSet := map[string]string{
				"tester": "aco",
			}

			configmap := corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ConfigMap",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-configmap", req.Name),
					Namespace: req.Namespace,
				},
				Data: map[string]string{
					"index.html": nSet.Spec.ReturnText,
				},
			}

			err = r.Client.Create(ctx, &configmap)
			if err != nil {
				return ctrl.Result{
					Requeue:      true,
					RequeueAfter: time.Second * 600,
				}, nil
			}

			deployment := appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-deployment", req.Name),
					Namespace: req.Namespace,
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &replcas,
					Selector: &(metav1.LabelSelector{
						MatchLabels: labelSet,
					}),
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: labelSet,
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  fmt.Sprintf("%s-container", req.Name),
									Image: "nginx:alpine",
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "firstmount",
											MountPath: "/usr/share/nginx/html",
										},
									},
								},
							},
							Volumes: []corev1.Volume{
								{
									Name: "firstmount",
									VolumeSource: corev1.VolumeSource{
										ConfigMap: &corev1.ConfigMapVolumeSource{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: fmt.Sprintf("%s-configmap", req.Name),
											},
										},
									},
								},
							},
						},
					},
				},
			}

			err = r.Client.Create(ctx, &deployment)
			if err != nil {
				return ctrl.Result{
					Requeue:      true,
					RequeueAfter: time.Second * 600,
				}, nil
			}

			fmt.Println("OKOK")
		}

	} else {

		// 找到 target 的 binding deployment
		fmt.Println(fmt.Println("Target bindedDeployment found"))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NginxsetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ngv2beta2.Nginxset{}).
		Complete(r)
}
