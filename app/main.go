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
`

type CmdLineArgs struct {
	Init       bool   `docopt:"init"`
	List       bool   `docopt:"list"`
	Remove     bool   `docopt:"remove"`
	AddUri     bool   `docopt:"add-uri"`
	AddMagnet  bool   `docopt:"add-magnet"`
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

	}

	if cmdLineArgs.List {

		config, err := loadConfig()
		if err != nil {
			log.Fatalln(err)
		}

		password, err := getPassword(config.User)
		if err != nil {
			log.Fatalln(err)
		}

		sid, err := synoLogin(config.Server, config.User, password)
		if err != nil {
			log.Fatalln(err)
		}

		torrentsList, err := listTorrents(config.Server, sid)
		if err != nil {
			log.Fatalln(err)
		}

		printTorrentTasks(torrentsList)

		err = synoLogout(config.Server)
		if err != nil {
			log.Fatalln(err)
		}

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

		config, err := loadConfig()
		if err != nil {
			log.Fatalln(err)
		}

		password, err := getPassword(config.User)
		if err != nil {
			log.Fatalln(err)
		}

		sid, err := synoLogin(config.Server, config.User, password)
		if err != nil {
			log.Fatalln(err)
		}

		err = addTorrentFromUri(config.Server, cmdLineArgs.TorrentUri, sid)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Task successfuly added to Download Station")

		err = synoLogout(config.Server)
		if err != nil {
			log.Fatalln(err)
		}

	}

}
