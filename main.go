package main

import (
	"github.com/adaptant-labs/geo/georeverse"
	"github.com/gorilla/handlers"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	reverser *georeverse.CountryReverser
)

func main() {
	var err error
	var port int
	var host string
	var dataPath string

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "Reverse Geocoding Service"
	app.Usage = "Microservice for providing ISO 3166-1 country code lookup from a lat/lng pair"
	app.Email = "labs@adaptant.io"
	app.Version = "0.0.1"
	app.Author = "Adaptant Labs"

	app.Flags = []cli.Flag {
		cli.IntFlag{
			Name:			"port",
			Value:			4041,
			Usage:			"Port to bind to",
			Destination: 	&port,
		},

		cli.StringFlag{
			Name:        	"host",
			Usage:       	"Host address to bind to",
			Value:       	"",
			Destination: 	&host,
		},

		cli.StringFlag{
			Name:        	"data",
			Usage:       	"Polygon definition file to use",
			Value:       	"./data/polygons.properties",
			Destination:	&dataPath,
		},
	}

	app.Action = func(c *cli.Context) error {
		reverser, err = georeverse.NewCountryReverser(dataPath)
		if err != nil {
			log.Fatal(err)
		}

		addr := host + ":" + strconv.Itoa(port)
		log.Printf("Listening on %s/georeverse...", addr)

		http.Handle("/georeverse", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(georeverseHandler)))

		return http.ListenAndServe(addr, nil)
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}