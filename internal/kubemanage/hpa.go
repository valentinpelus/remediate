package kubemanage

import (
	"context"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetHpa(podInfo map[string]interface{}, clientset *kubernetes.Clientset) bool {

	// Get HPA by it's name and check if it's present in the namespace
	deployment, err := getDeployment(podInfo, clientset)
	log.Info().Msgf("hpa.go Getting Deployment %s from namespace %s", deployment.deploymentName, deployment.deploymentNamespace)
	if err != nil {
		log.Error().Msgf("Error in getting HPA %s from namespace %s", deployment.deploymentName, deployment.deploymentNamespace)
	}

	hpa, err := clientset.AutoscalingV1().HorizontalPodAutoscalers(deployment.deploymentNamespace).Get(context.TODO(), deployment.deploymentName, metav1.GetOptions{})
	log.Trace().Msgf("hpa.go Getting HPA global :  %s", hpa)
	hpaInfo := make(map[string]interface{})
	hpaInfo["hpaName"] = hpa.Name
	hpaInfo["namespace"] = hpa.Namespace
	log.Info().Msgf("hpa.go Getting HPA %s from namespace %s", hpaInfo["hpaName"], hpaInfo["namespace"])
	if err != nil {
		log.Error().Msgf("Error in getting HPA %s from namespace %s", hpaInfo["hpaName"], hpaInfo["namespace"])
	}
	return false
}

/* func describeHpa(hpaName string, namespace string, clientset *kubernetes.Clientset) (interface{}, error) {

	// Describe HPA details before executing any remediations
	hpa, err := clientset.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).Get(context.TODO(), hpaName, metav1.GetOptions{})
	if err != nil {
		log.Error().Msgf("Error in getting HPA %s from namespace %s", hpaName, namespace)
		panic(err)
	}
	hpaMap := make(map[string]interface{})
	hpaMap["Name"] = hpa.Name
	hpaMap["Namespace"] = hpa.Namespace
	hpaMap["MinReplicas"] = hpa.Spec.MinReplicas
	hpaMap["MaxReplicas"] = hpa.Spec.MaxReplicas
	return hpaMap, nil
} */
