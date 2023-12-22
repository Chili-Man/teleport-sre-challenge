package main

import (
	"encoding/json"
	//	"fmt"
	"net/http"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/informers"
)

// namespaceFromPath returns the namespace from the provided path parameters.
// There are two cases to consider:
// 1. When the path parameters length is 1, then we assume that no namespace is provided
// 2. When there are 2 or more path parameters, then the second item is assumed to be the namespace
func namespaceFromPath(pathParameters []string) string {
	if params := len(pathParameters); params <= 1 {
		return corev1.NamespaceAll
	} else {
		return pathParameters[1]
	}
}

// healthzHandler is an HTTP handler for the healthz API.
type deploymentsHandler struct {
	client *kubernetes.Clientset
}

// Parses the URL parameters and routes to the correct handler if need be
func (h *deploymentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Need to figure out the path parameters if any
	urlPath := r.URL.EscapedPath()
	pathParameters := strings.FieldsFunc(urlPath, func(char rune) bool {
		return char == '/'
	})

	// Set JSON as response type
	w.Header().Set("Content-Type", "application/json")

	lister := deploymentsLister{
		deploymentsHandler: h,
		Namespace:          namespaceFromPath(pathParameters),
	}

	if err := listDeployments(w, r, lister); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
}

type deploymentsLister struct {
	*deploymentsHandler

	// Namespace of the deployments to list. If left blank, then all namespaces
	// are considered.
	Namespace string
}

// deploymentListerResponse for returning the results of listing the
// deployments. It maps the namespace to deployment names.
type deploymentListerResponse struct {
	Results map[string][]string
}

// listDeployments will list all of the deployments in the provided namespaces.
func listDeployments(w http.ResponseWriter, r *http.Request, d deploymentsLister) error {
	allDeployments, err := d.client.AppsV1().Deployments(d.Namespace).List(r.Context(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	// Generate the response
	response := deploymentListerResponse{
		Results: make(map[string][]string),
	}

	for _, deployment := range allDeployments.Items {
		_, ok := response.Results[deployment.ObjectMeta.Namespace]

		// If the namespace hasn't been added, we need to first allocate the list
		if !ok {
			response.Results[deployment.ObjectMeta.Namespace] = []string{deployment.ObjectMeta.Name}

			// Otherwise, lets add the deployment to the list
		} else {
			response.Results[deployment.ObjectMeta.Namespace] = append(response.Results[deployment.ObjectMeta.Namespace], deployment.ObjectMeta.Name)
		}
	}

	// Write the response
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	return enc.Encode(response.Results)
}

func namespaceExists(d deploymentsLister) bool {
	ns, err := d.client.CoreV1.Namespaces()
}
