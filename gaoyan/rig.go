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
	ConfigSent   bool
}

func (rig RIG) GetStat() (string, error) {

	miner := claymore.Miner{Address: fmt.Sprintf("%s:%d", rig.ID, rig.ClaymorePort)}
	info, err := miner.GetInfo()

	if err != nil {
		log.Println(err)
		return "", err
	}

	json_bytes, err := json.Marshal(info)
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
	config_topics := []string{fmt.Sprintf("haworkshopyc1/sensor/%s/totalhash/config", rig.ID)}
	config_payloads := []string{fmt.Sprintf("{\"name\": \"%s-totalhash\", \"unique_id\": \"%s-totalhash\", \"state_topic\": \"%s\",   \"unit_of_measurement\": \"B\",\"value_template\": \"{{ value_json.maincrypto.hashrate }}\" }", rig.ID, rig.ID, topic_state)}

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
