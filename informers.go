package main

import (
	"flag"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func fetchPodsUsingInformer() {
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
