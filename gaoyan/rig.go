package gaoyan

import (
	"encoding/json"
	"fmt"
	"log"

	claymore "../rpcclaymore"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type RIG struct {
	ID           string
	IP           string
	Username     string
	Password     string
	ClaymorePort int
}

func (rig RIG) GetStat() (string, error) {
	log.Println("GetStat")

	miner := claymore.Miner{Address: fmt.Sprintf("%s:%d", rig.ID, rig.ClaymorePort)}
	info, err := miner.GetInfo()

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	json_bytes, err := json.Marshal(info)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(json_bytes), nil
	//return "okok test", nil
}

func (rig RIG) PublishData(c mqtt.Client, json_data string) {
	log.Printf("PublishData for %s", rig.ID)

	topic_state := fmt.Sprintf("haworkshopyc1/sensor/%s/state", rig.ID)

	c.Publish(topic_state, 0, false, json_data)
}

func (rig RIG) PublishConfig(c mqtt.Client, json_data string) {
	log.Printf("PublishConfig for %s", rig.ID)
	topic_state := fmt.Sprintf("haworkshopyc1/sensor/%s/state", rig.ID)

	 mi := new(claymore.MinerInfo)
	
	json.Unmarshal([]byte(json_data),&mi)
	
	for i,g in range mi.GPUS{
		log.println(g)		
	}
}
