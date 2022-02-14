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
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func repl(c *cli.Context) error {
	log.Println("repl")
	return nil
}

func run(c *cli.Context) error {
	log.Println("run")

	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
		log.Println("debug mode on")
	}

	url := c.Args().Get(0)
	fetcher := request.NewDefaultFetcher(url)
	executor := runtime.NewAsyncRequestExecutor(fetcher)
	executor.Start()

	reader := input.NewStdinReader()
	reader.Start()

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
			{
				Name:   "repl",
				Usage:  "start the repl",
				Action: repl,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
