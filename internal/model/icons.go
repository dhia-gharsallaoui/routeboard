package model

import "strings"

var iconPatterns = map[string]string{
	"grafana":      "\U0001F4CA", // 📊
	"prometheus":   "\U0001F525", // 🔥
	"alertmanager": "\U0001F6A8", // 🚨
	"argocd":       "\U0001F680", // 🚀
	"argo-cd":      "\U0001F680", // 🚀
	"jenkins":      "\U0001F527", // 🔧
	"gitlab":       "\U0001F98A", // 🦊
	"gitea":        "\U0001F375", // 🍵
	"minio":        "\U0001FAA3", // 🪣
	"vault":        "\U0001F510", // 🔐
	"keycloak":     "\U0001F511", // 🔑
	"traefik":      "\U0001F500", // 🔀
	"nginx":        "\U0001F310", // 🌐
	"kibana":       "\U0001F4C8", // 📈
	"jaeger":       "\U0001F50D", // 🔍
	"longhorn":     "\U0001F402", // 🐂
	"rancher":      "\U0001F920", // 🤠
	"harbor":       "\U00002693", // ⚓
	"sonarqube":    "\U0001F52C", // 🔬
	"pgadmin":      "\U0001F418", // 🐘
	"redis":        "\U0001F534", // 🔴
	"rabbitmq":     "\U0001F430", // 🐰
}

const defaultIcon = "\U0001F310" // 🌐

func DetectIcon(serviceName, resourceName string) string {
	name := strings.ToLower(serviceName)
	if name == "" {
		name = strings.ToLower(resourceName)
	}
	for pattern, icon := range iconPatterns {
		if strings.Contains(name, pattern) {
			return icon
		}
	}
	return defaultIcon
}
