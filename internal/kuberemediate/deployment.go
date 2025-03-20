package kuberemediate

import (
	"context"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type deployment struct {
	deploymentName      string
	deploymentNamespace string
}

func getDeployment(podInfo map[string]interface{}, clientset *kubernetes.Clientset) (deployment, error) {

	// Get Pod by it's name and check if it's present in the namespace to retrieve the RS name, create a map if so
	deplPod, err := GetPod(podInfo["podName"].(string), podInfo["namespace"].(string), clientset)
	if err != nil {
		log.Error().Msgf("Error in getting pod %s from namespace %s", podInfo["podName"], podInfo["namespace"])
	}
	podMap := make(map[string]string)
	podMap["Name"] = deplPod["Name"]
	podMap["Namespace"] = deplPod["Namespace"]
	podMap["RSName"] = deplPod["rsName"]
	log.Info().Msgf("deployment.go Printing Pod global : %s ", podMap)

	// Get ReplicaSet by it's name and check if it's present in the namespace
	replica, err := getReplicaSet(podMap, clientset)
	if err != nil {
		log.Error().Msgf("Error in executing func to get RS")
	}
	log.Info().Msgf("deployment.go Printing RS global : %s ", replica)

	// Get Deployment from RS's name
	depl, err := clientset.AppsV1().Deployments(replica["Namespace"]).Get(context.Background(), replica["Deployment"], metav1.GetOptions{})
	log.Info().Msgf("deployment.go Printing Deployment global : %s ", depl)
	deploymentName := depl.Name
	deploymentNamespace := depl.Namespace
	log.Info().Msgf("deployment.go Getting Deployment %s from namespace %s", deploymentName, deploymentNamespace)
	if err != nil {
		log.Error().Msgf("Error in getting HPA %s from namespace %s", deploymentName, deploymentNamespace)
	}
	return deployment{deploymentName, deploymentNamespace}, nil
}
