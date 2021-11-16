package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type TorrentTask struct {
	Id       string
	User     string
	Title    string
	Size     string
	Status   string
	Progress string
}

func makeTorrentTasksTabulate(tasks []TorrentTask) [][]string {

	var tabulateTasks [][]string
	for _, taskItem := range tasks {
		strTask := []string{
			taskItem.Id,
			taskItem.User,
			taskItem.Title,
			taskItem.Status,
			taskItem.Size,
			taskItem.Progress,
		}
		tabulateTasks = append(tabulateTasks, strTask)
	}

	return tabulateTasks

}

func calculatePercentage(current int64, total int64) float64 {
	percentage := ((float64(current-total) / float64(total)) * 100) + 100
	return percentage

}

func testUri(uri string) error {

	client := resty.New()
	resp, err := client.R().Get(uri)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("uri returned http error %d [%s]", resp.StatusCode(), uri)
	}

	return nil
}
