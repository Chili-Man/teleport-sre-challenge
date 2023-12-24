package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	err := run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
	}
}

func run(args []string) error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	var port, kubeconfig string
	flag.StringVar(&port, "port", "8080", "server port")
	flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(homedir, ".kube", "config"), "path to the kubeconfig file")
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// Informer for Deployments

	/* Handlers */

	// add handlers
	http.Handle("/healthz", &healthzHandler{clientset: clientset})

	http.Handle("/deployments", &deploymentsHandler{client: clientset})
	http.Handle("/deployments/", &deploymentsHandler{client: clientset})

	fmt.Println("Starting server at :" + port)
	return http.ListenAndServe(":"+port, nil)
}
