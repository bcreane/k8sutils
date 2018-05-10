package main

import (
	"flag"
	"fmt"
	"github.com/bcreane/k8sutils"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func main() {
	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	var namespace = flag.String("namespace", "default", "namespace")
	var selector = flag.String("selector", "", "pod selector")
	var timeout = flag.Int("timeout", 30, "timeout in seconds")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Block up to timeout seconds for listed pods in namespace/selector to enter running state
	err = k8sutils.WaitForPodBySelectorRunning(clientSet, *namespace, *selector, *timeout)
	if err != nil {
		log.Errorf("\nThe pod never entered running phase\n")
		os.Exit(1)
	}
	fmt.Printf("\nAll pods in namespace=\"%s\" with selector=\"%s\" are running!\n", *namespace, *selector)
}
