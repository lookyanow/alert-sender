package main

import (
	"net/http"

	_ "net/http/pprof"

	"alert-sender/config"
	"alert-sender/handlers"
	log "alert-sender/logger"
	"alert-sender/phones"
	"alert-sender/sms"

	"go.uber.org/zap"
)

func main() {
	log := log.NewLogger()
	log.Info("alert-sender service starting")
	appConf := config.NewConfigFromEnv()

	rg, err := phones.NewReceiverGroup(appConf.PhoneFile)
	if err != nil {
		log.Fatal("Can't get phone groups: %s", zap.Error(err))
	}
	smsSender, err := sms.NewSender(appConf.SmsURL, appConf.SmsUser, appConf.SmsPassword)
	if err != nil {
		log.Fatal("Can't get SMS Api config: %s", zap.Error(err))
	}

	var phoneGroups = make(map[string][]string)
	for _, v := range *rg {
		phoneGroups[v.Name] = v.Phones
	}

	h := handlers.NewAppHandler(appConf, smsSender, phoneGroups)
	// Handlers
	http.HandleFunc("/live", h.LiveHandler)
	http.HandleFunc("/prometheus", h.PrometheusWebHookHandler)
	http.HandleFunc("/grafana", h.GrafanaWebHookHandler)
	http.HandleFunc("/ready", h.ReadyHandler)

	listenAddress := ":" + appConf.ServerPort

	log.Sugar().Infof("Listening on: %v", listenAddress)
	log.Fatal("Listener stoped with fatal error: %s", zap.Error(http.ListenAndServe(listenAddress, nil)))
}
