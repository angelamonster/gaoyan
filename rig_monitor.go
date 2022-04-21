package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	gaoyan "./gaoyan"
)

var rig_w0004 = gaoyan.RIG{ID: "w0004", IP: "192.168.0.204", Username: "user", Password: "1", ClaymorePort: 3334}
var rig_w0005 = gaoyan.RIG{ID: "w0005", IP: "192.168.0.205", Username: "user", Password: "9UNXmhyV", ClaymorePort: 3334}
var rig_w0007 = gaoyan.RIG{ID: "w0007", IP: "192.168.0.207", Username: "user", Password: "1", ClaymorePort: 3334}

var rigs = [3]gaoyan.RIG{rig_w0004, rig_w0005, rig_w0007}

func do_job() {
	log.Println("do some job")
	//gaoyan.RIG.Update()
	for _, rig := range rigs {
		json_string, err := rig.GetStat()
		if err != nil {
			panic(err)
		}
		log.Print(json_string)
		rig.PublishData()
	}
	//gaoyan.PublishData()
	//gaoyan.PublishConfig()

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

func init() {

}

func main() {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range sigs {
			log.Printf("captured %v, stop job..", sig)
			done <- true
		}
	}()

loop:
	for {
		select {
		case <-time.After(time.Millisecond * 5000):
			log.Println("after some time, do job")
			do_job()

		case <-done:
			log.Println("job stopped, exting...")
			break loop
		}
	}

}
