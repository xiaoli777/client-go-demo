package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"client-go-demo/common"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yaml2 "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

func main() {
	var (
		clientset *kubernetes.Clientset
		dpYaml []byte
		dpJson []byte
		dp = &v1.Deployment{}
		replicas int32
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

	// 修改结构体中 replicas 的数量
	replicas = 1
	dp.Spec.Replicas = &replicas

	// 查询 k8s 是否已有该 deployment
	if _, err = clientset.AppsV1().Deployments("default").Get(dp.Name, metav1.GetOptions{}); err != nil {

		if !errors.IsNotFound(err) {
			goto FAILED
		}
		// 不存在则创建
		if _, err = clientset.AppsV1().Deployments("default").Create(dp); err != nil {
			goto FAILED
		}
	} else {	// 已存在则更新
		if _, err = clientset.AppsV1().Deployments("default").Update(dp); err != nil {
			goto FAILED
		}
	}

	fmt.Printf("Successfully!\n")
	return

FAILED:
	fmt.Printf("Failed, the reson is %+v\n", err)
	return
}
