package main

import (
	"encoding/json"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/go-resty/resty/v2"
)

func synoApiInfo(server string, api string) (string, int, error) {

	url := fmt.Sprintf(
		"%s/webapi/query.cgi?api=SYNO.API.Info&version=1&method=query&query=%s",
		server,
		api,
	)

	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return "", 0, err
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

	err = json.Unmarshal([]byte(resp.Body()), &apiResponse)

	if err != nil {
		return "", 0, err
	}

	if !apiResponse.Success {
		return "", 0, fmt.Errorf("SYNO.API.Info Error %d", apiResponse.Error.Code)
	}

	return apiResponse.Data[api].Path, apiResponse.Data[api].MaxVersion, nil

}

func synoLogin(server string, user string, password string) (string, error) {

	// get api version and path
	api := "SYNO.API.Auth"
	//apiInfo, err := synoApiInfo(server, api)
	path, version, err := synoApiInfo(server, api)
	if err != nil {
		return "", err
	}

	// build login url
	url := fmt.Sprintf(
		"%s/webapi/%s?api=%s&version=%d&method=login&account=%s&passwd=%s&session=DownloadStation&format=sid",
		server,
		path,
		api,
		version,
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
		Data struct {
			Sid string `json:"sid,omitempty"`
		} `json:"data,omitempty"`
		Error struct {
			Code int `json:"code,omitempty"`
		} `json:"error,omitempty"`
		Success bool `json:"success,omitempty"`
	}{}

	err = json.Unmarshal([]byte(resp.Body()), &apiResponse)

	if err != nil {
		return "", err
	}

	if !apiResponse.Success {
		return "", fmt.Errorf("%s Error %d", api, apiResponse.Error.Code)
	}

	return apiResponse.Data.Sid, nil

}

func synoLogout(server string) error {

	// get api version and path
	api := "SYNO.API.Auth"
	//apiInfo, err := synoApiInfo(server, api)
	path, version, err := synoApiInfo(server, api)

	if err != nil {
		return err
	}

	// build login url
	url := fmt.Sprintf(
		"%s/webapi/%s?api=%s&version=%d&method=logout&session=DownloadStation",
		server,
		path, //apiInfo.Path,
		api,
		version, //apiInfo.MaxVersion,
	)

	// sent login request
	client := resty.New()
	resp, err := client.R().Get(url)
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

	err = json.Unmarshal([]byte(resp.Body()), &apiResponse)

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
	//apiInfo, err := synoApiInfo(server, api)
	path, version, err := synoApiInfo(server, api)

	if err != nil {
		return nil, err
	}

	// build login url
	url := fmt.Sprintf(
		"%s/webapi/%s?api=%s&version=%d&method=list&additional=transfer&_sid=%s",
		server,
		path, //apiInfo.Path,
		api,
		version, //apiInfo.MaxVersion,
		sid,
	)

	// sent login request
	client := resty.New()
	resp, err := client.R().Get(url)
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

	err = json.Unmarshal([]byte(resp.Body()), &apiResponse)

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
	path, version, err := synoApiInfo(server, api)

	if err != nil {
		return err
	}

	// build url
	url := fmt.Sprintf(
		"%s/webapi/%s?api=%s&version=%d&method=create&uri=%s&_sid=%s",
		server,
		path,
		api,
		version,
		uri,
		sid,
	)

	// send request
	client := resty.New()
	resp, err := client.R().Get(url)
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

	err = json.Unmarshal([]byte(resp.Body()), &apiResponse)

	if err != nil {
		return err
	}

	if !apiResponse.Success {
		return fmt.Errorf("%s Error %d", api, apiResponse.Error)
	}

	return nil

}
