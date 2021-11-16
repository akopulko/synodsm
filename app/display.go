package main

import (
	"fmt"

	"github.com/bndr/gotabulate"
)

func printTorrentTasks(tasks []TorrentTask) {

	if len(tasks) > 0 {
		tabulateTasks := makeTorrentTasksTabulate(tasks)
		tabulateOutput := gotabulate.Create(tabulateTasks)
		tabulateOutput.SetHeaders([]string{"ID", "User", "Title", "Status", "Size (GB)", "Progress (%)"})
		tabulateOutput.SetAlign("left")
		fmt.Println(tabulateOutput.Render("grid"))
	} else {
		fmt.Println("No tasks found")
	}

}
