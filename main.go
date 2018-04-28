package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"docker.io/go-docker/api/types"
	"docker.io/go-docker"
)

var working = false
var scheduleUpdate = false

// ContainerToHosts converts a Docker Container to a hosts-file ready string.
func ContainerToHosts(container types.Container, tld string) string {
	if len(container.NetworkSettings.Networks) < 1 {
		fmt.Printf("Warning: Container %s is not connected to any network.\n", container.ID)
		return ""
	}

	hosts := ""

	for _, details := range container.NetworkSettings.Networks {
		for _, name := range container.Names {
			// <name>.docker
			hosts += details.IPAddress + "\t" + strings.TrimPrefix(name, "/") + "." + tld + "\n"
		}
		// <id>.docker
		hosts += details.IPAddress + "\t" + container.ID[:10] + "." + tld + "\n"
		hosts += details.IPAddress + "\t" + container.ID + "." + tld + "\n"
		hosts += "\n"
		break // Only the first network
	}

	return hosts
}

// Update lists all Docker Containers and adds them to the hosts file.
func Update(docker docker.APIClient, file string, tld string, wait bool) {
	working = true
	scheduleUpdate = false

	// Wait some time, especially for the "die" event (and currently only used for that) as the container won't get removed immediately from the list.
	if wait {
		time.Sleep(500 * time.Millisecond)
	}

	containers, err := docker.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	hosts := ""

	// Prepare the hosts file
	dat, err := ioutil.ReadFile(file)
	if err == nil {
		hosts += strings.SplitN(string(dat), "### Docker Hosts ###\n\n", 2)[0]
	}
	if hosts != "" {
		hosts = strings.TrimRight(hosts, "\n") + "\n\n"
	}
	hosts += "### Docker Hosts ###\n\n"

	// Add containers
	count := 0
	for _, container := range containers {
		hosts += ContainerToHosts(container, tld)
		count++
	}

	// Write the hosts file
	err = ioutil.WriteFile(file, []byte(hosts), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Updated %s, %d active containers.\n", file, count)

	working = false
	if scheduleUpdate {
		scheduleUpdate = false
		Update(docker, file, tld, false)
	}
}

// Watch calls Update() automatically every time a container changes
func Watch(docker docker.APIClient, file string, tld string) {
	go Update(docker, file, tld, false)

	msgChan, errChan := docker.Events(context.Background(), types.EventsOptions{})

	go func() {
		err := <-errChan
		if err != nil {
			panic(err)
		}
	}()

	for {
		msg := <-msgChan

		if msg.Type == "container" && (msg.Action == "die" || msg.Action == "destroy" || msg.Action == "stop" || msg.Action == "start") {
			actor := msg.Actor.ID
			if name, ok := msg.Actor.Attributes["name"]; ok {
				actor = name
			}
			fmt.Printf("Got event from container %s: %s\n", actor, msg.Action)
			if working {
				scheduleUpdate = true
			} else {
				go Update(docker, file, tld, msg.Action == "die")
			}
		}
	}
}

func main() {
	docker, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	file := "/etc/hosts"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	tld := "docker"
	if len(os.Args) > 2 {
		tld = os.Args[2]
	}

	Watch(docker, file, tld)
}
