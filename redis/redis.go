package redis

import (
	"context"
	"fmt"
	"os/exec"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

var (
	cli           *client.Client
	once          sync.Once
	clockJumpTime = "+10"
)

type Instance struct {
	NodeID    int
	Container types.Container
}

func init() {
	once.Do(func() {
		if cli == nil {
			var err error
			cli, err = client.NewClientWithOpts(
				client.FromEnv,
				client.WithAPIVersionNegotiation(),
			)
			if err != nil {
				panic(err)
			}
		}
	})
}

func GetInstances() []*Instance {
	ctx := context.Background()

	filters := filters.NewArgs()
	filters.Add("status", "running")
	filters.Add("ancestor", "redis:6.0.18-alpine3.17")

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters,
	})
	if err != nil {
		panic(err)
	}

	Instances := make([]*Instance, len(containers))
	for i := range containers {
		Instances[i] = &Instance{
			NodeID:    i,
			Container: containers[len(containers)-1-i],
		}
	}
	return Instances
}

func (instance *Instance) Start() {
	ctx := context.Background()
	if err := cli.ContainerStart(ctx, instance.Container.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

func (instance *Instance) Restart() {
	ctx := context.Background()
	if err := cli.ContainerRestart(ctx, instance.Container.ID, container.StopOptions{}); err != nil {
		panic(err)
	}
}

func (instance *Instance) Crash() {
	ctx := context.Background()
	if err := cli.ContainerStop(ctx, instance.Container.ID, container.StopOptions{}); err != nil {
		panic(err)
	}
}

func (instance *Instance) ClockJump() {
	ctx := context.Background()
	containerJson, err := cli.ContainerInspect(ctx, instance.Container.ID)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("sudo", "./watchmaker", "-pid", fmt.Sprintf("%d", containerJson.State.Pid), "-sec_delta", clockJumpTime)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func (instance *Instance) Cleanup() {
	ctx := context.Background()
	defer func() {
		if err := cli.Close(); err != nil {
			panic(err)
		}
	}()
	if err := cli.ContainerStop(ctx, instance.Container.ID, container.StopOptions{}); err != nil {
		panic(err)
	}
}
