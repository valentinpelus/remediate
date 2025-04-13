package kuberemediate

import (
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/valentinpelus/remediate/internal/kubemanage"
	"k8s.io/client-go/kubernetes"
)

type actionState struct {
	mutex      sync.Mutex
	isRunning  bool
	alertLabel string
	podName    string
	nodeName   string
}

func podNotSchedulable(podName string, namespace string, clientset *kubernetes.Clientset) (bool, error) {

	actionState := actionState{
		isRunning:  false,
		alertLabel: "podnotschedulable",
		podName:    podName,
		nodeName:   "",
	}
	actionState.mutex.Lock()
	log.Info().Msgf("podnotschedulable.go Get node from pod %s", podName)
	pod, err := kubemanage.GetPod(podName, namespace, clientset)
	if err != nil {
		log.Error().Msgf("podnotschedulable.go Error in getting pod %s : %s", podName, err)
		return false, err
	}

	log.Info().Msgf("podnotschedulable.go Query logs for pod %s", pod["NodeName"])

	nodeName := pod["NodeName"]
	// Cordoning node
	node, err := kubemanage.CordonNode(nodeName, clientset)
	if err != nil {
		log.Error().Msgf("podnotschedulable.go Error in cordonning node %s : %s", nodeName, err)
		return false, err
	}

	// Check if node is cordoned
	if node {
		log.Info().Msgf("podnotschedulable.go Node cordoned %s", nodeName)

		// Draining node
		nodeDrain, err := kubemanage.DrainNode(nodeName, clientset)
		if err != nil {
			log.Error().Msgf("podnotschedulable.go Error in draining node %s : %s", nodeName, err)
			return false, err
		}

		// Check if node is drained
		if nodeDrain {
			log.Info().Msgf("podnotschedulable.go Node drained %s", nodeName)

			// Deleting node
			nodeDelete, err := kubemanage.DeleteNode(nodeName, clientset)
			if err != nil {
				log.Error().Msgf("podnotschedulable.go Error in deleting node %s : %s", nodeName, err)
				return false, err
			}
			// Check if node is deleted
			if nodeDelete {
				log.Info().Msgf("podnotschedulable.go Node deleted %s", nodeName)
			}
		}
	}

	//time.Sleep(5 * time.Second)
	log.Info().Msgf("podnotschedulable.go Node cordoned, drained and deleted %s", nodeName)
	return true, nil
}
