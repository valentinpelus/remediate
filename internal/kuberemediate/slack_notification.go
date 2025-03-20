package kuberemediate

import (
	conf "github.com/valentinpelus/remediate/internal/config"
	notif "github.com/valentinpelus/remediate/internal/notification"

	"github.com/rs/zerolog/log"
)

func postMessageSlack(alertInfo map[string]interface{}, confPath *string) {
	conf.LoadConfSlack(*confPath)

	log.Info().Msgf("Alert notif %s %s", alertInfo["alertName"], alertInfo["namespace"])

	url := conf.ConfigurationSlack.WebhookUrl
	username := conf.ConfigurationSlack.SlackClient.UserName
	channel := conf.ConfigurationSlack.SlackClient.Channel
	clusterName := conf.ConfigurationSlack.ClusterName
	alertName := alertInfo["alertName"].(string)
	podName := alertInfo["podName"].(string)
	namespace := alertInfo["namespace"].(string)

	//log.Info().Msgf("Slack url %s", url)
	log.Info().Msgf("Slack clusterName %s", clusterName)

	slackDetail := "*Alert triggered*: " + alertName +
		"\r\n *Cluster*: " + clusterName +
		"\r\n *Namespace*: " + namespace +
		"\r\n *Pod*: " + podName +
		"\r\n *Namespace*: " + namespace

	// Loading slack
	sc := notif.SlackClient{
		WebHookUrl: url,
		UserName:   username,
		Channel:    channel,
	}

	sr := notif.SlackJobNotification{
		Title:     "Remediate has been triggered",
		Text:      "Remediate has been triggered",
		Details:   slackDetail,
		Color:     "#5581d9",
		IconEmoji: "necron",
	}

	err := sc.SendJobNotification(sr)
	if err != nil {
		log.Fatal().Err(err)
	}
}
