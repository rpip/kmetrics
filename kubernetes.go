package main

import (
	"log"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// Pod is an extraction of relevant fields from the Kubernetes pods structure
type Pod struct {
	Name             string `json: name`
	ApplicationGroup string `json: applicationGroup`
	RunningPodsCount int    `json: runningPodsCount`
}

// kubernetes will hold the connection info and methods for interacting with the cluster
type kubeClient struct {
	api typev1.CoreV1Interface
}

func NewKubeClient(kubeconfig string) *kubeClient {
	return &kubeClient{
		api: getKubeClient(kubeconfig),
	}
}

// GetPods returns a list of pods running in the cluster in the given namespace
func (k8s *kubeClient) GetPods(namespace string) ([]Pod, error) {
	listOptions := metav1.ListOptions{}
	k8sPods, err := k8s.api.Pods(namespace).List(listOptions)

	if err != nil {
		return nil, errors.Wrapf(err, "Could not retrieve pods in %s namespace", namespace)
	}

	// first group the pods by the service name
	groupedPods := make(map[string][]corev1.Pod)
	for _, p := range k8sPods.Items {
		groupedPods[p.Labels["service"]] = append(groupedPods[p.Labels["service"]], p)
	}

	var pods []Pod
	for serviceName, runningPods := range groupedPods {
		pod1 := runningPods[0]

		pod := Pod{
			Name:             serviceName,
			ApplicationGroup: pod1.Labels["applicationGroup"],
			RunningPodsCount: len(runningPods),
		}
		pods = append(pods, pod)
	}
	return pods, nil
}

func getKubeClient(kubeconfig string) typev1.CoreV1Interface {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	return clientset.CoreV1()
}
