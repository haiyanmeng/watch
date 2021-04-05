package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"net/http"
	_ "net/http/pprof"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	qps := flag.Int("qps", 20, "REST config qps")
	burst := flag.Int("burst", 20, "REST config burst")
	count := flag.Int("count", 100, "watch count")
	flag.Parse()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	fmt.Printf("qps: %d;  burst: %d; count: %d\n", *qps, *burst, *count)

	var config *rest.Config
	var err error
	if *kubeconfig == "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err)
		}
	}
	config.QPS = float32(*qps)
	config.Burst = *burst

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	// listNS(clientset, *count)
	// createNS(clientset, *count)
	// deleteNS(clientset, *count)
	createWatchForNS(clientset, *count)
}

func listNS(clientset *kubernetes.Clientset, count int) {
	start := time.Now()
	for i := 0; i < count; i++ {
		_, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}

		// for _, ns := range nsList.Items {
		// 	fmt.Println(ns.Name)
		// }

	}
	duration := time.Since(start).Seconds()
	fmt.Printf("time to list namespaces %d times: %vs\n", count, duration)
}

func createWatchForNS(clientset *kubernetes.Clientset, count int) {
	start := time.Now()
	options := metav1.ListOptions{
		AllowWatchBookmarks: true,
		Watch:               true,
		// ResourceVersion:     "1182xs1",
	}
	for i := 0; i < count; i++ {
		_, err := clientset.CoreV1().Namespaces().Watch(context.TODO(), options)
		if err != nil {
			fmt.Println(err)
		}
	}
	duration := time.Since(start).Seconds()
	fmt.Printf("time to create %d watches: %vs\n", count, duration)
}

func createNS(clientset *kubernetes.Clientset, count int) {
	start := time.Now()
	for i := 0; i < count; i++ {
		ns := v1.Namespace{}
		ns.Name = fmt.Sprintf("ns%d", i)
		_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), &ns, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
	}
	duration := time.Since(start).Seconds()
	fmt.Printf("time to create %d namespaces: %vs\n", count, duration)
}

func deleteNS(clientset *kubernetes.Clientset, count int) {
	start := time.Now()
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("ns%d", i)
		err := clientset.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			panic(err)
		}
	}
	duration := time.Since(start).Seconds()
	fmt.Printf("time to delete %d namespaces: %vs\n", count, duration)
}
