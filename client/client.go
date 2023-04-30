package client

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-redsync/redsync/v4"
	"github.com/lawrshen/redlock-test/config"
	redissrv "github.com/lawrshen/redlock-test/redis"
)

const (
	clientNum       = 4
	abnormalNum     = 3
	finalStockValue = 1
)

func normal(instance *redissrv.Instance) {
}

func crash(instance *redissrv.Instance) {
	instance.Restart()
}

func clockjump(instance *redissrv.Instance) {
	instance.ClockJump()
}

func client(actFn func(*redissrv.Instance)) (getlockNum int) {
	cfg := config.NewConfig()
	defer cfg.Cleanup()
	rs := redsync.New(cfg.Pools...)
	mutex := rs.NewMutex("redlock-test")

	var wg sync.WaitGroup
	for i := 0; i < clientNum; i++ {
		wg.Add(1)
		go func(clientName string) {
			defer wg.Done()

			if err := mutex.Lock(); err != nil {
				log.Printf("%s get lock error: %v\n", clientName, err)
				return
			} else {
				for j := 0; j < abnormalNum; j++ {
					actFn(cfg.Instances[j])
				}
			}

			defer func() {
				if _, err := mutex.Unlock(); err != nil {
					log.Printf("%s unlock error: %v\n", clientName, err)
					return
				}
			}()

			if service() == http.StatusOK {
				getlockNum++
			}
		}(fmt.Sprintf("client[%d]", i))
	}
	wg.Wait()
	return getlockNum
}

func service() int {
	resp, err := http.Get("http://localhost:8080")

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	return resp.StatusCode
}

func reset() {
	resp, err := http.Get("http://localhost:8080/reset")

	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}
