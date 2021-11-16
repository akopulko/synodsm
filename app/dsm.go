package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/dustin/go-humanize"
	"github.com/go-resty/resty/v2"
)

type SynologyApiCallData struct {
	Server string
	Path   string
	Params map[string]string
}

func callSynologyApi(apiData SynologyApiCallData) ([]byte, error) {

	// build url for api call
	apiCallUrl, err := url.Parse(fmt.Sprintf("%s/webapi/%s", apiData.Server, apiData.Path))
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	for key, value := range apiData.Params {
		params.Add(key, value)
	}
	apiCallUrl.RawQuery = params.Encode()

	//fmt.Println(apiCallUrl.String())

	// make api call via resty
	client := resty.New()
	resp, err := client.R().Get(apiCallUrl.String())
	if err != nil {
		return nil, err
	}

	//fmt.Println(resp)

	return resp.Body(), nil
}

func getSynologyApiInfo(server string, api string) (string, string, error) {

	apiRequest := SynologyApiCallData{
		Server: server,
		Path:   "query.cgi",
		Params: map[string]string{
			"api":     "SYNO.API.Info",
			"version": "1",
			"method":  "query",
			"query":   api,
		},
	}

	apiResponse := struct {
		Data map[string]struct {
			Path       string `json:"path,omitempty"`
			MinVersion int    `json:"minVersion,omitempty"`
			MaxVersion int    `json:"maxVersion,omitempty"`
		} `json:"data,omitempty"`
		Error struct {
			Code int `json:"code,omitempty"`
		} `json:"error,omitempty"`
		Success bool `json:"success,omitempty"`
	}{}

	resp, err := callSynologyApi(apiRequest)
	if err != nil {
		return "", "", err
	}

	err = json.Unmarshal(resp, &apiResponse)
	if err != nil {
		return "", "", err
	}

	if !apiResponse.Success {
		return "", "", fmt.Errorf("SYNO.API.Info Error %d", apiResponse.Error.Code)
	}

	return apiResponse.Data[api].Path, fmt.Sprintf("%d", apiResponse.Data[api].MaxVersion), nil

}

func Login(server string, user string, password string) (string, error) {

	api := "SYNO.API.Auth"
	path, version, err := getSynologyApiInfo(server, api)
	if err != nil {
		return "", err
	}

	apiRequest := SynologyApiCallData{
		Server: server,
		Path:   path,
		Params: map[string]string{
			"api":     api,
			"version": version,
			"method":  "login",
			"account": user,
			"passwd":  password,
			"session": "DownloadStation",
			"format":  "sid",
		},
	}

	resp, err := callSynologyApi(apiRequest)
	if err != nil {
		return "", err
	}

	// parse login response
	apiResponse := struct {
		Data struct {
			Sid string `json:"sid,omitempty"`
		} `json:"data,omitempty"`
		Error struct {
			Code int `json:"code,omitempty"`
		} `json:"error,omitempty"`
		Success bool `json:"success,omitempty"`
	}{}

	err = json.Unmarshal(resp, &apiResponse)

	if err != nil {
		return "", err
	}

	if !apiResponse.Success {
		return "", fmt.Errorf("%s Error %d", api, apiResponse.Error.Code)
	}

	return apiResponse.Data.Sid, nil

}

func Logout(server string) error {

	api := "SYNO.API.Auth"
	path, version, err := getSynologyApiInfo(server, api)
	if err != nil {
		return err
	}

	apiRequest := SynologyApiCallData{
		Server: server,
		Path:   path,
		Params: map[string]string{
			"api":     api,
			"version": version,
			"method":  "logout",
			"session": "DownloadStation",
		},
	}

	resp, err := callSynologyApi(apiRequest)
	if err != nil {
		return err
	}

	// parse login response
	apiResponse := struct {
		Error struct {
			Code int `json:"code,omitempty"`
		} `json:"error,omitempty"`
		Success bool `json:"success,omitempty"`
	}{}

	err = json.Unmarshal(resp, &apiResponse)

	if err != nil {
		return err
	}

	if !apiResponse.Success {
		return fmt.Errorf("%s Error %d", api, apiResponse.Error)
	}

	return nil

}

func listTorrents(server string, sid string) ([]TorrentTask, error) {

	// get api version and path
	api := "SYNO.DownloadStation.Task"
	path, version, err := getSynologyApiInfo(server, api)
	if err != nil {
		return nil, err
	}

	apiRequest := SynologyApiCallData{
		Server: server,
		Path:   path,
		Params: map[string]string{
			"api":        api,
			"version":    version,
			"method":     "list",
			"additional": "transfer",
			"_sid":       sid,
		},
	}

	resp, err := callSynologyApi(apiRequest)
	if err != nil {
		return nil, err
	}

	//parse tasks list response
	apiResponse := struct {
		Error struct {
			Code int `json:"code,omitempty"`
		} `json:"error,omitempty"`
		Success bool `json:"success,omitempty"`
		Data    struct {
			Tasks []struct {
				Id         string `json:"id,omitempty"`
				Type       string `json:"type,omitempty"`
				Username   string `json:"username,omitempty"`
				Title      string `json:"title,omitempty"`
				Size       int64  `json:"size,omitempty"`
				Status     string `json:"status,omitempty"`
				Additional struct {
					Transfer struct {
						SizeDownloaded int64 `json:"size_downloaded"`
					} `json:"transfer,omitempty"`
				} `json:"additional,omitempty"`
			}
		}
	}{}

	err = json.Unmarshal(resp, &apiResponse)

	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("%s Error %d", api, apiResponse.Error)
	}

	// build result array of tasks informaton
	var arrayTorrentTask []TorrentTask
	for _, task := range apiResponse.Data.Tasks {
		if task.Type == "bt" {
			//progress := ((float64(task.Additional.Transfer.SizeDownloaded-task.Size) / float64(task.Size)) * 100) + 100
			progress := calculatePercentage(task.Additional.Transfer.SizeDownloaded, task.Size)
			torrentTask := TorrentTask{
				Id:       task.Id,
				Title:    task.Title,
				User:     task.Username,
				Size:     humanize.Bytes(uint64(task.Size)),
				Status:   task.Status,
				Progress: fmt.Sprintf("%.0f %%", progress),
			}
			arrayTorrentTask = append(arrayTorrentTask, torrentTask)
		}
	}

	return arrayTorrentTask, nil

}

func addTorrentFromUri(server string, uri string, sid string) error {

	// get api version and path
	api := "SYNO.DownloadStation.Task"
	path, version, err := getSynologyApiInfo(server, api)
	if err != nil {
		return err
	}

	apiRequest := SynologyApiCallData{
		Server: server,
		Path:   path,
		Params: map[string]string{
			"api":     api,
			"version": version,
			"method":  "create",
			"uri":     uri,
			"_sid":    sid,
		},
	}

	resp, err := callSynologyApi(apiRequest)
	if err != nil {
		return err
	}

	// parse  response
	apiResponse := struct {
		Error struct {
			Code int `json:"code,omitempty"`
		} `json:"error,omitempty"`
		Success bool `json:"success,omitempty"`
	}{}

	err = json.Unmarshal(resp, &apiResponse)

	if err != nil {
		return err
	}

	if !apiResponse.Success {
		return fmt.Errorf("%s Error %d", api, apiResponse.Error)
	}

	return nil

}

func manageTorrentTask(server string, taskId string, action string, sid string) error {

	if !isValidAction(action) {
		return fmt.Errorf("provided action '%s' is invalid", action)
	}

	api := "SYNO.DownloadStation.Task"
	path, version, err := getSynologyApiInfo(server, api)
	if err != nil {
		return err
	}

	apiRequest := SynologyApiCallData{
		Server: server,
		Path:   path,
		Params: map[string]string{
			"api":     api,
			"version": version,
			"method":  action,
			"id":      taskId,
			"_sid":    sid,
		},
	}

	resp, err := callSynologyApi(apiRequest)
	if err != nil {
		return err
	}

	apiResponse := struct {
		Error struct {
			Code int `json:"code,omitempty"`
		} `json:"error,omitempty"`
		Success bool `json:"success,omitempty"`
		Data    []struct {
			Error int    `json:"error,omitempty"`
			Id    string `json:"id,omitempty"`
		}
	}{}

	err = json.Unmarshal(resp, &apiResponse)

	if err != nil {
		return err
	}

	if !apiResponse.Success {
		return fmt.Errorf("%s Error %d", api, apiResponse.Error)
	}

	if apiResponse.Data[0].Error != 0 {
		return fmt.Errorf("%s Error %d", api, apiResponse.Data[0].Error)
	}

	return nil

}
