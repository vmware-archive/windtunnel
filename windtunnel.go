package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"

	"github.com/codegangsta/cli"
)

func kill(c *cli.Context) {
	endpoint := c.String("e")
	requests := c.Int("r")

	done := make(chan bool)

	for i := 0; i < requests; i++ {
		fmt.Printf("Killing instance %v...\n", i)

		// Use goroutines to kill instances in parallel!
		go func() {
			resp, err := http.Get("http://" + endpoint + "/killSwitch")
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
			} else {
				fmt.Printf("Response: %v\n", resp.Status)
			}
			done <- true
		}()
	}

	// Allow all goroutines to finish executing.
	for i := 0; i < requests; i++ {
		<-done
	}
}

func health(c *cli.Context) {
	endpoint := c.String("e")
	requests := c.Int("r")

	done := make(chan bool)
	var healthy uint64 = 0

	for i := 0; i < requests; i++ {
		go func() {
			resp, err := http.Get("http://" + endpoint + "/health")
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
			} else {
				fmt.Printf("Response: %v\n", resp.Status)

				if resp.StatusCode == 200 {
					atomic.AddUint64(&healthy, 1)
					runtime.Gosched()
				}
			}
			done <- true
		}()
	}

	for i := 0; i < requests; i++ {
		<-done
	}
	fmt.Printf("Healthy Requests: %v\n", atomic.LoadUint64(&healthy))
}

func main() {
	app := cli.NewApp()
	app.Name = "windtunnel"
	app.Usage = "A tool for stress testing cloud application platforms."

	app.Commands = []cli.Command{
		{
			Name:      "kill",
			ShortName: "k",
			Usage:     "hit the kill endpoint n times",
			Action:    kill,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "e",
					Usage: "endpoint that will kill the app",
				},
				cli.IntFlag{
					Name:  "r",
					Value: 1,
					Usage: "number of requests to submit",
				},
			},
		},
		{
			Name:      "health",
			ShortName: "h",
			Usage:     "poll the health endpoint n times",
			Action:    health,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "e",
					Usage: "endpoint that will provide a health indicator",
				},
				cli.IntFlag{
					Name:  "r",
					Value: 1,
					Usage: "number of requests to submit",
				},
			},
		},
	}

	app.Run(os.Args)
}
