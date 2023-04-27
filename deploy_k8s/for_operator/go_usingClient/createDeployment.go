package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Get specific Deployment
	deployName := "nginx"
	namespace := "default"
	targetDeployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		panic("Deployment not found.")
	} else {
		fmt.Println(targetDeployment.ObjectMeta.Name)
	}

	// List Deployment

	clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "app=nginx"})

	// Create Deployment
	replcas := int32(1)
	labelSet := map[string]string{"tester": "aco"}
	clientset.
		AppsV1().
		Deployments(namespace).
		Create(
			context.TODO(),
			&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deploy-02",
					Namespace: "default",
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
									Name:  fmt.Sprintf("%s-container", "deploy-02"),
									Image: "nginx:alpine",
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "firstmount",
											MountPath: "/usr/share/nginx/html",
										},
									},
								},
							},
							// Volumes: []corev1.Volume{
							// 	{
							// 		Name: "firstmount",
							// 		VolumeSource: corev1.VolumeSource{
							// 			ConfigMap: &corev1.ConfigMapVolumeSource{
							// 				LocalObjectReference: corev1.LocalObjectReference{
							// 					Name: fmt.Sprintf("%s-configmap", "deploy-02"),
							// 				},
							// 			},
							// 		},
							// 	},
							// },
						},
					},
				},
			},
			metav1.CreateOptions{},
		)
}
