package rig

import (
	"encoding/json"
	"fmt"
	"log"

	claymore "github.com/ivanbeldad/rpc-claymore"
)

func Update() {
	log.Println("update")

	miner := claymore.Miner{Address: "w0004:3334"}
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
