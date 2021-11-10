package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/docopt/docopt-go"
)

var usageTestTable = []struct {
	argv      []string    // Given command line args
	validArgs bool        // Are they supposed to be valid?
	opts      docopt.Opts // Expected options parsed
}{
	{
		[]string{"init", "http://server", "user", "password"},
		true,
		docopt.Opts{
			"init":          true,
			"list":          false,
			"remove":        false,
			"add":           false,
			"<server>":      "http://server",
			"<user>":        "user",
			"<password>":    "password",
			"<task_id>":     nil,
			"<torrent_url>": nil,
		},
	},
	{
		[]string{"list"},
		true,
		docopt.Opts{
			"init":          false,
			"list":          true,
			"remove":        false,
			"add":           false,
			"<server>":      nil,
			"<user>":        nil,
			"<password>":    nil,
			"<task_id>":     nil,
			"<torrent_url>": nil,
		},
	},
	{
		[]string{"list", "dfdfd"},
		false,
		docopt.Opts{},
	},
	{
		[]string{"init", "dfdfd"},
		false,
		docopt.Opts{},
	},
}

func TestUsage(t *testing.T) {
	for _, tt := range usageTestTable {
		validArgs := true
		parser := &docopt.Parser{
			HelpHandler: func(err error, usage string) {
				if err != nil {
					validArgs = false // Triggered usage, args were invalid.
				}
			},
		}
		opts, err := parser.ParseArgs(usage, tt.argv, "")
		fmt.Println(opts)
		if validArgs != tt.validArgs {
			t.Fail()
		}
		if tt.validArgs && err != nil {
			t.Fail()
		}
		if tt.validArgs && !reflect.DeepEqual(opts, tt.opts) {
			t.Errorf("result (1) doesn't match expected (2) \n%v \n%v", opts, tt.opts)
		}
	}
}
