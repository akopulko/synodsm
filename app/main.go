package main

import (
	"fmt"
	"log"

	"github.com/docopt/docopt-go"
)

var usage = `
Usage:
	synodsm init <server> <user> <password>
	synodsm list
	synodsm remove <task_id>
	synodsm add <torrent_url>
`

type CmdLineArgs struct {
	Init       bool   `docopt:"init"`
	List       bool   `docopt:"list"`
	Remove     bool   `docopt:"remove"`
	Add        bool   `docopt:"add"`
	Server     string `docopt:"<server>"`
	User       string `docopt:"<user>"`
	Password   string `docopt:"<password>"`
	TaskID     string `docopt:"<task_id>"`
	TorrentURL string `docopt:"<torrent_url>"`
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

		err = saveConfig(cmdLineArgs.Server, cmdLineArgs.User)
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

		fmt.Println(config)
		fmt.Println(password)

	}

}
