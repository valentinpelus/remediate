package kubemanage

import (
	"bytes"
	"context"
	"io"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/rs/zerolog/log"
)

var podLabelTarget string

func DeletePod(podInfo map[string]interface{}, clientset *kubernetes.Clientset) bool {

	gracePeriod := int64(0)
	if checkPodPresence(podInfo, clientset) {
		log.Info().Msgf("pod.go Deleting pod %s", podInfo["podName"])
		if err := clientset.CoreV1().Pods(podInfo["namespace"].(string)).Delete(context.TODO(), podInfo["podName"].(string), metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod}); err != nil {
			log.Info().Msgf("pod.go Error in deletion of pod %s", podInfo["podName"])
			return false
		}
		return true
	}
	return false
}

func GetPod(podName string, namespace string, clientset *kubernetes.Clientset) (map[string]string, error) {

	// Get pod by it's name and check if it's present in the namespace, it will help to target the required project's pods with it's label
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	podMap := make(map[string]string)
	podMap["Name"] = pod.Name
	podMap["Kind"] = pod.Kind
	podMap["Namespace"] = pod.Namespace
	podMap["NodeName"] = pod.Spec.NodeName
	podMap["Status"] = string(pod.Status.Phase)
	if len(pod.Labels["app.kubernetes.io/instance"]) > 0 {
		podMap["LabelInstance"] = pod.Labels["app.kubernetes.io/instance"]
	}
	if len(pod.Labels["app.kubernetes.io/name"]) > 0 {
		podMap["LabelName"] = pod.Labels["app.kubernetes.io/name"]
	}
	if len(pod.Labels["project"]) > 0 {
		podMap["LabelProject"] = pod.Labels["project"]
	}
	if len(pod.OwnerReferences) > 0 {
		podMap["rsName"] = pod.OwnerReferences[0].Name
	}

	return podMap, nil
}

func checkPodPresence(podInfo map[string]interface{}, clientset *kubernetes.Clientset) bool {

	podName := podInfo["podName"].(string)
	namespace := podInfo["namespace"].(string)
	// Get pod by it's name and check if it's present in the namespace, it will help to target the required project's pods with it's label
	podQuery, err := GetPod(podName, namespace, clientset)
	if err != nil {
		log.Error().Msgf("pod.go Error in getting pod %s from namespace %s", podName, namespace)
		return false
	}

	// We will check the labels of our pod to find the right target to list all pods concerning the same project
	if _, ok := podQuery["LabelProject"]; ok {
		// If the label project exist we retrieve it and set it as target"
		podLabelTarget = "project=" + podQuery["LabelProject"]
	} else {
		// If there is no label project we set the label app.kubernetes.io/name as target
		podLabelTarget = "app.kubernetes.io/name=" + podQuery["app.kubernetes.io/name"]
	}

	log.Info().Msgf("pod.go Pod name %s and Label Target : %s", podName, podLabelTarget)

	// Listing Pods from chosen namespace, targeting the right label
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: podLabelTarget})
	log.Info().Msgf("pod.go Listing pods from namespace %s", pods)
	if err != nil {
		log.Warn().Msgf("pod.go Error issued during pod listing on namespace. Backing off")
		return false
	}
	log.Info().Msgf("pod.go Searching pod %s in namespace %s", podName, namespace)
	// Using func checkQuotaPod to check the ammount of healthy pod on our project

	if getPodCountForNamespace(namespace, podLabelTarget, clientset) >= 2 {
		log.Info().Msgf("pod.go More than 2 pod on namespace %s can proceed to actions", namespace)
		// Making sure we are targeting running pod and not backoff/restarting one
		for _, podsList := range (pods).Items {
			if (podsList.Name == podName) && (podsList.Status.Phase == "Running") {
				log.Info().Msgf("Found pod %s in namespace %s in status %s", podName, namespace, podsList.Status.Phase)
				return true
			}
		}
	} else {
		log.Warn().Msgf("pod.go Warning issued during deletion, the ratio between unhealthy pod and healthy one is not optimal, therefore no action will be taken. Backing off")
		return false
	}
	return false
}

func getPodCountForNamespace(namespace string, podLabelTarget string, clientset *kubernetes.Clientset) int {

	// Using this func to check if we have more than one pod on our namespace before taking any action, avoiding creating chain reaction
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: podLabelTarget})
	if err != nil {
		log.Error().Msgf("pod.go Error in getting number pods from namespace %s", namespace)
		return 0
	}
	log.Info().Msgf("pod.go Checking number of pod on namespace %s before taking actions", namespace)

	// If we have more than the ammount of pod returned by the alert - 2 and if we have at minimum 2 pod on the project running, then we can proceed
	// It aims to avoid deleting all pods of the same project directly
	return len(pods.Items)
}

func GetLogPod(podInfo map[string]interface{}, clientset *kubernetes.Clientset, follow bool) string {

	podName := podInfo["podName"].(string)
	namespace := podInfo["namespace"].(string)
	// Pod Log options
	count := int64(100)
	podLogOptions := v1.PodLogOptions{
		Container: "test5",
		Follow:    false,
		TailLines: &count,
	}
	log.Info().Msgf("pod.go LogOptions : %s", &podLogOptions)
	log.Info().Msgf("pod.go Getting Pod Log %s from namespace %s", podName, namespace)

	// Get pod by it's name and check if it's present in the namespace, it will help to target the required project's pods with it's label
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &podLogOptions)
	if req != nil {
		log.Error().Msgf("pod.go Error in getting pod %s from namespace %s", podName, namespace)
	}
	stream, err := req.Stream(context.TODO())
	if err != nil {
		log.Error().Msgf("pod.go Error in opening stream for pod %s from namespace %s", podName, namespace)
	}
	defer stream.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, stream)
	if err != nil {
		log.Error().Msgf("pod.go Error in copying pod logs for pod %s from namespace %s", podName, namespace)
	}
	str := buf.String()

	log.Info().Msgf("pod.go Getting Pod Log %s", str)
	return str

	//return nil, nil
}
