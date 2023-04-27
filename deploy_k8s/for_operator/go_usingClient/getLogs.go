package main

import (
	"bytes"
	"context"
	"fmt"
	"io"

	v1 "k8s.io/api/core/v1"
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

	// designate namespace
	namespace := "test"

	ctx := context.TODO()

	// get pods list
	podList, err :=
		clientset.
			CoreV1().
			Pods(namespace).
			List(
				ctx,
				metav1.ListOptions{
					LabelSelector: "app=api",
				},
			)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("======================")
	fmt.Println(len(podList.Items))
	for _, item := range podList.Items {
		fmt.Println(item.Name)
	}
	fmt.Println("======================")

	// find what pod to extract logs
	podName := "nginx"

	req :=
		clientset.
			CoreV1().
			Pods(namespace).
			GetLogs(
				podName,
				&(v1.PodLogOptions{
					Container: "nginx",
				}),
			)
	podLogs, err := req.Stream(ctx)

	if err != nil {
		panic(err.Error())
	}

	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		panic(err.Error())
	}
	str := buf.String()

	fmt.Println(str)
}
