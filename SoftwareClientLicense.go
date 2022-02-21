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

type HistoricDataState string

const (
	HistoricDataStateNone     HistoricDataState = "None"
	HistoricDataStateStart    HistoricDataState = "Start"
	HistoricDataStateRunning  HistoricDataState = "Running..."
	HistoricDataStateDone     HistoricDataState = "Done"
	HistoricDataStateDelete   HistoricDataState = "Delete"
	HistoricDataStateDeleting HistoricDataState = "Deleting..."
	HistoricDataStateDeleted  HistoricDataState = "Deleted"

	StateActive          string = "Active"
	CustomFieldGuidState string = "c6986444-5147-4f25-97e8-6857aa4a3189"
)

type SoftwareClientLicense struct {
	CompanyId                 int64
	CompanyName               string
	SoftwareClientLicenseGuid string
	dataState                 string
	AsOf                      civil.Date
	Values                    map[string]string
}

type GetSoftwareClientLicensesConfig struct {
	SoftwarePackageGuid       string
	RequiredKeys              *[]string
	CompanyId                 *int64
	SoftwareClientLicenseGuid *string
}

func (service *Service) GetSoftwareClientLicenses(config *GetSoftwareClientLicensesConfig) (*[]SoftwareClientLicense, *errortools.Error) {
	_softwareClientLicenses := []struct {
		CompanyId                  int64
		CompanyName                string
		SoftwareClientLicenseGuid  string
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
	values.Set("software_package_guid", config.SoftwarePackageGuid)

	if config.CompanyId != nil {
		values.Set("company_id", fmt.Sprintf("%v", *config.CompanyId))
	}
	if config.SoftwareClientLicenseGuid != nil {
		values.Set("software_client_license_guid", *config.SoftwareClientLicenseGuid)
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("creds?%s", values.Encode())),
		ResponseModel: &_softwareClientLicenses,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	softwareClientLicenses := []SoftwareClientLicense{}

	for _, _softwareClientLicense := range _softwareClientLicenses {
		if _softwareClientLicense.SoftwareClientLicenseState != StateActive {
			fmt.Printf("%s skipped (not active)\n", _softwareClientLicense.CompanyName)
			continue
		}

		/*if _softwareClientLicense.AsOf == nil {
			fmt.Printf("%s skipped (no as of date)\n", _softwareClientLicense.CompanyName)
			continue
		}*/

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
				fmt.Printf("%s skipped (missing keys: %s)\n", _softwareClientLicense.CompanyName, strings.Join(missingKeys, ", "))
				continue
			}
		}

		softwareClientLicenses = append(softwareClientLicenses, SoftwareClientLicense{
			CompanyId:                 _softwareClientLicense.CompanyId,
			CompanyName:               _softwareClientLicense.CompanyName,
			SoftwareClientLicenseGuid: _softwareClientLicense.SoftwareClientLicenseGuid,
			dataState:                 _softwareClientLicense.DataState,
			AsOf:                      *_softwareClientLicense.AsOf,
			Values:                    values,
		})
	}

	return &softwareClientLicenses, nil
}

func (softwareClientLicense *SoftwareClientLicense) DataState() string {
	return softwareClientLicense.dataState
}

func (softwareClientLicense *SoftwareClientLicense) isHistoricDataState(state HistoricDataState, default_ bool) bool {
	if softwareClientLicense == nil {
		return default_
	}

	if softwareClientLicense.dataState == "" {
		return default_
	}

	if softwareClientLicense.dataState == string(state) {
		return true
	}

	return false
}

func (softwareClientLicense *SoftwareClientLicense) IsHistoricDataNone() bool {
	return softwareClientLicense.isHistoricDataState(HistoricDataStateNone, true)
}

func (softwareClientLicense *SoftwareClientLicense) IsHistoricDataStart() bool {
	return softwareClientLicense.isHistoricDataState(HistoricDataStateStart, false)
}

func (softwareClientLicense *SoftwareClientLicense) IsHistoricDataRunning() bool {
	return softwareClientLicense.isHistoricDataState(HistoricDataStateRunning, false)
}

func (softwareClientLicense *SoftwareClientLicense) IsHistoricDataDone() bool {
	return softwareClientLicense.isHistoricDataState(HistoricDataStateDone, false)
}

func (softwareClientLicense *SoftwareClientLicense) IsHistoricDataDelete() bool {
	return softwareClientLicense.isHistoricDataState(HistoricDataStateDelete, false)
}

func (softwareClientLicense *SoftwareClientLicense) IsHistoricDataDeleting() bool {
	return softwareClientLicense.isHistoricDataState(HistoricDataStateDeleting, false)
}

func (softwareClientLicense *SoftwareClientLicense) IsHistoricDataDeleted() bool {
	return softwareClientLicense.isHistoricDataState(HistoricDataStateDeleted, false)
}

func (service *Service) setHistoricDataState(softwareClientLicense *SoftwareClientLicense, state HistoricDataState) *errortools.Error {
	body := struct {
		CompanyId                 int64  `json:"company_id"`
		SoftwareClientLicenseGuid string `json:"software_client_license_guid"`
		State                     string `json:"state"`
	}{
		softwareClientLicense.CompanyId,
		softwareClientLicense.SoftwareClientLicenseGuid,
		string(state),
	}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.url("creds/state"),
		BodyModel: body,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}

func (service *Service) SetHistoricDataRunning(softwareClientLicense *SoftwareClientLicense) *errortools.Error {
	return service.setHistoricDataState(softwareClientLicense, HistoricDataStateRunning)
}

func (service *Service) SetHistoricDataDone(softwareClientLicense *SoftwareClientLicense) *errortools.Error {
	return service.setHistoricDataState(softwareClientLicense, HistoricDataStateDone)
}

func (service *Service) SetHistoricDataDeleting(softwareClientLicense *SoftwareClientLicense) *errortools.Error {
	return service.setHistoricDataState(softwareClientLicense, HistoricDataStateDeleting)
}

func (service *Service) SetHistoricDataDeleted(softwareClientLicense *SoftwareClientLicense) *errortools.Error {
	return service.setHistoricDataState(softwareClientLicense, HistoricDataStateDeleted)
}
