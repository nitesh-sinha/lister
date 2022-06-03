package main

import (
	"time"

	//"context"

	"flag"
	"fmt"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// func main() {
// 	// When running in a client host machine, read K8S API info
// 	// from kubeconfig file
// 	kubeconfig := flag.String("kubeconfig", "/Users/nitesh.sinha/.kube/dev.config",
// 		"Path to your kubeconfig")
// 	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
// 	if err != nil {
// 		fmt.Printf("Error building from flags: %s \n", err.Error())
// 		// Running in K8S cluster as a pod. So, read API keys from
// 		// the serviceaccount files mounted inside the pod
// 		// by kubelet when that pod is launched
// 		config, err = rest.InClusterConfig()
// 		if err != nil {
// 			fmt.Printf("Error read incluster config: %s", err.Error())
// 		}
// 	}
// 	//fmt.Println(config)
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		fmt.Printf("Error creating client set: %s \n", err.Error())
// 	}
// 	ctx := context.Background()
// 	// When running inside a K8S pod, by default this code will not
// 	// be able to list pods or deployments in default namespace
// 	// since it will be running using default serviceaccount credentials
// 	// whcih doesn't have access. Role and Rolebinding has to be explicitly
// 	// created to allow default serviceaccount to list pod and deployment
// 	// kubectl create role pod-depl --resource pods,deployments --verb list
// 	// kubectl create rolebinding pod-depl-binding --role pod-depl --serviceaccount default:default
// 	pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
// 	if err != nil {
// 		fmt.Printf("Error listing pods: %s \n", err.Error())
// 	}
// 	fmt.Println("No. of pods = ", len(pods.Items))
// 	for _, pod := range pods.Items {
// 		fmt.Println(pod.Name)
// 	}
// 	deployments, err := clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
// 	if err != nil {
// 		fmt.Printf("Error listing deployments: %s \n", err.Error())
// 	}
// 	fmt.Println("No. of deployments = ", len(deployments.Items))
// 	for _, dep := range deployments.Items {
// 		fmt.Println(dep.Name)
// 	}
// }

func main() {
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

	informerFactory := informers.NewSharedInformerFactory(clientset, 10*time.Second)
	// Use a podinformer interface to register event handlers
	// so that our client code is informed when new
	// pods are added, updated or deleted
	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(new interface{}) {
			// business logic here
		},
		UpdateFunc: func(old, new interface{}) {
			// business logic here
		},
		DeleteFunc: func(obj interface{}) {
			// some more business logic
		},
	})
	// Start the pod informer and call List to K8S API server
	// to load the in mem store
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
	// Fetch lister pods from in mem store
	// The following is same as `kubectl get pods -l app=lister`
	listerReq, _ := labels.NewRequirement("app", selection.Equals, []string{"lister"})
	selector := labels.NewSelector().Add(*listerReq)
	pods, err := podInformer.Lister().List(selector)
	if err != nil {
		// handle it
	}
	fmt.Println("Pods with label `app=lister` are:")
	for _, pod := range pods {
		fmt.Println(pod.Name)
	}
}
