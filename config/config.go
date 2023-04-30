package config

import (
	"context"
	"fmt"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	redissrv "github.com/lawrshen/redlock-test/redis"
)

var (
	ports = [...]int{3001, 3002, 3003, 3004, 3005} // this is meant to be constant! Please don't mutate it!
)

type config struct {
	Pools     []redis.Pool
	Instances []*redissrv.Instance
}

func NewConfig() *config {
	cfg := &config{
		Pools:     NewPools(),
		Instances: redissrv.GetInstances(),
	}
	return cfg
}

func NewPools() []redis.Pool {
	var pools []redis.Pool
	for _, port := range ports {
		client := goredislib.NewClient(&goredislib.Options{
			Addr: fmt.Sprintf("localhost:%d", port),
		})
		pools = append(pools, goredis.NewPool(client))
	}
	return pools
}

func (cfg *config) Cleanup() {
	for _, pool := range cfg.Pools {
		ctx := context.Background()
		conn, err := pool.Get(ctx)
		if err != nil {
			panic(err)
		}
		err = conn.Close()
		if err != nil {
			panic(err)
		}
	}

	// for _, instance := range cfg.Instances {
	// 	instance.Cleanup()
	// }
}
