package main

import (
	"flag"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func getGroupVersionResource() {
	var res string

	flag.StringVar(&res, "resource", "", "K8S resource name")
	flag.Parse()

	configFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	matchVersionFlags := cmdutil.NewMatchVersionFlags(configFlags)
	restMapper, err := cmdutil.NewFactory(matchVersionFlags).ToRESTMapper()
	if err != nil {
		fmt.Println("Error building REST mapper: ", err.Error())
		return
	}
	groupVerRes, err := restMapper.ResourceFor(schema.GroupVersionResource{
		Resource: res,
	})
	if err != nil {
		fmt.Println("Error getting Group, Version and Resource: ", err.Error())
		return
	}
	fmt.Printf("Group= %s, Version=%s, Resource=%s \n", groupVerRes.Group, groupVerRes.Version,
		groupVerRes.Resource)

}
