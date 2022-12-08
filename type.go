package main

import "time"

type AlertManagerPayloadObject struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []struct {
		Status string `json:"status"`
		Labels struct {
			Alertname string `json:"alertname"`
			Env       string `json:"env"`
			Group     string `json:"group"`
			Instance  string `json:"instance"`
			Job       string `json:"job"`
			Loc       string `json:"loc"`
			Resp      string `json:"resp"`
			Severity  string `json:"severity"`
			Theme     string `json:"theme"`
			Type      string `json:"type"`
		} `json:"labels"`
		Annotations struct {
			Summary string `json:"summary"`
		} `json:"annotations"`
		StartsAt     time.Time `json:"startsAt"`
		EndsAt       time.Time `json:"endsAt"`
		GeneratorURL string    `json:"generatorURL"`
		Fingerprint  string    `json:"fingerprint"`
	} `json:"alerts"`
	GroupLabels struct {
		Alertname string `json:"alertname"`
	} `json:"groupLabels"`
	CommonLabels struct {
		Alertname string `json:"alertname"`
		Env       string `json:"env"`
		Group     string `json:"group"`
		Instance  string `json:"instance"`
		Job       string `json:"job"`
		Loc       string `json:"loc"`
		Resp      string `json:"resp"`
		Severity  string `json:"severity"`
		Theme     string `json:"theme"`
		Type      string `json:"type"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
		Summary string `json:"summary"`
	} `json:"commonAnnotations"`
	ExternalURL     string `json:"externalURL"`
	Version         string `json:"version"`
	GroupKey        string `json:"groupKey"`
	TruncatedAlerts int    `json:"truncatedAlerts"`
}
