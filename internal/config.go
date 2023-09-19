package internal

import (
	"os"
	"strings"
)

type MapperConfig struct {
	newHttpEntrypointName  string
	newHttpsEntrypointName string
	addMiddleware          string
	newCertResolver        string
	removeHttpRouters      bool
	removeWwwRouters       bool
	removeCoolify          bool
	removeWwwMiddlewares   bool
	ignoreMiddlewareSites  []string
}

func ReadMapperConfigFromENV() MapperConfig {
	return MapperConfig{
		newCertResolver:        getStringFromEnv("NEW_CERT_RESOLVER"),
		newHttpEntrypointName:  getStringFromEnv("NEW_HTTP_ENTRYPOINT"),
		newHttpsEntrypointName: getStringFromEnv("NEW_HTTPS_ENTRYPOINT"),
		addMiddleware:          getStringFromEnv("ADD_MIDDLEWARE_NAME"),
		ignoreMiddlewareSites:  getStringArrayFromEnv("IGNORE_MIDDLEWARE_SITES"),
		removeCoolify:          getBooleanFromEnv("FILTER_COOLIFY"),
		removeWwwMiddlewares:   getBooleanFromEnv("FILTER_WWW_MIDDLEWARE"),
		removeHttpRouters:      getBooleanFromEnv("FILTER_HTTP_ROUTERS"),
		removeWwwRouters:       getBooleanFromEnv("FILTER_WWW_ROUTERS"),
	}
}

func getStringFromEnv(name string) string {
	fullName := "TRAEFIK_MAPPER_" + name
	return os.Getenv(fullName)
}

func getStringArrayFromEnv(name string) []string {
	stringValue := getStringFromEnv(name)
	if stringValue == "" {
		return nil
	}
	return strings.Split(stringValue, ";")
}

func getBooleanFromEnv(name string) bool {
	return getStringFromEnv(name) == "true"
}
