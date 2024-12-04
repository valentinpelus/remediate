package kuberemediate

import (
	conf "github.com/valentinpelus/remediate/internal/config"
	notif "github.com/valentinpelus/remediate/internal/notification"

	"github.com/rs/zerolog/log"
)

func postMessageSlack(alertName string, namespace string, confPath *string) {
	conf.LoadConfSlack(*confPath)

	log.Info().Msgf("Alert notif %s %s", alertName, namespace)

	url := conf.ConfigurationSlack.WebhookUrl
	username := conf.ConfigurationSlack.SlackClient.UserName
	channel := conf.ConfigurationSlack.SlackClient.Channel
	clusterName := conf.ConfigurationSlack.ClusterName

	log.Info().Msgf("Slack url %s", url)
	log.Info().Msgf("Slack clusterName %s", clusterName)

	// Loading slack
	sc := notif.SlackClient{
		WebHookUrl: url,
		UserName:   username,
		Channel:    channel,
	}

	sr := notif.SlackJobNotification{
		Title: "Remediate - The auto-remediation has been triggered",
		Text:  "Remediate - The auto-remediation has been triggered",
		Details: "*Alert remediated*: " + alertName +
			"\r\n *Cluster*: " + clusterName +
			"\r\n *Namespace*: " + namespace,
		Color:     "#5581d9",
		IconEmoji: "necron",
	}

	err := sc.SendJobNotification(sr)
	if err != nil {
		log.Fatal().Err(err)
	}
}
