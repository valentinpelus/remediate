package kuberemediate

import (
	"github.com/rs/zerolog/log"
	"github.com/valentinpelus/remediate/internal/kubemanage"
	"k8s.io/client-go/kubernetes"
)

func podNotSchedulable(podName string, namespace string, clientset *kubernetes.Clientset) (bool, error) {
	log.Info().Msgf("podnotschedulable.go Get node from pod %s", podName)
	pod, err := kubemanage.GetPod(podName, namespace, clientset)
	if err != nil {
		log.Error().Msgf("podnotschedulable.go Error in getting pod %s : %s", podName, err)
		return false, err
	}

	log.Info().Msgf("podnotschedulable.go Query logs for pod %s", pod["NodeName"])

	nodeName := pod["NodeName"]
	node, err := kubemanage.CordonNode(nodeName, clientset)
	if err != nil {
		log.Error().Msgf("podnotschedulable.go Error in cordonning node %s : %s", nodeName, err)
		return false, err
	}
	if node {
		log.Info().Msgf("podnotschedulable.go Node cordoned %s", nodeName)
	}
	nodeDrain, err := kubemanage.DrainNode(nodeName, clientset)
	if err != nil {
		log.Error().Msgf("podnotschedulable.go Error in draining node %s : %s", nodeName, err)
		return false, err
	}
	if nodeDrain {
		log.Info().Msgf("podnotschedulable.go Node drained %s", nodeName)
	}
	nodeDelete, err := kubemanage.DeleteNode(nodeName, clientset)
	if err != nil {
		log.Error().Msgf("podnotschedulable.go Error in deleting node %s : %s", nodeName, err)
		return false, err
	}
	if nodeDelete {
		log.Info().Msgf("podnotschedulable.go Node deleted %s", nodeName)
	}
	//time.Sleep(5 * time.Second)
	log.Info().Msgf("podnotschedulable.go Node cordoned, drained and deleted %s", nodeName)
	return true, nil
}
