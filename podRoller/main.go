package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"client-go-demo/common"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yaml2 "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

func main(){
	var (
		clientset *kubernetes.Clientset
		dpYaml []byte
		dpJson []byte
		dp = &v1.Deployment{}
		k8sDP *v1.Deployment
		containers []corev1.Container
		tomcatContainer corev1.Container
		containerPort []corev1.ContainerPort
		tomcatPort corev1.ContainerPort
		replicas int32
		podList *corev1.PodList
		pod corev1.Pod
		err error
	)

	// 初始化 k8s 客户端
	if clientset, err = common.InitClient(); err != nil {
		goto FAILED
	}

	// 读取 yaml 文件
	if dpYaml, err = ioutil.ReadFile("./nginx-deployment.yaml"); err != nil {
		goto FAILED
	}

	// 转换成 json 格式
	if dpJson, err = yaml2.ToJSON(dpYaml); err != nil {
		goto FAILED
	}

	// 解析 json 格式, 并存储在结构体 dp 中
	if err = json.Unmarshal(dpJson, dp); err != nil {
		goto FAILED
	}

	// 更新 label, image,
	dp.Spec.Template.Labels["image"] = "tomcat"
	dp.Spec.Template.Labels["time"] = strconv.Itoa(int(time.Now().Unix()))

	tomcatContainer.Name = "tomcat"
	tomcatContainer.Image = "tomcat:latest"
	containers = append(containers, tomcatContainer)
	dp.Spec.Template.Spec.Containers = containers

	for _, c := range dp.Spec.Template.Spec.Containers {
		tomcatPort.ContainerPort = 8080
		containerPort = append(containerPort, tomcatPort)
		c.Ports = containerPort
	}

	replicas = 5
	dp.Spec.Replicas = &replicas

	// 更新 deployment
	if _, err = clientset.AppsV1().Deployments("default").Update(dp); err != nil {
		goto FAILED
	}

	for {
		// 获取 k8s 中的 deployment
		if k8sDP, err = clientset.AppsV1().Deployments("default").Get(dp.Name, metav1.GetOptions{}); err != nil {
			goto RETRY
		}

		// 进行滚动升级是否完成
		if k8sDP.Status.Replicas == *(k8sDP.Spec.Replicas) &&
			k8sDP.Status.ReadyReplicas == *(k8sDP.Spec.Replicas) &&
			k8sDP.Status.UpdatedReplicas == *(k8sDP.Spec.Replicas) &&
			k8sDP.Status.AvailableReplicas == *(k8sDP.Spec.Replicas) &&
			k8sDP.Status.ObservedGeneration == k8sDP.Generation {
			break
		}

		// 输出 deployment 的信息
		fmt.Printf("K8s Deployment Information:\n")
		fmt.Printf("Status Replicas: %d\n", k8sDP.Status.Replicas)
		fmt.Printf("Status ReadyReplicas: %d\n", k8sDP.Status.ReadyReplicas)
		fmt.Printf("Status UpdatedReplicas: %d\n", k8sDP.Status.UpdatedReplicas)
		fmt.Printf("Status AvailableReplicas: %d\n", k8sDP.Status.AvailableReplicas)
		fmt.Printf("Status ObservedGeneration: %d\n", k8sDP.Status.ObservedGeneration)
		fmt.Printf("Spec Replicas: %d\n", *(k8sDP.Spec.Replicas))
		fmt.Printf("Spec Generation: %d\n", k8sDP.Generation)

	RETRY:
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Successful Roller Updating!")

	// 输出 Pod 状态
	if podList, err = clientset.CoreV1().Pods("default").List(metav1.ListOptions{LabelSelector: "image=tomcat"}); err == nil {
		for _, pod = range podList.Items {
			podName := pod.Name
			podStatus := string(pod.Status.Phase)

			if podStatus == string(corev1.PodRunning) {
				// why the pod is in this state
				if pod.Status.Reason != "" {
					podStatus = pod.Status.Reason
					goto KO
				}

				// Current service state of pod
				for _, cond := range pod.Status.Conditions {
					if cond.Type == corev1.PodReady {
						if cond.Status != corev1.ConditionTrue {
							podStatus = cond.Reason
						}
						goto KO
					}
				}
			}
		KO:
			fmt.Printf("podName: %s, and its Status: %s.\n", podName, podStatus)
		}
	}

	return

FAILED:
	fmt.Printf("Failed, the reson is %+v\n", err)
	return
}
