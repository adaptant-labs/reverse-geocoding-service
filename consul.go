package main

import (
	consul "github.com/hashicorp/consul/api"
	"log"
)

func ConsulServiceRegister(config *ServiceConfig) error {
	log.Println("Registering with Consul server..")

	consulConfig := consul.DefaultConfig()
	consulConfig.Address = config.consul.agentAddress

	client, err := consul.NewClient(consulConfig)
	if err != nil {
		return err
	}

	agent := client.Agent()

	service := consul.AgentServiceRegistration{
		Name:    "reverse-geocoding",
		Address: config.host,
		Port:    config.port,
	}

	return agent.ServiceRegister(&service)
}
