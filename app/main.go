package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/docopt/docopt-go"
)

var usage = `
Usage:
	synodsm init <server> <user> <password>
	synodsm list
	synodsm remove <task_id>
	synodsm add-uri <torrent_file_uri>
	synodsm pause <task_id>
	synodsm resume <task_id>


`

type CmdLineArgs struct {
	Init       bool   `docopt:"init"`
	List       bool   `docopt:"list"`
	Remove     bool   `docopt:"remove"`
	AddUri     bool   `docopt:"add-uri"`
	Pause      bool   `docopt:"pause"`
	Resume     bool   `docopt:"resume"`
	Server     string `docopt:"<server>"`
	User       string `docopt:"<user>"`
	Password   string `docopt:"<password>"`
	TaskID     string `docopt:"<task_id>"`
	TorrentUri string `docopt:"<torrent_file_uri>"`
}

func main() {

	// parse docopt
	args, err := docopt.ParseDoc(usage)

	if err != nil {
		log.Fatalln(err)
		return
	}

	// bind cmd line arguments
	var cmdLineArgs CmdLineArgs
	err = args.Bind(&cmdLineArgs)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// handle init command
	if cmdLineArgs.Init {

		configFilePath, err := getUserConfigFilePath()
		if err != nil {
			log.Fatalln(err)
		}

		err = saveConfig(cmdLineArgs.Server, cmdLineArgs.User, configFilePath)
		if err != nil {
			log.Fatalln(err)
		}

		err = saveCredentials(cmdLineArgs.User, cmdLineArgs.Password)
		if err != nil {
			log.Fatalln(err)
		}

	} else {

		config, err := loadConfig()
		if err != nil {
			log.Fatalln(err)
		}

		password, err := getPassword(config.User)
		if err != nil {
			log.Fatalln(err)
		}

		sid, err := Login(config.Server, config.User, password)
		if err != nil {
			log.Fatalln(err)
		}

		// handle other commands
		if cmdLineArgs.List {
			torrentsList, err := listTorrents(config.Server, sid)
			if err != nil {
				log.Fatalln(err)
			}
			printTorrentTasks(torrentsList)
		}

		if cmdLineArgs.AddUri {
			// check valid uri
			_, err := url.ParseRequestURI(cmdLineArgs.TorrentUri)
			if err != nil {
				log.Fatalln(err)
			}
			// check uri live
			err = testUri(cmdLineArgs.TorrentUri)
			if err != nil {
				log.Fatalln(err)
			}
			err = addTorrentFromUri(config.Server, cmdLineArgs.TorrentUri, sid)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("Task successfuly added to Download Station")
		}

		if cmdLineArgs.Remove {
			err = manageTorrentTask(config.Server, cmdLineArgs.TaskID, "delete", sid)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("Task %s successfuly removed\n", cmdLineArgs.TaskID)
		}

		if cmdLineArgs.Pause {
			err = manageTorrentTask(config.Server, cmdLineArgs.TaskID, "pause", sid)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("Task %s successfuly paused\n", cmdLineArgs.TaskID)
		}

		if cmdLineArgs.Resume {
			err = manageTorrentTask(config.Server, cmdLineArgs.TaskID, "resume", sid)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("Task %s successfuly resumed\n", cmdLineArgs.TaskID)
		}

		err = Logout(config.Server)
		if err != nil {
			log.Fatalln(err)
		}
	}

}
