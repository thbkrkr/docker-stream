package main

import (
	"encoding/json"
	"flag"
	"os"
	"strings"
	"time"

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
	hostname  string

	docker *client.Client
	kafka  *qlient.Qlient
	pub    chan []byte

	noExecEvent    bool
	defaultHeaders = map[string]string{"User-Agent": name + "/" + gitCommit}
	tick           = time.Duration(1) * time.Minute
)

type Container struct {
	Timestamp int64
	Host      string
	ID        string
	Name      string
	Image     string
	Status    string
}

func main() {
	flag.BoolVar(&noExecEvent, "no-exec-event", true, "Filter exec_* event")
	flag.Parse()

	hostname = os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname = os.Getenv("HOST")
	}
	if hostname == "" {
		hostname = "unknown"
	}

	kafka, err := qlient.NewClientFromEnv(name)
	if err != nil {
		log.Fatal(err)
	}

	pub, err = kafka.AsyncPub()
	if err != nil {
		log.Fatal(err)
	}

	ps()
	go func() {
		tick := time.NewTicker(tick)
		for range tick.C {
			ps()
		}
	}()

	stream()
}

func ps() {
	containers, err := list()
	if err != nil {
		log.Error(err)
	} else {
		push(containers)
	}
}

func list() ([]types.Container, error) {
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

func push(containers []types.Container) {
	for _, container := range containers {
		log.Infof("%+v", container)
		msg, _ := json.Marshal(Container{
			Timestamp: time.Now().Unix(),
			Host:      hostname,
			ID:        container.ID,
			Image:     container.Image,
			Name:      container.Names[0],
			Status:    container.State,
		})
		pub <- []byte(msg)
	}
}

func stream() {
	events, errs := docker.Events(context.Background(), types.EventsOptions{})

	go func() {
		for err := range errs {
			log.Error(err)
		}
	}()

	for event := range events {
		// Discard exec event
		if noExecEvent && strings.HasPrefix(event.Status, "exec_") {
			continue
		}

		// Discard event not relative to containers
		if event.Type != "container" {
			continue
		}

		msg, _ := json.Marshal(Container{
			Timestamp: time.Now().Unix(),
			Host:      hostname,
			ID:        event.ID,
			Image:     event.Actor.Attributes["image"],
			Name:      event.Actor.Attributes["name"],
			Status:    event.Status,
		})
		pub <- []byte(msg)
	}
}
