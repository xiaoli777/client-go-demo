package main

import (
	"client-go-demo/common"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	var (
		clientset *kubernetes.Clientset
		req *rest.Request
		res rest.Result
		podList *corev1.PodList
		logs []byte
		err error
	)

	if clientset, err = common.InitClient(); err != nil {
		goto FAILED
	}

	if podList, err = clientset.CoreV1().Pods("default").List(metav1.ListOptions{LabelSelector: "app=nginx"}); err != nil {
		goto FAILED
	}

	if len(podList.Items) == 0 {
		fmt.Printf("There is no nginx pod on K8s.\n")
		return
	}

	req = clientset.CoreV1().Pods("default").GetLogs(podList.Items[0].Name, &corev1.PodLogOptions{Container: "nginx"})

	if res = req.Do(); res.Error() != nil {
		err = res.Error()
		goto FAILED
	}

	if logs, err = res.Raw(); err != nil {
		goto FAILED
	}

	fmt.Printf("%s Logs: \n", podList.Items[0].Name)
	fmt.Println(string(logs))
	fmt.Println("Finished!")
	return

FAILED:
	fmt.Printf("Failed, the reson is %+v\n", err)
	return
}
