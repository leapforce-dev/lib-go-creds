package creds

import (
	"fmt"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

const (
	apiName string = "creds"
)

type Service struct {
	domain      string
	apiKey      string
	httpService *go_http.Service
}

type ServiceConfig struct {
	Domain string
	ApiKey string
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.Domain == "" {
		return nil, errortools.ErrorMessage("Domain not provided")
	}

	if serviceConfig.ApiKey == "" {
		return nil, errortools.ErrorMessage("ApiKey not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		domain:      serviceConfig.Domain,
		apiKey:      serviceConfig.ApiKey,
		httpService: httpService,
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	header := http.Header{}
	header.Set("X-Api-Key", service.apiKey)
	requestConfig.NonDefaultHeaders = &header

	return service.httpService.HttpRequest(requestConfig)
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", service.domain, path)
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return ""
}

func (service *Service) ApiCallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) ApiReset() {
	service.httpService.ResetRequestCount()
}
