package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	var delete bool
	flag.BoolVar(&delete, "delete", false, "Executes delete operations. Default false.")
	kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "absolute path to the kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Println("k8s build config err:", err, *kubeconfig)
		os.Exit(1)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("k8s create client err:", err)
		os.Exit(1)

	}
	fmt.Printf("1. connected to a k8s cluster\n\n")

	informerFactory := informers.NewSharedInformerFactory(kubeClient, 1*time.Minute)
	podInformer := informerFactory.Core().V1().Pods().Informer()
	namespaceInformer := informerFactory.Core().V1().Namespaces().Informer()
	stopCh := make(chan struct{})
	defer close(stopCh)
	informerFactory.Start(stopCh)

	if ok := cache.WaitForCacheSync(stopCh, podInformer.HasSynced); !ok {
		fmt.Println("pod Informer err syncing")
		os.Exit(1)
	}
	if ok := cache.WaitForCacheSync(stopCh, namespaceInformer.HasSynced); !ok {
		fmt.Println("namespace Informer err syncing")
		os.Exit(1)
	}

	root := newProgram(kubeClient,
		informerFactory.Core().V1().Namespaces().Lister(),
		informerFactory.Core().V1().Pods().Lister(),
	)
	root.run(context.Background(), delete)
}
