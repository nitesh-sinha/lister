package main

import "fmt"

func main() {
	//fmt.Println("Fetching pods and deployments using basic clientset: ")
	//fetchPodsUsingClientSet()
	// fmt.Println("Fetching pods using Informer - More efficient")
	// fetchPodsUsingInformer()
	fmt.Println("Fetching group version resource")
	getGroupVersionResource()
}
