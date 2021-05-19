package common

import (
	"io/ioutil"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func InitClient() (clientset *kubernetes.Clientset, err error) {
	var (
		restConf *rest.Config
	)

	if restConf, err = GetRestConf(); err != nil {
		return nil, err
	}

	// 生成 k8s client
	if clientset, err = kubernetes.NewForConfig(restConf); err != nil {
		return nil, err
	}

	return clientset, nil
}

func GetRestConf() (restConf *rest.Config, err error) {
	var (
		kubeconfig []byte
	)

	// 读取 kubeconfig 文件
	if kubeconfig, err = ioutil.ReadFile("./k8s.conf"); err != nil {
		return nil, err
	}

	// 通过 kubeconfig 生成 client
	if restConf, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig); err != nil {
		return nil, err
	}

	return restConf, nil
}