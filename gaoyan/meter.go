package gaoyan

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	modbus "github.com/goburrow/modbus"

	//modbusclient "github.com/dpapathanasiou/go-modbus"

	"encoding/json"
	//modbus "github.com/thinkgos/gomodbus"
)

//"encoding/json"
// "fmt"
// "log"

// claymore "../rpcclaymore"
// mqtt "github.com/eclipse/paho.mqtt.golang"

type METER struct {
	Name       string `json:"name"`
	ConfigSent bool
}
type METERInfo struct {
	Timestamp int64   `json:"timestamp"`
	VA        float64 `json:"va"`
	VB        float64 `json:"vb"`
	VC        float64 `json:"vc"`
	IA        float64 `json:"ia"`
	IB        float64 `json:"ib"`
	IC        float64 `json:"ic"`
	P         float64 `json:"p"`
	PA        float64 `json:"pa"`
	PB        float64 `json:"pb"`
	PC        float64 `json:"pc"`
	Q         float64 `json:"q"` // reactive power 无功功率 Q Var
	QA        float64 `json:"qa"`
	QB        float64 `json:"qb"`
	QC        float64 `json:"qc"`
	S         float64 `json:"s"` // Aparent Power 视在功率 S VA
	SA        float64 `json:"sa"`
	SB        float64 `json:"sb"`
	SC        float64 `json:"sc"`
	FA        float64 `json:"fa"`
	FB        float64 `json:"fb"`
	FC        float64 `json:"fc"`
	E         float64 `json:"e"`
}

func byte16_to_float64(results []byte, pos int) float64 {

	var v uint16 = uint16(results[pos])<<8 + uint16(results[pos+1])
	var s int32 = 0x10000

	if v < 0x8000 {
		return float64(v)
	} else {
		return float64(int32(v) - s)
	}
}

func byte32_to_float64(results []byte, pos int) float64 {

	var v uint32 = uint32(results[pos])<<24 + uint32(results[pos+1])<<16 + uint32(results[pos+2])<<8 + uint32(results[pos+3])
	var s int64 = 0x100000000

	if v < 0x80000000 {
		return float64(v)
	} else {
		return float64(int64(v) - s)
	}

}

func (m METER) Read(host string, port int) (string, error) {
	info := new(METERInfo)
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", host, port))
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 0x01
	//handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	// Connect manually so that multiple requests are handled in one connection session
	err := handler.Connect()
	defer handler.Close()

	if err != nil {
		log.Println("connect failed, ", err)
		return "", err
	} else {
		var length uint16 = 0x1E + 1
		client := modbus.NewClient(handler)
		results, err := client.ReadInputRegisters(0x00, length)
		// results, err = client.WriteMultipleRegisters(1, 2, []byte{0, 3, 0, 4})
		// results, err = client.WriteMultipleCoils(5, 10, []byte{4, 3})
		if err != nil {
			log.Println("read failed, ", err)
			return "", err
		} else {
			info.VA = byte16_to_float64(results, 0x00*2) * 0.1
			info.VB = byte16_to_float64(results, 0x01*2) * 0.1
			info.VC = byte16_to_float64(results, 0x02*2) * 0.1

			info.IA = byte16_to_float64(results, (0x03+0)*2) * 0.01
			info.IB = byte16_to_float64(results, (0x03+1)*2) * 0.01
			info.IC = byte16_to_float64(results, (0x03+2)*2) * 0.01
			info.P = byte16_to_float64(results, (0x07)*2)
			info.PA = byte16_to_float64(results, (0x08+0)*2)
			info.PB = byte16_to_float64(results, (0x08+1)*2)
			info.PC = byte16_to_float64(results, (0x08+2)*2)
			info.Q = byte16_to_float64(results, (0x0b)*2)
			info.QA = byte16_to_float64(results, (0x0c+0)*2)
			info.QB = byte16_to_float64(results, (0x0c+1)*2)
			info.QC = byte16_to_float64(results, (0x0c+2)*2)
			info.S = byte16_to_float64(results, (0x0f)*2)
			info.SA = byte16_to_float64(results, (0x10+0)*2)
			info.SB = byte16_to_float64(results, (0x10+1)*2)
			info.SC = byte16_to_float64(results, (0x10+2)*2)
			info.FA = byte16_to_float64(results, (0x1a+0)*2) * 0.01
			info.FB = byte16_to_float64(results, (0x1a+1)*2) * 0.01
			info.FC = byte16_to_float64(results, (0x1a+2)*2) * 0.01

			info.E = byte32_to_float64(results, (0x1D)*2) * 0.01

			info.Timestamp = time.Now().Unix()
		}
	}

	json_bytes, err := json.Marshal(info)

	return string(json_bytes), nil

}

func (m METER) PublishData(c mqtt.Client, json_data string) {
	//log.Printf("PublishData for %s", rig.ID)

	topic_state := fmt.Sprintf("haworkshopyc1/sensor/powermeter%s/state", m.Name)

	c.Publish(topic_state, 0, false, json_data)
}

func (m METER) PublishConfig(c mqtt.Client) {
	log.Printf("Publish Config for %s", m.Name)
	topic_state := fmt.Sprintf("haworkshopyc1/sensor/powermeter%s/state", m.Name)

	config_topics := []string{}
	config_payloads := []string{}

	cat := []string{"va", "vb", "vc", "ia", "ib", "ic", "p", "pa", "pb", "pc", "q", "qa", "qb", "qc", "s", "sa", "sb", "sc", "fa", "fb", "fc", "e"}
	unit := []string{"V", "V", "V", "A", "A", "A", "W", "W", "W", "W", "Var", "Var", "Var", "Var", "VA", "VA", "VA", "VA", "HZ", "HZ", "HZ", "kWh"}
	for i, c := range cat {
		config_topics = append(config_topics, fmt.Sprintf("haworkshopyc1/sensor/powermeter%s/%s/config", m.Name, c))
		config_payloads = append(config_payloads, fmt.Sprintf("{\"device_class\": \"power\", \"name\": \"power-meter-%s-%s\", \"unique_id\": \"power-meter-%s-%s\", \"state_topic\": \"%s\",   \"unit_of_measurement\": \"%s\" ,  \"value_template\": \"{{ value_json.%s }}\"  , \"expire_after\":120 }}", m.Name, c, m.Name, c, topic_state, unit[i], c))
	}

	for i, topic := range config_topics {
		//log.Println(config_payloads[i])
		c.Publish(topic, 2, true, config_payloads[i])
	}

}
