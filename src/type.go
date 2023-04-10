package main

import "time"

type Assets struct {
	Banner string
}

// ConfigObject is the object that contains the config information
type ConfigObject struct {
	IgnoreConfigFlag bool             `yaml:"configFlag"` // If the configFlag is true, config.yaml | config.yml is to be ignored
	ConfigFlag       bool             `yaml:"configFlag"` // If the configFlag is true, config.yaml | config.yml exists
	SMTPConfig       SMTPConfigObject `yaml:"smtpConfig"` // SMTP config information
}

type SMTPConfigObject struct {
	SMTPServer  []SMTPServerObject `yaml:"smtpServer"`  // SMTP server information (host, port, username, password, fromEmail)
	TargetEmail []string           `yaml:"targetEmail"` // Target email address
}

type SMTPServerObject struct {
	Host      string `yaml:"host"`      // SMTP server host
	Port      int    `yaml:"port"`      // SMTP server port
	Username  string `yaml:"username"`  // SMTP server username
	Password  string `yaml:"password"`  // SMTP server password
	FromEmail string `yaml:"fromEmail"` // SMTP server from email address
}

// AlertManagerPayloadObject is the object that contains the payload information
type AlertManagerPayloadObject struct {
	Receiver    string        `json:"receiver"`
	Status      string        `json:"status"`
	Alert       []AlertObject `json:"alerts"`
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

type AlertObject struct {
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
}
