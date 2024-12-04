package kuberemediate

/* func getHpa(podInfo map[string]interface{}, clientset *kubernetes.Clientset) bool {

	// Get HPA by it's name and check if it's present in the namespace
	hpa, err := clientset.AutoscalingV2beta2().HorizontalPodAutoscalers(podInfo["namespace"].(string)).Get(context.TODO(), podInfo["hpaName"].(string), metav1.GetOptions{})
	if err != nil {
		log.Error().Msgf("Error in getting HPA %s from namespace %s", hpa["hpaName"], hpa["namespace"])
		panic(err)
	}
	return false
} */

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
