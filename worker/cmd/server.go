package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"

	"github.com/henrylee2cn/erpc/v6"

	"github.com/watercompany/multiplex/worker"
)

var portNum int

func init() {
	flag.IntVar(&portNum, "port", 9090, "port number for worker to use")
}

func main() {
	runtime.GOMAXPROCS(1)
	flag.Parse()
	deployWorker(portNum)
}

func deployWorker(portNum int) {
	erpc.SetLoggerLevel("OFF")()
	// graceful
	go erpc.GraceSignal()

	wg := new(sync.WaitGroup)

	wg.Add(1)

	go func() {
		// server peer
		srv := erpc.NewPeer(erpc.PeerConfig{
			CountTime:   true,
			ListenPort:  uint16(portNum),
			PrintDetail: false,
		})
		srv.SetTLSConfig(erpc.GenerateTLSConfigForServer())

		// router
		srv.RouteCall(new(worker.ProgramWorker))

		// listen and serve
		err := srv.ListenAndServe()
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	fmt.Printf("listen and serve: %v\n", portNum)
	wg.Wait()
}
