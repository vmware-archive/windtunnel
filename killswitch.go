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
	instances := c.Int("i")

	done := make(chan bool)

	for i := 0; i < instances; i++ {
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
	for i := 0; i < instances; i++ {
		<-done
	}
}

func health(c *cli.Context) {
	endpoint := c.String("e")
	instances := c.Int("i")

	done := make(chan bool)
	var healthy uint64 = 0

	for i := 0; i < instances; i++ {
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

	for i := 0; i < instances; i++ {
		<-done
	}
	fmt.Printf("Healthy Requests: %v\n", atomic.LoadUint64(&healthy))
}

func main() {
	app := cli.NewApp()
	app.Name = "killswitch"
	app.Usage = "Kill CF applications with extreme prejudice!"

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
					Name:  "i",
					Value: 1,
					Usage: "number of instances to kill",
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
					Usage: "endpoint that will kill the app",
				},
				cli.IntFlag{
					Name:  "i",
					Value: 1,
					Usage: "number of instances to kill",
				},
			},
		},
	}

	app.Run(os.Args)
}
