package gaoyan

import (
	"encoding/json"
	"fmt"
	"log"

	claymore "github.com/ivanbeldad/rpc-claymore"
)

type RIG struct {
	ID           string
	IP           string
	Username     string
	Password     string
	ClaymorePort int
}

func GetStat(rig RIG) {
	log.Println("GetStat")

	miner := claymore.Miner{Address: fmt.Sprintf("%s:%d", rig.ID, rig.ClaymorePort)}
	info, err := miner.GetInfo()

	if err != nil {
		log.Fatal(err)
	}

	json_bytes, err := json.Marshal(info)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(json_bytes))
}

func (rig RIG) PublishConfig() {
	log.Printf("PublishConfig for %s", rig.ID)

}
func (rig RIG) PublishData() {
	log.Printf("PublishData for %s", rig.ID)

}
