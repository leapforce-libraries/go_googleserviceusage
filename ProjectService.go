package googleserviceusage

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
)

type ProjectServicesResponse struct {
	ProjectServices []ProjectService `json:"projectServices"`
	NextPageToken   string           `json:"nextPageToken"`
}

type ProjectService struct {
	Name   string `json:"name"`
	Config struct {
		Name          string `json:"name"`
		Title         string `json:"title"`
		Documentation struct {
			Summary string `json:"summary"`
		} `json:"documentation"`
		Quota struct {
		} `json:"quota"`
		Authentication struct {
		} `json:"authentication"`
		Usage struct {
			Requirements []string `json:"requirements"`
		} `json:"usage"`
		MonitoredResources []struct {
			Type        string `json:"type"`
			DisplayName string `json:"displayName"`
			Description string `json:"description"`
			Labels      []struct {
				Key         string `json:"key"`
				Description string `json:"description"`
			} `json:"labels"`
			LaunchStage string `json:"launchStage"`
		} `json:"monitoredResources"`
		Monitoring struct {
			ConsumerDestinations []struct {
				MonitoredResource string   `json:"monitoredResource"`
				Metrics           []string `json:"metrics"`
			} `json:"consumerDestinations"`
		} `json:"monitoring"`
	} `json:"config"`
	State  string `json:"state"`
	Parent string `json:"parent"`
}

type ProjectServicesConfig struct {
	ProjectNumber string
	Filter        *string
}

func (service *Service) ProjectServices(config *ProjectServicesConfig) (*[]ProjectService, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("Config must not be nil")
	}
	var projectServices []ProjectService

	var pageToken = ""

	var values = url.Values{}
	if config != nil {
		if config.Filter != nil {
			values.Set("filter", *config.Filter)
		}
	}

	for {
		if pageToken != "" {
			values.Set("pageToken", pageToken)
		}

		var projectServicesResponse ProjectServicesResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("projects/%s/services?%s", config.ProjectNumber, values.Encode())),
			ResponseModel: &projectServicesResponse,
		}

		_, _, e := service.googleService.HttpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		projectServices = append(projectServices, projectServicesResponse.ProjectServices...)

		if projectServicesResponse.NextPageToken == "" {
			break
		}

		pageToken = projectServicesResponse.NextPageToken
	}

	return &projectServices, nil
}
