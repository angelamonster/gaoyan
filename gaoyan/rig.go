package gaoyan

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	claymore "github.com/ivanbeldad/rpc-claymore"
)

type Crypto struct {
	HashRate       int `json:"hashrate"`
	Shares         int `json:"shares"`
	RejectedShares int `json:"rejected"`
	InvalidShares  int `json:"invalid"`
}

type PoolInfo struct {
	Address  string `json:"adress"`
	Switches int    `json:"switches"`
}

// GPU Information about each concrete GPU
type GPU struct {
	HashRate    int `json:"hashrate"`
	AltHashRate int `json:"althashrate"`
	Temperature int `json:"temperature"`
	FanSpeed    int `json:"fanspeed"`
}

// MinerInfo Information about the miner
type MinerInfo struct {
	Version    string   `json:"version"`
	UpTime     int      `json:"uptime"`
	MainCrypto Crypto   `json:"maincrypto"`
	AltCrypto  Crypto   `json:"altcrypto"`
	MainPool   PoolInfo `json:"mainpool"`
	AltPool    PoolInfo `json:"altpool"`
	GPUS       []GPU
	Timestamp  int64 `json:"timestamp"`
	HighTemp   int   `json:"hightemperature"`
}

type RIG struct {
	ID           string
	IP           string
	Username     string
	Password     string
	ClaymorePort int
	ConfigSent   bool
}

func (rig RIG) GetStat() (string, error) {

	miner := claymore.Miner{Address: fmt.Sprintf("%s:%d", rig.ID, rig.ClaymorePort)}
	info, err := miner.GetInfo()

	if err != nil {
		log.Println(err)
		return "", err
	}

	var mi = new(MinerInfo)
	mi.HighTemp = 0
	for _, g := range info.GPUS {
		var gpu GPU = GPU{HashRate: g.HashRate, AltHashRate: g.AltHashRate, Temperature: g.Temperature, FanSpeed: g.FanSpeed}
		mi.GPUS = append(mi.GPUS, gpu)
		if g.Temperature > mi.HighTemp {
			mi.HighTemp = g.Temperature
		}
	}
	mi.MainCrypto = Crypto{HashRate: info.MainCrypto.HashRate, Shares: info.MainCrypto.Shares, RejectedShares: info.MainCrypto.RejectedShares, InvalidShares: info.MainCrypto.InvalidShares}
	mi.AltCrypto = Crypto{HashRate: info.AltCrypto.HashRate, Shares: info.AltCrypto.Shares, RejectedShares: info.AltCrypto.RejectedShares, InvalidShares: info.AltCrypto.InvalidShares}
	mi.MainPool = PoolInfo{Address: info.MainPool.Address, Switches: info.MainPool.Switches}
	mi.AltPool = PoolInfo{Address: info.AltPool.Address, Switches: info.AltPool.Switches}
	mi.Version = info.Version
	mi.UpTime = info.UpTime
	mi.Timestamp = time.Now().Unix()

	json_bytes, err := json.Marshal(mi)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(json_bytes), nil
	//return "okok test", nil
}

func (rig RIG) PublishData(c mqtt.Client, json_data string) {
	//log.Printf("PublishData for %s", rig.ID)

	topic_state := fmt.Sprintf("haworkshopyc1/sensor/%s/state", rig.ID)

	c.Publish(topic_state, 0, false, json_data)
}

func (rig RIG) PublishConfig(c mqtt.Client, json_data string) {
	log.Printf("PublishConfig for %s", rig.ID)
	topic_state := fmt.Sprintf("haworkshopyc1/sensor/%s/state", rig.ID)

	mi := new(claymore.MinerInfo)

	json.Unmarshal([]byte(json_data), &mi)

	//     #topic_totalpower_config = "haworkshopyc1/sensor/{}/totalpower/config".format(self.id)
	//     #totalpower_config = '{{"device_class": "power", "name": "{}-totalpower", "unique_id": "{}-totalpower", "state_topic": "{}",   "unit_of_measurement": "W","value_template": "{{{{ value_json.{}.totalpower }}}}" }}'.format(self.id,self.id,topic_state,self.id)
	//     #mqtt.client.publish(topic_totalpower_config, payload=totalpower_config, qos=2,retain=True)     # 发送消息

	//     topic_totalhash_config = "haworkshopyc1/sensor/{}/totalhash/config".format(self.id)
	//     totalhash_config = '{{"name": "{}-totalhash", "unique_id": "{}-totalhash", "state_topic": "{}",   "unit_of_measurement": "B","value_template": "{{{{ value_json.hash }}}}"  , "expire_after":"120" }}'.format(self.id,self.id,self.topic_state)
	//     mqtt.client.publish(topic_totalhash_config, payload=totalhash_config, qos=2,retain=True)     # 发送消息
	config_topics := []string{fmt.Sprintf("haworkshopyc1/sensor/%s/totalhash/config", rig.ID),
		fmt.Sprintf("haworkshopyc1/sensor/%s/hightemperature/config", rig.ID)}
	config_payloads := []string{fmt.Sprintf("{\"name\": \"%s-totalhash\", \"unique_id\": \"%s-totalhash\", \"state_topic\": \"%s\",   \"unit_of_measurement\": \"B\",\"value_template\": \"{{ value_json.maincrypto.hashrate }}\" }", rig.ID, rig.ID, topic_state),
		fmt.Sprintf("{\"device_class\": \"temperature\", \"name\": \"%s-hightemperature\", \"unique_id\": \"%s-hightemperature\", \"state_topic\": \"%s\",   \"unit_of_measurement\": \"°C\",\"value_template\": \"{{ value_json.hightemperature }}\" }", rig.ID, rig.ID, topic_state)}

	for i, _ := range mi.GPUS {
		config_topics = append(config_topics, fmt.Sprintf("haworkshopyc1/sensor/%s-%d/temp/config", rig.ID, i))
		config_payloads = append(config_payloads, fmt.Sprintf("{\"device_class\": \"temperature\", \"name\": \"%s-%d-temp\", \"unique_id\": \"%s-%d-temp\", \"state_topic\": \"%s\",   \"unit_of_measurement\": \"°C\" ,  \"value_template\": \"{{ value_json.GPUS[%d].temperature }}\"  , \"expire_after\":120 }", rig.ID, i, rig.ID, i, topic_state, i))

		config_topics = append(config_topics, fmt.Sprintf("haworkshopyc1/sensor/%s-%d/hash/config", rig.ID, i))
		config_payloads = append(config_payloads, fmt.Sprintf("{\"name\": \"%s-%d-hash\", \"unique_id\": \"%s-%d-hash\", \"state_topic\": \"%s\",   \"unit_of_measurement\": \"B\" ,  \"value_template\": \"{{ value_json.GPUS[%d].hashrate }}\"  , \"expire_after\":120 }", rig.ID, i, rig.ID, i, topic_state, i))

		config_topics = append(config_topics, fmt.Sprintf("haworkshopyc1/sensor/%s-%d/fan/config", rig.ID, i))
		config_payloads = append(config_payloads, fmt.Sprintf("{\"name\":  \"%s-%d-fan\", \"unique_id\": \"%s-%d-fan\",  \"state_topic\": \"%s\",   \"unit_of_measurement\": \"%%\" ,  \"value_template\": \"{{ value_json.GPUS[%d].fanspeed }}\"  , \"expire_after\":120 }", rig.ID, i, rig.ID, i, topic_state, i))

	}

	for i, topic := range config_topics {
		//log.Println(config_payloads[i])
		c.Publish(topic, 2, true, config_payloads[i])
	}

}
