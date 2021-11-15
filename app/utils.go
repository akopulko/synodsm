package main

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
