package kuberemediate

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Response struct {
	Status string `json:"status"`
	Data   []Data `json:"data"`
}
type Labels struct {
	AdminAlert  string `json:"admin_alert"`
	Alertgroup  string `json:"alertgroup"`
	AlertName   string `json:"alertname"`
	ClusterName string `json:"cluster_name"`
	Namespace   string `json:"namespace"`
	Pod         string `json:"pod"`
}

type Data struct {
	Labels Labels `json:"labels,omitempty"`
}

var (
	Client HTTPClient
)

var alertExtractList [][]string

func GetAlertList(server string, SupportedAlert [][]string) (array [][]string) {
	// Initialisation of GET request
	res, err := http.Get(server)
	if err != nil {
		log.Fatal().Msgf("alert.go Error in GET request %s ", err)
	}
	// Closing request
	defer res.Body.Close()

	// Reading Body content
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal().Msgf("alert.go Error in reading body %s ", err)
	}

	// Serialising return of Body into JSON
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal().Msgf("alert.go Error in reading body %s ", err)
	}

	wg := new(sync.WaitGroup)
	ch := make(chan [][]string)
	alertExtractList = nil
	// Parsing Json return to match Alertname with the slice sent to the function
	for _, alertsData := range response.Data {
		// Adding a waitgroup to wait for all goroutines to finish
		wg.Add(1)
		go parseAlertList(alertsData, SupportedAlert, ch, wg)
	}

	// Waiting for all goroutines to finish
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Reading from the channel and appending the result to the alertExtractList
	for alertMatch := range ch {
		alertExtractList = append(alertExtractList, alertMatch...)
	}

	return alertExtractList
}

func parseAlertList(alertsData Data, SupportedAlert [][]string, ch chan [][]string, wg *sync.WaitGroup) (array [][]string) {
	defer wg.Done()
	var alertMatch [][]string
	alertName := alertsData.Labels.AlertName
	log.Info().Msgf("alert.go alertName %s podName %s", alertName, alertsData.Labels.Pod)

	// Iterating over the SupportedAlert slice to find a match
	for i := range SupportedAlert {
		if alertName == SupportedAlert[i][0] && len(alertsData.Labels.Pod) > 0 {
			// If we find a match we append the podname, namespace, alertaction and alertname to our slice
			podName := alertsData.Labels.Pod
			namespace := alertsData.Labels.Namespace
			alertAction := SupportedAlert[i][1]
			alertMatch := append(alertMatch, []string{podName, namespace, alertAction, alertName})
			log.Info().Msgf("alert.go Alert matching a rule in remediate, alertName %s for podName %s, executing following action %s", alertName, podName, alertAction)
			// Sending the alertMatch to the channel and return array
			ch <- alertMatch
			return
		} else {
			log.Info().Msgf("alert.go No alerts matching any enabled rules in remediate, continuing...")
		}
	}
	return
}

func ParseMatchList(alertPodExtractList [][]string, confPath *string, clientset *kubernetes.Clientset) {
	for i := range alertPodExtractList {
		log.Info().Msgf("main.go Iteration : %v", i)
		// Creating a map to store the pod information
		podInfo := make(map[string]interface{})
		podInfo["podName"] = alertPodExtractList[i][0]
		podInfo["namespace"] = alertPodExtractList[i][1]
		podInfo["alertAction"] = alertPodExtractList[i][2]
		podInfo["alertName"] = alertPodExtractList[i][3]
		podInfo["podCount"] = len(alertPodExtractList)

		podName := podInfo["podName"].(string)
		namespace := podInfo["namespace"].(string)
		// Check if podName and namespace are not empty
		if (len(podName) > 0) && (len(namespace) > 0) {
			log.Info().Msgf("alert.go Detecting pod %s in namespace %s", podInfo["podName"], podInfo["namespace"])
			// Parse returned alertPodExtractList to determine which action should be done with remediate
			switch podInfo["alertAction"] {
			case "deletePod":
				log.Info().Msgf("alert.go Delete pod %s in namespace %s in error", podInfo["podName"], podInfo["namespace"])
				triggeredAction := DeletePod(podInfo, clientset)
				time.Sleep(5 * time.Second)
				if triggeredAction {
					postMessageSlack(podInfo["alertName"].(string), namespace, confPath)
				}
			case "enrichAlert":
				//kuberemediate.DescribeDeployment(podName, clientset, namespace)

			}
		}
	}
}
