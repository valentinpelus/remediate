package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	conf "github.com/valentinpelus/remediate/internal/config"
	kuberemediate "github.com/valentinpelus/remediate/internal/kuberemediate"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func main() {

	// Init of zerolog library and config
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	confPath := flag.String("conf", "config.yaml", "Config path")
	debug := flag.Bool("debug", false, "sets log level to debug")
	//tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Loading configuration files
	conf.LoadConfKube(*confPath)
	conf.LoadConfAlert(*confPath)

	// Loading kubeconfig file with context
	kube_config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// Creates the clientset
	clientset, err := kubernetes.NewForConfig(kube_config)
	if err != nil {
		panic(err.Error())
	}

	// Init http client
	Client = &http.Client{}

	// Init AMUrl to allow alerts query
	jsonUrl := conf.Conf.QueryURL + "/api/v1/alerts"

	ListSupportedAlert := conf.EnabledAlertList

	fmt.Println(ListSupportedAlert)

	for {
		time.Sleep(20 * time.Second)
		log.Info().Msgf("main.go Check ongoing")
		// Querying Alertmanager to check if alert is firing for backend size divergence and proceed to deletion if needed
		alertPodExtractList := kuberemediate.GetAlertList(jsonUrl, ListSupportedAlert)
		fmt.Println("Print of return alertPodExtractList : ", alertPodExtractList)
		// Looping in the alerts list returned by GetVMAlerMatch
		kuberemediate.ParseAlertList(alertPodExtractList, confPath, clientset)
	}
}
