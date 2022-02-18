package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cachance7/grapl/input"
	grapljson "github.com/cachance7/grapl/json"
	"github.com/cachance7/grapl/request"
	"github.com/cachance7/grapl/runtime"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

// TODO This is a messy function where everything just kinda came together.
// Clean it up.
func run(c *cli.Context) error {

	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
		log.Println("debug mode on")
	}

	url := c.Args().Get(0)
	if len(url) == 0 {
		return cli.Exit("url is missing", 1)
	}
	fetcher := request.NewDefaultFetcher(url)
	executor := runtime.NewAsyncRequestExecutor(fetcher)
	executor.Start()

	reader := input.NewStdinReader()
	reader.Start()

	fmt.Printf("Listening on %s\n", url)
	for {

		msg, err := reader.Read()
		if err != nil {
			panic(err)
		}

		var queryObj map[string]interface{}
		if err := json.Unmarshal(msg.Data, &queryObj); err != nil {
			log.Printf("input not json, assuming queryObj")
			queryObj = make(map[string]interface{})
			queryObj["query"] = string(msg.Data)
		}
		query, err := json.Marshal(queryObj)

		log.Println("writing query")
		log.Println(string(query))

		res, err := executor.Put(query)
		if err != nil {
			panic("could not put request")
		}

		data, err := res.Read()
		if err != nil {
			panic(err)
		}

		var obj map[string]interface{}
		json.Unmarshal(data, &obj)

		out, err := grapljson.Pretty(obj)
		if err != nil {
			log.Fatalln(err)
			break
		}

		fmt.Println("~~~~~~~~~~")
		fmt.Println(string(out))
		fmt.Println()
	}

	return nil
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Usage:   "Turns on debug logs",
				Aliases: []string{"d"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "start the app",
				Action: run,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
