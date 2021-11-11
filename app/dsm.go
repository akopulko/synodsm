package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type SynoApiInfoData struct {
	Path       string
	MinVersion int
	MaxVersion int
}

type SynoApiInfoResponse struct {
	Data    map[string]SynoApiInfoData
	Error   int
	Success bool
}

func synoApiInfo(server string, api string) (*SynoApiInfoData, error) {

	apiInfoData := SynoApiInfoData{}

	url := fmt.Sprintf(
		"%s/webapi/query.cgi?api=SYNO.API.Info&version=1&method=query&query=%s",
		server,
		api,
	)

	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return &apiInfoData, err
	}

	var apiResponse SynoApiInfoResponse
	err = json.Unmarshal([]byte(resp.Body()), &apiResponse)

	if err != nil {
		return &apiInfoData, err
	}

	if !apiResponse.Success {
		return &apiInfoData, fmt.Errorf("SYNO.API.Info Error %d", apiResponse.Error)
	}

	apiInfoData = apiResponse.Data[api]

	return &apiInfoData, nil

}

func synoLogin(server string, user string, password string) (string, error) {

	// get api version and path
	api := "SYNO.API.Auth"
	apiInfo, err := synoApiInfo(server, api)
	if err != nil {
		return "", err
	}

	// build login url
	url := fmt.Sprintf(
		"%s/webapi/%s?api=%s&version=%d&method=login&account=%s&passwd=%s&session=DownloadStation&format=sid",
		server,
		apiInfo.Path,
		api,
		apiInfo.MaxVersion,
		user,
		password,
	)

	// sent login request
	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return "", err
	}

	// parse login response
	apiResponse := struct {
		Data    struct{ Sid string }
		Error   int
		Success bool
	}{}

	err = json.Unmarshal([]byte(resp.Body()), &apiResponse)

	if err != nil {
		return "", err
	}

	if !apiResponse.Success {
		return "", fmt.Errorf("%s Error %d", api, apiResponse.Error)
	}

	return apiResponse.Data.Sid, nil

}

func synoLogout(server string) error {

	// get api version and path
	api := "SYNO.API.Auth"
	apiInfo, err := synoApiInfo(server, api)
	if err != nil {
		return err
	}

	// build login url
	url := fmt.Sprintf(
		"%s/webapi/%s?api=%s&version=%d&method=logout&session=DownloadStation",
		server,
		apiInfo.Path,
		api,
		apiInfo.MaxVersion,
	)

	// sent login request
	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return err
	}

	// parse login response
	apiResponse := struct {
		Error   int
		Success bool
	}{}

	err = json.Unmarshal([]byte(resp.Body()), &apiResponse)

	if err != nil {
		return err
	}

	if !apiResponse.Success {
		return fmt.Errorf("%s Error %d", api, apiResponse.Error)
	}

	return nil

}
