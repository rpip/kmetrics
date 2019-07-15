package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Pod is an extraction of relevant fields from the Kubernetes pods structure
type Pod struct {
	Name             string `json:"name"`
	ApplicationGroup string `json:"applicationGroup"`
	RunningPodsCount int32  `json:"runningPodsCount"`
}

// kubernetes will hold the connection info and methods for interacting with the cluster
type kubeClient struct {
	client *kubernetes.Clientset //typev1.CoreV1Interface
}

func NewKubeClient(kubeconfig string) *kubeClient {
	return &kubeClient{
		client: getKubeClient(kubeconfig),
	}
}

// GetServices returns a list of services running in the cluster in the given namespace
func (k8s *kubeClient) GetServices(namespace string, group string) ([]Pod, error) {
	listOptions := metav1.ListOptions{}
	if group != "" {
		listOptions.LabelSelector = fmt.Sprintf("applicationGroup=%s", group)
	}

	deploymentClient := k8s.client.ExtensionsV1beta1().Deployments(namespace)
	deployments, err := deploymentClient.List(listOptions)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not retrieve pods in %s namespace", namespace)
	}

	var pods []Pod
	for _, deployment := range deployments.Items {
		pod := Pod{
			Name:             deployment.Name,
			ApplicationGroup: deployment.Labels["applicationGroup"],
			RunningPodsCount: *deployment.Spec.Replicas,
		}
		pods = append(pods, pod)
	}
	return pods, nil
}

func getKubeClient(kubeconfig string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	return clientset
}
