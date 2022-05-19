package gaoyan

import (
	"fmt"
	"log"
	"math"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	modbus "github.com/goburrow/modbus"
	//modbusclient "github.com/dpapathanasiou/go-modbus"
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
	Timestamp int64     `json:"ts"` //timestamp
	V_Phase   []float64 `json:"v"`
	VA        float64   `json:"-"`
	VB        float64   `json:"-"`
	VC        float64   `json:"-"`
	I_Phase   []float64 `json:"i"`
	IA        float64   `json:"-"`
	IB        float64   `json:"-"`
	IC        float64   `json:"-"`
	P         float64   `json:"P"`
	P_Phase   []float64 `json:"p"`
	PA        float64   `json:"-"`
	PB        float64   `json:"-"`
	PC        float64   `json:"-"`
	Q         float64   `json:"Q"` // reactive power 无功功率 Q Var
	Q_Phase   []float64 `json:"q"`
	QA        float64   `json:"-"`
	QB        float64   `json:"-"`
	QC        float64   `json:"-"`
	S         float64   `json:"S"` // Aparent Power 视在功率 S VA
	S_Phase   []float64 `json:"s"`
	SA        float64   `json:"-"`
	SB        float64   `json:"-"`
	SC        float64   `json:"-"`
	F_Phase   []float64 `json:"f"`
	FA        float64   `json:"-"`
	FB        float64   `json:"-"`
	FC        float64   `json:"-"`
	E         float64   `json:"E"`
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

func (m METER) Read(host string, port int) (*METERInfo, error) {
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
		return nil, err
	} else {
		var length uint16 = 0x1E + 1
		client := modbus.NewClient(handler)
		results, err := client.ReadInputRegisters(0x00, length)
		// results, err = client.WriteMultipleRegisters(1, 2, []byte{0, 3, 0, 4})
		// results, err = client.WriteMultipleCoils(5, 10, []byte{4, 3})
		if err != nil {
			log.Println("read failed, ", err)
			return nil, err
		} else {
			info.VA = math.Round(byte16_to_float64(results, 0x00*2)) / 10
			info.VB = math.Round(byte16_to_float64(results, 0x01*2)) / 10
			info.VC = math.Round(byte16_to_float64(results, 0x02*2)) / 10
			info.V_Phase = append(info.V_Phase, info.VA)
			info.V_Phase = append(info.V_Phase, info.VB)
			info.V_Phase = append(info.V_Phase, info.VC)
			info.IA = byte16_to_float64(results, (0x03+0)*2) / 100
			info.IB = byte16_to_float64(results, (0x03+1)*2) / 100
			info.IC = byte16_to_float64(results, (0x03+2)*2) / 100
			info.I_Phase = append(info.I_Phase, info.IA)
			info.I_Phase = append(info.I_Phase, info.IB)
			info.I_Phase = append(info.I_Phase, info.IC)
			info.P = byte16_to_float64(results, (0x07)*2)
			info.PA = byte16_to_float64(results, (0x08+0)*2)
			info.PB = byte16_to_float64(results, (0x08+1)*2)
			info.PC = byte16_to_float64(results, (0x08+2)*2)
			info.P_Phase = append(info.P_Phase, info.PA)
			info.P_Phase = append(info.P_Phase, info.PB)
			info.P_Phase = append(info.P_Phase, info.PC)
			info.Q = byte16_to_float64(results, (0x0b)*2)
			info.QA = byte16_to_float64(results, (0x0c+0)*2)
			info.QB = byte16_to_float64(results, (0x0c+1)*2)
			info.QC = byte16_to_float64(results, (0x0c+2)*2)
			info.Q_Phase = append(info.Q_Phase, info.QA)
			info.Q_Phase = append(info.Q_Phase, info.QB)
			info.Q_Phase = append(info.Q_Phase, info.QC)
			info.S = byte16_to_float64(results, (0x0f)*2)
			info.SA = byte16_to_float64(results, (0x10+0)*2)
			info.SB = byte16_to_float64(results, (0x10+1)*2)
			info.SC = byte16_to_float64(results, (0x10+2)*2)
			info.S_Phase = append(info.S_Phase, info.SA)
			info.S_Phase = append(info.S_Phase, info.SB)
			info.S_Phase = append(info.S_Phase, info.SC)
			info.FA = byte16_to_float64(results, (0x1a+0)*2) / 100
			info.FB = byte16_to_float64(results, (0x1a+1)*2) / 100
			info.FC = byte16_to_float64(results, (0x1a+2)*2) / 100
			info.F_Phase = append(info.F_Phase, info.FA)
			info.F_Phase = append(info.F_Phase, info.FB)
			info.F_Phase = append(info.F_Phase, info.FC)
			info.E = byte32_to_float64(results, (0x1D)*2) / 100

			info.Timestamp = time.Now().Unix()
		}
	}

	return info, nil

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
	tpl := []string{"v[0]", "v[1]", "v[2]", "i[0]", "i[1]", "i[2]", "P", "p[0]", "p[1]", "p[2]", "Q", "q[0]", "q[1]", "q[2]", "S", "s[0]", "s[1]", "s[2]", "f[0]", "f[1]", "f[2]", "E"}
	unit := []string{"V", "V", "V", "A", "A", "A", "W", "W", "W", "W", "Var", "Var", "Var", "Var", "VA", "VA", "VA", "VA", "HZ", "HZ", "HZ", "kWh"}
	for i, c := range cat {
		config_topics = append(config_topics, fmt.Sprintf("haworkshopyc1/sensor/powermeter%s/%s/config", m.Name, c))
		config_payloads = append(config_payloads, fmt.Sprintf("{\"device_class\": \"power\", \"name\": \"power-meter-%s-%s\", \"unique_id\": \"power-meter-%s-%s\", \"state_topic\": \"%s\",   \"unit_of_measurement\": \"%s\" ,  \"value_template\": \"{{ value_json.%s }}\"  , \"expire_after\":120 }", m.Name, c, m.Name, c, topic_state, unit[i], tpl[i]))
	}

	for i, topic := range config_topics {
		//log.Println(config_payloads[i])
		c.Publish(topic, 2, true, config_payloads[i])
	}

}
