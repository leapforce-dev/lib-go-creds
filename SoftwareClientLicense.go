package creds

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"cloud.google.com/go/civil"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type SoftwareClientLicense struct {
	SoftwareClientLicenseGuid  string
	ObjectID                   string
	ObjectName                 string
	SoftwareClientLicenseState string
	DataState                  string
	AsOf                       civil.Date
	Values                     map[string]string
}

type GetSoftwareClientLicensesConfig struct {
	SoftwarePackageGuid       string
	RequiredKeys              *[]string
	CompanyID                 *int64
	SoftwareClientLicenseGuid *string
}

func (service *Service) GetSoftwareClientLicenses(config *GetSoftwareClientLicensesConfig) (*[]SoftwareClientLicense, *errortools.Error) {
	_softwareClientLicenses := []struct {
		SoftwareClientLicenseGuid  string
		ObjectID                   string
		ObjectName                 string
		SoftwareClientLicenseState string
		DataState                  string
		AsOf                       *civil.Date
		Key1                       string
		Value1                     string
		Key2                       string
		Value2                     string
		Key3                       string
		Value3                     string
		Key4                       string
		Value4                     string
		Key5                       string
		Value5                     string
		Key6                       string
		Value6                     string
		Key7                       string
		Value7                     string
	}{}

	if config == nil {
		return nil, errortools.ErrorMessage("Config is required")
	}

	var values url.Values = url.Values{}
	values.Set("pkg", config.SoftwarePackageGuid)

	if config.CompanyID != nil {
		values.Set("obj", fmt.Sprintf("%v", *config.CompanyID))
	}
	if config.SoftwareClientLicenseGuid != nil {
		values.Set("swr", *config.SoftwareClientLicenseGuid)
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		URL:           service.url(fmt.Sprintf("creds?%s", values.Encode())),
		ResponseModel: &_softwareClientLicenses,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	softwareClientLicenses := []SoftwareClientLicense{}

	for _, _softwareClientLicense := range _softwareClientLicenses {
		if _softwareClientLicense.SoftwareClientLicenseState != "Active" {
			fmt.Printf("%s skipped (not active)\n", _softwareClientLicense.ObjectName)
			continue
		}

		if _softwareClientLicense.AsOf == nil {
			fmt.Printf("%s skipped (no as of date)\n", _softwareClientLicense.ObjectName)
			continue
		}

		values := make(map[string]string)
		if strings.Trim(_softwareClientLicense.Key1, " ") != "" {
			values[strings.Trim(_softwareClientLicense.Key1, " ")] = strings.Trim(_softwareClientLicense.Value1, " ")
		}
		if strings.Trim(_softwareClientLicense.Key2, " ") != "" {
			values[strings.Trim(_softwareClientLicense.Key2, " ")] = strings.Trim(_softwareClientLicense.Value2, " ")
		}
		if strings.Trim(_softwareClientLicense.Key3, " ") != "" {
			values[strings.Trim(_softwareClientLicense.Key3, " ")] = strings.Trim(_softwareClientLicense.Value3, " ")
		}
		if strings.Trim(_softwareClientLicense.Key4, " ") != "" {
			values[strings.Trim(_softwareClientLicense.Key4, " ")] = strings.Trim(_softwareClientLicense.Value4, " ")
		}
		if strings.Trim(_softwareClientLicense.Key5, " ") != "" {
			values[strings.Trim(_softwareClientLicense.Key5, " ")] = strings.Trim(_softwareClientLicense.Value5, " ")
		}
		if strings.Trim(_softwareClientLicense.Key6, " ") != "" {
			values[strings.Trim(_softwareClientLicense.Key6, " ")] = strings.Trim(_softwareClientLicense.Value6, " ")
		}
		if strings.Trim(_softwareClientLicense.Key7, " ") != "" {
			values[strings.Trim(_softwareClientLicense.Key7, " ")] = strings.Trim(_softwareClientLicense.Value7, " ")
		}

		if config.RequiredKeys != nil {
			missingKeys := []string{}
			for _, key := range *config.RequiredKeys {
				_, ok := values[key]
				if !ok {
					missingKeys = append(missingKeys, key)
				}
			}

			if len(missingKeys) > 0 {
				fmt.Printf("%s skipped (missing keys: %s)\n", _softwareClientLicense.ObjectName, strings.Join(missingKeys, ", "))
				continue
			}
		}

		softwareClientLicenses = append(softwareClientLicenses, SoftwareClientLicense{
			SoftwareClientLicenseGuid:  _softwareClientLicense.SoftwareClientLicenseGuid,
			ObjectID:                   _softwareClientLicense.ObjectID,
			ObjectName:                 _softwareClientLicense.ObjectName,
			SoftwareClientLicenseState: _softwareClientLicense.SoftwareClientLicenseState,
			DataState:                  _softwareClientLicense.DataState,
			AsOf:                       *_softwareClientLicense.AsOf,
			Values:                     values,
		})
	}

	return &softwareClientLicenses, nil
}
