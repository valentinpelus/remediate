package kubemanage

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	policyv1beta "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type err struct {
	err error
}

func GetNode(nodeName string, clientset *kubernetes.Clientset) (node *v1.Node, err error) {

	log.Info().Msgf("node.go Getting Node %s from cluster", nodeName)
	nodeQuery, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Error().Msgf("node.go Error in getting node %s: %v", nodeName, err)
		return nil, err
	}

	return nodeQuery, nil
}

func CordonNode(nodeName string, clientset *kubernetes.Clientset) (bool, error) {

	log.Info().Msgf("node.go Cordoning Node %s", nodeName)
	node, err := GetNode(nodeName, clientset)
	if err != nil {
		log.Error().Msgf("node.go Error in getting node %s: %v", nodeName, err)
		return false, err
	}

	node.Spec.Unschedulable = true
	cordonAction, err := clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		log.Error().Msgf("node.go Error in cordoning node %s: %v", nodeName, err)
		return false, err
	}

	// Watch node state until it's fully cordoned
	for {
		log.Info().Msgf("node.go Waiting for node %s to be cordoned", nodeName)
		node, err := GetNode(nodeName, clientset)
		if err != nil {
			log.Error().Msgf("node.go Error in getting node %s: %v", nodeName, err)
			return false, err
		}
		if node.Spec.Unschedulable == true {
			break
		} else {
			log.Info().Msgf("node.go Node schedule status : %t", cordonAction.Spec.Unschedulable)
			time.Sleep(5 * time.Second)
			continue
		}
	}

	log.Info().Msgf("node.go Node %s cordonned", cordonAction.Name)

	return true, nil
}

func UncordonNode(nodeName string, clientset *kubernetes.Clientset) (bool, error) {

	log.Info().Msgf("node.go Uncordoning Node %s", nodeName)
	node, err := GetNode(nodeName, clientset)
	if err != nil {
		log.Error().Msgf("node.go Error in getting node %s: %v", nodeName, err)
		return false, err
	}

	node.Spec.Unschedulable = false
	uncordonAction, err := clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		log.Error().Msgf("node.go Error in uncordoning node %s: %v", nodeName, err)
		return false, err
	}

	// Watch node state until it's fully uncordoned
	for {
		log.Info().Msgf("node.go Waiting for node %s to be uncordoned", nodeName)
		node, err := GetNode(nodeName, clientset)
		if err != nil {
			log.Error().Msgf("node.go Error in getting node %s: %v", nodeName, err)
			return false, err
		}
		if node.Spec.Unschedulable == false {
			break
		} else {
			time.Sleep(5 * time.Second)
			continue
		}
	}

	log.Trace().Msgf("node.go Node uncordon status %s", &uncordonAction.Status)

	return true, nil
}

func DrainNode(nodeName string, clientset *kubernetes.Clientset) (bool, error) {

	log.Info().Msgf("node.go Draining Node %s", nodeName)
	node, err := GetNode(nodeName, clientset)
	if err != nil {
		log.Error().Msgf("node.go Error in getting node %s: %v", nodeName, err)
		return false, err
	}

	log.Info().Msgf("node.go Node status %s", node.Status.Conditions)

	// List all pods on the node
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		log.Error().Msgf("node.go Error listing pods on node %s: %v", nodeName, err)
		return false, err
	}

	for _, pod := range pods.Items {
		// Skip daemonset-managed pods
		for _, ownerRef := range pod.ObjectMeta.OwnerReferences {
			if ownerRef.Kind == "DaemonSet" {
				log.Info().Msgf("node.go Skipping eviction for DaemonSet pod %s", pod.Name)
				continue
			}
		}
		// Skip pods with emptyDir volumes if --delete-empty-dir-data is not set
		hasEmptyDir := false
		for _, volume := range pod.Spec.Volumes {
			if volume.EmptyDir != nil {
				hasEmptyDir = true
				break
			}
		}
		if hasEmptyDir {
			log.Info().Msgf("node.go Evicting pod %s with emptyDir volumes", pod.Name)
		}
		// Evict the pod
		log.Info().Msgf("node.go Evicting pod %s from node %s", pod.Name, nodeName)
		eviction := &policyv1beta.Eviction{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pod.Name,
				Namespace: pod.Namespace,
			},
			DeleteOptions: &metav1.DeleteOptions{
				GracePeriodSeconds: new(int64), // Set grace period to 0 for immediate deletion
			},
		}

		err := clientset.CoreV1().Pods(pod.Namespace).Evict(context.TODO(), eviction)
		if err != nil {
			log.Error().Msgf("node.go Error evicting pod %s: %v", pod.Name, err)
			return false, err
		}
	}

	log.Info().Msgf("node.go Node %s drained successfully", nodeName)

	return true, nil
}

func DeleteNode(nodeName string, clientset *kubernetes.Clientset) (bool, error) {
	log.Info().Msgf("node.go Deleting Node %s", nodeName)
	node, err := GetNode(nodeName, clientset)
	if err != nil {
		log.Error().Msgf("node.go Error in getting node %s: %v", nodeName, err)
		return false, err
	}

	log.Info().Msgf("node.go Preparing deletion of node %s", node.Name)

	err = clientset.CoreV1().Nodes().Delete(context.TODO(), nodeName, metav1.DeleteOptions{})
	if err != nil {
		log.Error().Msgf("node.go Error in deleting node %s: %v", nodeName, err)
		return false, err
	}

	log.Info().Msgf("node.go Node %s deleted successfully", nodeName)

	return true, nil
}
