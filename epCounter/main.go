package main

import (
	"fmt"
	"time"

	"client-go-demo/common"

	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func main() {
	var (
		clientset *kubernetes.Clientset
		epsList *core_v1.EndpointsList
		err error
	)

	// 初始化 k8s 客户端
	if clientset, err = common.InitClient(); err != nil {
		return
	}

	for {
		if epsList, err = clientset.CoreV1().Endpoints("default").List(meta_v1.ListOptions{}); err != nil {
			return
		}
		fmt.Printf("There are %d endpoints in the cluster\n", len(epsList.Items))
		for i := 0; i < len(epsList.Items); i++ {
			fmt.Println(epsList.Items[i].ObjectMeta.Name)
			for j :=0; j < len(epsList.Items[i].Subsets); j++ {
				fmt.Println("Addresses", epsList.Items[i].Subsets[j].Addresses)
				fmt.Println("NotReadyAddresses", epsList.Items[i].Subsets[j].NotReadyAddresses)
			}
		}
		fmt.Printf("\n")
		time.Sleep(10 * time.Second)
	}

	return
}