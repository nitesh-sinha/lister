package main

import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func fetchPodsUsingClientSet() {
	// When running in a client host machine, read K8S API info
	// from kubeconfig file
	kubeconfig := flag.String("kubeconfig", "/Users/nitesh.sinha/.kube/config",
		"Path to your kubeconfig")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("Error building from flags: %s \n", err.Error())
		// Running in K8S cluster as a pod. So, read API keys from
		// the serviceaccount files mounted inside the pod
		// by kubelet when that pod is launched
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("Error read incluster config: %s", err.Error())
		}
	}
	//fmt.Println(config)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating client set: %s \n", err.Error())
	}
	ctx := context.Background()
	// When running inside a K8S pod, by default this code will not
	// be able to list pods or deployments in default namespace
	// since it will be running using default serviceaccount credentials
	// whcih doesn't have access. Role and Rolebinding has to be explicitly
	// created to allow default serviceaccount to list pod and deployment
	// kubectl create role pod-depl --resource pods,deployments --verb list
	// kubectl create rolebinding pod-depl-binding --role pod-depl --serviceaccount default:default
	pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing pods: %s \n", err.Error())
	}
	fmt.Println("No. of pods = ", len(pods.Items))
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}
	deployments, err := clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing deployments: %s \n", err.Error())
	}
	fmt.Println("No. of deployments = ", len(deployments.Items))
	for _, dep := range deployments.Items {
		fmt.Println(dep.Name)
	}
}
