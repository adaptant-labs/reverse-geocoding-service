package main

import (
	"github.com/adaptant-labs/geo/georeverse"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/oschwald/maxminddb-golang"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	reverser *georeverse.CountryReverser
	geodb *maxminddb.Reader
)

type ConsulConfig struct {
	enabled		bool
	agentAddress	string
}

type ServiceConfig struct {
	consul		*ConsulConfig
	dataPath	string
	host		string
	port		int
}

func defaultServiceConfig() *ServiceConfig {
	return &ServiceConfig{
		port: 4041,
		host: "",
		dataPath: "./data/polygons.properties",
		consul: &ConsulConfig{
			enabled: true,
		},
	}
}

func main() {
	var err error

	config := defaultServiceConfig()

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "Reverse Geocoding Service"
	app.Usage = "Microservice for providing ISO 3166-1 country code lookup from a lat/lng pair"
	app.Email = "labs@adaptant.io"
	app.Version = "0.0.2"
	app.Author = "Adaptant Labs"
	app.Copyright = "Adaptant Solutions AG"

	app.Flags = []cli.Flag {
		cli.IntFlag{
			Name:		"port",
			Value:		config.port,
			Usage:		"Port to bind to",
			Destination:	&config.port,
		},

		cli.StringFlag{
			Name:		"host",
			Usage:		"Host address to bind to",
			Value:		config.host,
			Destination:	&config.host,
		},

		cli.StringFlag{
			Name:		"data",
			Usage:		"Polygon definition file to use",
			Value:		config.dataPath,
			Destination:	&config.dataPath,
		},

		cli.BoolTFlag{
			Name:		"use-consul",
			Usage:		"Use Consul for Service Registration",
			Destination:	&config.consul.enabled,
		},

		cli.StringFlag{
			Name:		"consul-agent",
			Usage:		"Consul Agent to connect to",
			Value:		"localhost:8500",
			Destination:	&config.consul.agentAddress,
		},
	}

	app.Action = func(c *cli.Context) error {
		reverser, err = georeverse.NewCountryReverser(config.dataPath)
		if err != nil {
			return err
		}

		geodb, err = maxminddb.Open("data/GeoLite2-Country.mmdb")
		if err != nil {
			return err
		}

		if config.consul.enabled == true {
			err = ConsulServiceRegister(config)
			if err != nil {
				return err
			}
		}

		router := mux.NewRouter()
		router.HandleFunc("/georeverse", georeverseHandler)
		router.HandleFunc("/georeverse/{ipAddress}", georeverseIPHandler)

		loggedRouter := handlers.LoggingHandler(os.Stdout, router)

		addr := config.host + ":" + strconv.Itoa(config.port)

		log.Printf("Listening on %s/georeverse...", addr)

		return http.ListenAndServe(addr, loggedRouter)
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
