package main

import (
	"fmt"

	"github.com/bndr/gotabulate"
)

func printTorrentTasks(tasks []TorrentTask) {

	tabulateTasks := makeTorrentTasksTabulate(tasks)
	tabulateOutput := gotabulate.Create(tabulateTasks)
	tabulateOutput.SetHeaders([]string{"ID", "User", "Title", "Status", "Size (GB)", "Progress (%)"})
	tabulateOutput.SetAlign("left")
	fmt.Println(tabulateOutput.Render("grid"))

}
