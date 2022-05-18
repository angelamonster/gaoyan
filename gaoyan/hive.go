package gaoyan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type EarningStat struct {
	MeanReward float64 `json:"meanReward"`
	Reward     float64 `json:"reward"`
	Timestamp  string  `json:"timestamp"`
}
type PayOut struct {
	Amount          float64 `json:"amount"`
	ApproveUUID     string  `json:"approveUUID"`
	Coin            string  `json:"coin"`
	CreatedAt       string  `json:"createdAt"`
	Meta            string  `json:"meta"`
	PayoutDirection string  `json:"payoutDirection"`
	Status          string  `json:"status"`
	TxHash          string  `json:"txHash"`
	TxMeta          string  `json:"txMeta"`
	Type            string  `json:"type"`
	UpdatedAt       string  `json:"updatedAt"`
	UserAddress     string  `json:"userAddress"`
	Uuid            string  `json:"uuid"`
}

type BILL struct {
	EarningStats       []EarningStat `json:"earningStats"`
	ExpectedReward24H  float64       `json:"expectedReward24H"`
	ExpectedRewardWeek float64       `json:"ExpectedRewardWeek"`
	PendingPayouts     []PayOut      `json:"pendingPayouts"`
	SucceedPayouts     []PayOut      `json:"succeedPayouts"`
	TotalPaid          float64       `json:"totalPaid"`
	TotalUnpaid        float64       `json:"totalUnpaid"`
}

type SharesStatusStats struct {
	InvalidCount string  `json:"invalidCount"`
	InvalidRate  float64 `json:"invalidRate"`
	LastShareDt  string  `json:"lastShareDt"`
	StaleCount   string  `json:"staleCount"`
	StaleRate    float64 `json:"staleRate"`
	ValidCount   string  `json:"validCount"`
	ValidRate    float64 `json:"validRate"`
}

type Stat struct {
	Hashrate            string            `json:"hashrate"`
	Hashrate24h         string            `json:"hashrate24h"`
	OnlineWorkerCount   string            `json:"onlineWorkerCount"`
	ReportedHashrate    string            `json:"reportedHashrate"`
	ReportedHashrate24h string            `json:"reportedHashrate24h"`
	SharesStatusStats   SharesStatusStats `json:"SharesStatusStats"`
}

type HIVE struct {
	Address    string `json:"address"`
	ConfigSent bool
}

type HIVEInfo struct {
	OnlineWorkerCount string  `json:"w"`
	Hashrate          float64 `json:"hr"`
	ReportedHashrate  float64 `json:"rhr"`
	TotalPaid         float64 `json:"tp"`
	TotalUnpaid       float64 `json:"tup"`
	ExpectedReward24H float64 `json:"e24h"`
}

func (hive HIVE) Read() (*HIVEInfo, error) {
	var info = new(HIVEInfo)

	url_bill := fmt.Sprintf("https://hiveon.net/api/v1/stats/miner/%s/ETH/billing-acc", hive.Address)
	resp, err := http.Get(url_bill)
	if err != nil {
		log.Println(err)
		return nil, err
	} else {

		html, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(html))
		var bill = new(BILL)
		err := json.Unmarshal(html, bill)
		if err != nil {
			log.Println(err)
			return nil, err
		} else {
			//log.Println(bill)
			log.Printf("unpaid:%f,forcast 24h:%f\n", bill.TotalUnpaid, bill.ExpectedReward24H)
			info.ExpectedReward24H = bill.ExpectedReward24H
			info.TotalUnpaid = bill.TotalUnpaid
			info.TotalPaid = bill.TotalPaid
		}
	}

	url_stat := fmt.Sprintf("https://hiveon.net/api/v1/stats/miner/%s/ETH", hive.Address)
	resp, err = http.Get(url_stat)
	if err != nil {
		log.Println(err)
		return nil, err
	} else {
		html, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(html))
		var sta = new(Stat)
		err := json.Unmarshal(html, sta)
		if err != nil {
			log.Println(err)

			return nil, err
		} else {
			log.Printf("online:%s\n", sta.OnlineWorkerCount)
			info.OnlineWorkerCount = sta.OnlineWorkerCount
			h, err := strconv.Atoi(sta.Hashrate)
			rh, rerr := strconv.Atoi(sta.Hashrate)
			if err != nil && rerr != nil {
				info.Hashrate = float64(h / 1000000)
				info.ReportedHashrate = float64(rh / 1000000)
			}
		}
	}

	return info, nil
}

func (hive HIVE) PublishData(c mqtt.Client, json_data string) {
	//log.Printf("PublishData for %s", rig.ID)
	topic_state := fmt.Sprintf("haworkshopyc1/sensor/hive/state")

	c.Publish(topic_state, 0, false, json_data)
}

func (hive HIVE) PublishConfig(c mqtt.Client) {
	log.Printf("Publish Config for hive pool")
	topic_state := fmt.Sprintf("haworkshopyc1/sensor/hive/state")

	config_topics := []string{}
	config_payloads := []string{}

	cat := []string{"onlineWorkerCount", "hashrate", "reportedHashrate", "totalPaid", "totalUnpaid", "expectedReward24H"}
	scat := []string{"w", "hr", "rhr", "tp", "tup", "e24h"}
	unit := []string{"U", "MH", "MH", "ETH", "ETH", "ETH"}
	for i, c := range cat {
		config_topics = append(config_topics, fmt.Sprintf("haworkshopyc1/sensor/hive/%s/config", c))
		config_payloads = append(config_payloads, fmt.Sprintf("{\"name\": \"hive-%s\", \"unique_id\": \"hive-%s\", \"state_topic\": \"%s\",   \"unit_of_measurement\": \"%s\" ,  \"value_template\": \"{{ value_json.%s }}\"  , \"expire_after\":600 }", c, c, topic_state, unit[i], scat[i]))
	}

	for i, topic := range config_topics {
		//log.Println(config_payloads[i])
		c.Publish(topic, 2, true, config_payloads[i])
	}
}
