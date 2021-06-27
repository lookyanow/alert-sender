package handlers

import (
	"alert-sender/config"
	"alert-sender/sms"
	"encoding/json"
	"net/http"
	"strings"

	"alert-sender/logger"

	"github.com/prometheus/alertmanager/template"
	"go.uber.org/zap"
)

type AppHandler struct {
	cfg            *config.Config
	smsSender      *sms.Sender
	recieverGroups map[string][]string
}

type GrafanaAlert struct {
	DashboardId int               `json:"dashboardId"`
	EvalMatches []EvalMatch       `json:"evalMatches"`
	ImageUrl    string            `json:"imageUrl"`
	Message     string            `json:"message"`
	OrgId       int               `json:"orgId"`
	PanelId     int               `json:"panelId"`
	RuleId      int               `json:"ruleId"`
	RuleName    string            `json:"ruleName"`
	RuleUrl     string            `json:"ruleUrl"`
	State       string            `json:"state"`
	Tags        map[string]string `json:"tags"`
	Title       string            `json:"title"`
}

type EvalMatch struct {
	Value  int    `json:"value"`
	Metric string `json:"metric"`
}

func NewAppHandler(cfg *config.Config, sms *sms.Sender, recieverGroups map[string][]string) *AppHandler {
	return &AppHandler{cfg: cfg, smsSender: sms, recieverGroups: recieverGroups}
}

func (a *AppHandler) PrometheusWebHookHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.NewLogger()
	defer r.Body.Close()
	data := template.Data{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Error("PromWebHookHandler unmarshaling request payload error: %s", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	for _, alert := range data.Alerts {
		log.Sugar().Infof("Alert: status=%s,Labels=%v,Annotations=%v", alert.Status, alert.Labels, alert.Annotations)

		var group string
		value, ok := alert.Labels["group"]
		if ok {
			group = value
		} else {
			group = "default"
		}
		severity := alert.Labels["severity"]
		numbers := a.recieverGroups[group]

		switch strings.ToUpper(severity) {
		case "CRITICAL":
			err := a.smsSender.SendSMS(alert.Annotations["description"], alert.Status, numbers)
			if err != nil {
				log.Error("Error in sending sms: %s", zap.Error(err))
			}
		default:
			log.Sugar().Infof("No action on severity: %s", severity)
		}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("success"))
}

func (a *AppHandler) GrafanaWebHookHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.NewLogger()
	defer r.Body.Close()

	// Unmarshaling alert data structure
	data := GrafanaAlert{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Error("GrafanaWebHookHandler unmarshaling request payload error: %s", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	// sending sms if group name found in alert tags
	group := getGroupName(data.Tags)
	if group != "" && len(a.recieverGroups[group]) != 0 {
		groupNumbers := a.recieverGroups[group]
		err := a.smsSender.SendSMS(data.Message, data.State, groupNumbers)
		if err != nil {
			log.Error("Error in sending sms: %s", zap.Error(err))
		}
		log.Sugar().Infof("Alert: %s, status: %s, group: %s, phones: %s", data.RuleName, data.State, group, groupNumbers)
	}

	// sending sms if phones tag found
	numbers := getPhoneNumbers(data.Tags)
	if numbers != nil {
		err := a.smsSender.SendSMS(data.Message, data.State, numbers)
		if err != nil {
			log.Error("Error in sending sms: %s", zap.Error(err))
		}
		log.Sugar().Infof("Alert: %s, status: %s, phones: %s", data.RuleName, data.State, numbers)
	} else {
		log.Sugar().Infof("Alert: %s, status: %s, has empty recipient phone list", data.RuleName, data.State)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("success"))
}

func (a *AppHandler) LiveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Alive!"))
}

func (a *AppHandler) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Ready!"))
}
