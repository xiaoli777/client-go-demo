package main

import (
	"fmt"
	"time"

	"client-go-demo/common"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func main() {
	var (
		clientset *kubernetes.Clientset
		epsList *corev1.EndpointsList
		err error
	)

	// 初始化 k8s 客户端
	if clientset, err = common.InitClient(); err != nil {
		fmt.Printf("Init client Failed, the reason is: %+v\n", err)
		return
	}

	for {
		// 通过 k8s API 获取 Endpoints
		if epsList, err = clientset.CoreV1().Endpoints("default").List(metav1.ListOptions{}); err != nil {
			fmt.Printf("Get Endpoint List Failed, the reason is: %+v\n", err)
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