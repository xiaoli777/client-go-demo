package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"client-go-demo/common"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	yaml2 "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

func main() {
	var (
		clientset *kubernetes.Clientset
		dpYaml []byte
		dpJson []byte
		dp = &v1.Deployment{}
		containers []corev1.Container
		tomcatContainer corev1.Container
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

	// 解析 json 格式， 并存储在结构体 dp 中
	if err = json.Unmarshal(dpJson, dp); err != nil {
		goto FAILED
	}

	// 修改结构体中 container 信息
	tomcatContainer.Name = "tomcat"
	tomcatContainer.Image = "tomcat:latest"
	containers = append(containers, tomcatContainer)

	dp.Spec.Template.Spec.Containers = containers

	// 更新 deployment, 替换容器的 image
	if _, err = clientset.AppsV1().Deployments("default").Update(dp); err != nil {
		goto FAILED
	}

	fmt.Printf("Successfully!\n")
	return

FAILED:
	fmt.Printf("Failed, the reson is %+v\n", err)
	return
}
