package main

import (
	"context"
	"fmt"
	// "io"
	"log"
	"os"
	// "net/http"
	// "encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	// "github.com/docker/docker/libnetwork/drivers/host"
	"github.com/docker/go-connections/nat"
	"github.com/docker/docker/pkg/stdcopy"
)


func GetAllContainer(cli *client.Client){
	ctx := context.Background()
	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, container := range containers {
		fmt.Println(container.ID[:12], container.Image)
	}
}

func GetContainer(cli *client.Client, id string){
	ctx := context.Background()
	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, container := range containers {
		containerID := container.ID[:12]
		if containerID == id{
			fmt.Println(containerID, container.Image)
		}
	}
}

func CreateContainer(cli *client.Client, docker *Container){
	//readerにはlogが格納される
	ctx := context.Background()
	pullimage := "docker.io/library/" + docker.image
	_, err := cli.ImagePull(ctx, pullimage, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	// io.Copy(os.Stdout, reader)

	if docker.hostport == ""{
		docker.hostport = "80"
	}
	config := &container.Config{
		ExposedPorts: nat.PortSet{nat.Port(docker.hostport): struct{}{}},
		Image: docker.image,
	}

	port := &container.HostConfig{	
		PortBindings: nat.PortMap{
			nat.Port(docker.hostport): []nat.PortBinding{
				{
					HostPort: docker.localport,
				},
			},
		},
	}
	resp, err := cli.ContainerCreate(context.Background(), config, port, nil, nil, "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Setting container ", resp.ID[:12], "... ")
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}
	fmt.Println("Up")
	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}

func UpdateAllContainer(cli *client.Client){
	ctx := context.Background()
	options := container.StopOptions{}
	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, container := range containers {
		fmt.Print("Restart container ", container.ID[:12], "... ")
		if err := cli.ContainerRestart(ctx, container.ID, options); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Up")
	}
}

func UpdateContainer(cli *client.Client, id string){
	ctx := context.Background()
	options := container.StopOptions{}
	fmt.Print("Restart container ", id, "... ")
		if err := cli.ContainerRestart(ctx, id, options); err != nil {
			log.Fatal(err)
		}
	fmt.Println("Up")
}

func DeleteAllContainer(cli *client.Client){
	ctx := context.Background()
	options := container.StopOptions{}
	deloptions := container.RemoveOptions{}
	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, container := range containers {
		fmt.Print("Stopping container ", container.ID[:12], "... ")
		if err := cli.ContainerStop(ctx, container.ID, options); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Success")

		fmt.Print("Deleting container ", container.ID[:12], "... ")
		if err := cli.ContainerRemove(ctx, container.ID, deloptions); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Success")
	}
}

func DeleteContainer(cli *client.Client, id string){
	ctx := context.Background()
    
    fmt.Print("Stopping container ", id, "... ")
    if err := cli.ContainerStop(ctx, id, container.StopOptions{}); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Success")

    fmt.Print("Deleting container ", id, "... ")
    if err := cli.ContainerRemove(ctx, id, container.RemoveOptions{}); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Success")
}