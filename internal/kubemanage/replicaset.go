package kubemanage

import (
	"context"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getReplicaSet(podMap map[string]string, clientset *kubernetes.Clientset) (map[string]string, error) {

	podNamespace := podMap["Namespace"]
	podRS := podMap["RSName"]
	log.Info().Msgf("deployment.go Getting RS %s from namespace %s", podRS, podNamespace)
	// Get HPA by it's name and check if it's present in the namespace
	rs, err := clientset.AppsV1().ReplicaSets(podNamespace).Get(context.TODO(), podRS, metav1.GetOptions{})
	if err != nil {
		log.Error().Msgf("Error in getting RS %s from namespace %s", podRS, podNamespace)
	}
	log.Info().Msgf("deployment.go Getting RS in getReplicaSet %s from namespace %s deployment %s", rs.Name, rs.Namespace, rs.OwnerReferences[0].Name)
	rsMap := make(map[string]string)
	rsMap["Name"] = rs.Name
	rsMap["Namespace"] = rs.Namespace
	rsMap["Deployment"] = rs.OwnerReferences[0].Name

	return rsMap, nil
}
