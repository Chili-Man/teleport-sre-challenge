package main

import (
	"fmt"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// healthzHandler is an HTTP handler for the healthz API.
type healthzHandler struct {
	clientset *kubernetes.Clientset
}

// ServeHTTP implements http.Handler
func (h *healthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pods, err := h.clientset.CoreV1().Pods(corev1.NamespaceAll).List(r.Context(), metav1.ListOptions{})

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	fmt.Fprintf(w, "healthz: OK\nThere are %d pods in the cluster\n", len(pods.Items))
}
