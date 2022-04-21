package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	rig "./package"
)

func do_job() {
	log.Println("do some job")
	rig.Update()

	// rig_205a = HIVERIG("w0005", "192.168.0.205", "user", "9UNXmhyV")
	// rig_207b = HIVERIG("w0007", "192.168.0.207", "user", "1")
	// rig_204c = HIVERIG("w0004", "192.168.0.204", "user", "1")

	// miner := claymore.Miner{Address: "w0004:3334"}
	// info, err := miner.GetInfo()

	// if err != nil {
	// 	//log.Fatal(err)
	// 	fmt.Println(err)
	// }

	// json_bytes, err := json.Marshal(info)
	// if err != nil {
	// 	//log.Fatal(err)
	// 	fmt.Println(err)
	// }

	// fmt.Println(string(json_bytes))
}

func main() {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range sigs {
			log.Printf("captured %v, stopping profiler and exiting..", sig)
			done <- true
		}
	}()

loop:
	for {
		select {
		case <-time.After(time.Millisecond * 1000):
			do_job()
			fmt.Println("after some time, do job")

		case <-done:
			fmt.Println("exting")
			break loop
		}
	}

}
