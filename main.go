package main

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	qlient "github.com/thbkrkr/qli/client"
	"golang.org/x/net/context"
)

var (
	name      = "docker-stream"
	buildDate = "dev"
	gitCommit = "dev"

	docker *client.Client
	kafka  *qlient.Qlient

	defaultHeaders = map[string]string{"User-Agent": name + "/" + gitCommit}
)

type Container struct {
	ID     string
	Name   string
	Image  string
	Status string
}

func main() {
	kafka, err := qlient.NewClientFromEnv(name)
	if err != nil {
		log.Fatal(err)
	}

	pub, err := kafka.AsyncPub()
	if err != nil {
		log.Fatal(err)
	}

	containers, err := listContainers()
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range containers {
		//log.Infof("container:\n %+v", container)
		log.Info(container.ID)
		log.Info(container.Names[0])
		log.Info(container.Image)
		log.Info(container.State)
		log.Info(container.Status)
		msg, _ := json.Marshal(Container{
			ID:     container.ID,
			Image:  container.Image,
			Name:   container.Names[0],
			Status: container.Status,
		})
		pub <- []byte(msg)
	}

	events, errs := docker.Events(context.Background(), types.EventsOptions{})
	go func() {
		for err := range errs {
			log.Error(err)
		}
	}()
	for event := range events {
		//log.Infof("%+v", event)
		//log.Info(event.Type)
		if event.Type == "container" {
			log.Info(event.ID)
			log.Info(event.Actor.Attributes["name"])
			log.Info(event.Actor.Attributes["image"])
			log.Info(event.Status)
			log.Info(event.Actor.Attributes["commit"])

			msg, _ := json.Marshal(Container{
				ID:     event.ID,
				Image:  event.Actor.Attributes["image"],
				Name:   event.Actor.Attributes["name"],
				Status: event.Status,
			})
			pub <- []byte(msg)
		}
	}
}

func listContainers() ([]types.Container, error) {
	if docker == nil {
		c, err := client.NewClient("unix:///var/run/docker.sock", "v1.30", nil, defaultHeaders)
		if err != nil {
			return nil, err
		}
		docker = c
	}

	containers, err := docker.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	return containers, nil
}
