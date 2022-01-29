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
	code        string
	httpService *go_http.Service
}

type ServiceConfig struct {
	Domain string
	Code   string
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.Domain == "" {
		return nil, errortools.ErrorMessage("Domain not provided")
	}

	if serviceConfig.Code == "" {
		return nil, errortools.ErrorMessage("Code not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		domain:      serviceConfig.Domain,
		code:        serviceConfig.Code,
		httpService: httpService,
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	header := http.Header{}
	header.Set("X-Check-Code", service.code)
	requestConfig.NonDefaultHeaders = &header

	return service.httpService.HTTPRequest(requestConfig)
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", service.domain, path)
}

func (service *Service) APIName() string {
	return apiName
}

func (service *Service) APIKey() string {
	return ""
}

func (service *Service) APICallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) APIReset() {
	service.httpService.ResetRequestCount()
}
