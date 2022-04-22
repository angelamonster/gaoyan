package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	claymore "./rpcclaymore"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//{"P":"2588","PA":"755","PB":"836","PC":"1001","AP":"2694","APA":"778","APB":"865","APC":"1050","VA":"225.1","VB":"223.1","VC":"224.1","FA":"49.99","FB":"49.99","FC":"49.99","E":"27759.83"}
type METER struct {
	P         int     `json:"P"`
	PA        int     `json:"PA"`
	PB        int     `json:"PB"`
	PC        int     `json:"PC"`
	AP        int     `json:"AP"`
	APA       int     `json:"APA"`
	APB       int     `json:"APB"`
	APC       int     `json:"APC"`
	VA        float64 `json:"VA"`
	VB        float64 `json:"VB"`
	VC        float64 `json:"VC"`
	FA        float64 `json:"FA"`
	FB        float64 `json:"FB"`
	FC        float64 `json:"FC"`
	E         float64 `json:"E"`
	Timestamp int64   `json:"timestamp"`
}

//{"P":3188,"PSUNGROW":2228.0,"PGINLONG":960.0,"E":26048.0 ,"ESUNGROW":24139.0 ,"EGINLONG":1909.0 ,"VA":234.4 ,"VB":234.4 ,"VC":235.0,"F":50.03,"timestamp":1649911033}
type SOLAR struct {
	P         float64 `json:"P"`
	PSUNGROW  float64 `json:"PSUNGROW"`
	PGINLONG  float64 `json:"PGINLONG"`
	E         float64 `json:"E"`
	ESUNGROW  float64 `json:"ESUNGROW"`
	EGINLONG  float64 `json:"EGINLONG"`
	VA        float64 `json:"VA,omitempty"`
	VB        float64 `json:"VB,omitempty"`
	VC        float64 `json:"VC,omitempty"`
	F         float64 `json:"F,omitempty"`
	Timestamp int64   `json:"timestamp"`
}

var w0004 claymore.MinerInfo
var w0005 claymore.MinerInfo
var w0007 claymore.MinerInfo

var meter METER

var solar SOLAR

//var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
func onMessageReceived(client mqtt.Client, msg mqtt.Message) {
	//fmt.Printf("TOPIC: %s\n", msg.Topic())
	//fmt.Printf("MSG: %s\n", msg.Payload())
	//var jsonStr = []byte(`{"uptime": "11:29", "hash": 305956000, "rig": "w0005", "timestamp": 1649845484}`)
	//var rig RIG
	//if err := json.Unmarshal(jsonStr, &i_msg); err != nil {
	//haworkshopyc1/sensor/w0004/state
	if w000_pos := strings.Index(msg.Topic(), "w000"); w000_pos != -1 {
		rig_id := msg.Topic()[w000_pos+4 : w000_pos+5]
		//fmt.Printf("%s", rig_id)
		switch rig_id {
		case "4":
			if err := json.Unmarshal(msg.Payload(), &w0004); err != nil {
				fmt.Println(err)
			}
			queue = dequeue(queue)
			queue = enqueue(queue, 4)
		case "5":
			if err := json.Unmarshal(msg.Payload(), &w0005); err != nil {
				fmt.Println(err)
			}
			queue = dequeue(queue)
			queue = enqueue(queue, 5)
		case "7":
			if err := json.Unmarshal(msg.Payload(), &w0007); err != nil {
				fmt.Println(err)
			}
			queue = dequeue(queue)
			queue = enqueue(queue, 7)
		default:
			fmt.Printf("TOPIC: %s\n", msg.Topic())
			fmt.Printf("MSG: %s\n", msg.Payload())
		}
	}

	// haworkshopyc1/sensor/powermeteryc1/state
	if pos := strings.Index(msg.Topic(), "powermeteryc1"); pos != -1 {
		if err := json.Unmarshal(msg.Payload(), &meter); err != nil {
			fmt.Println(err)
		}
		queue = dequeue(queue)
		queue = enqueue(queue, 1)
	}

	// "haworkshopyc1/sensor/solaryc1/state"
	if pos := strings.Index(msg.Topic(), "solaryc1"); pos != -1 {
		if err := json.Unmarshal(msg.Payload(), &solar); err != nil {
			fmt.Println(err)
		}
		queue = dequeue(queue)
		queue = enqueue(queue, 2)
	}

	//if err := json.Unmarshal(msg.Payload(), &rig); err != nil {
	// if err := json.Unmarshal(msg.Payload(), &(rigs[rig_id])); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	//for _, movie := range moviesBack {
	// 	//	fmt.Println(movie.Title)
	// 	//}
	// 	//secNow := time.Now().Unix()
	// 	//fmt.Printf("%s : %dM - %d\n", rig.ID, rig.Hash/1000000, secNow-int64(rig.Timestamp))
	// 	//rigs[rig.ID] = rig
	// }
}

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

var sig_exit = false
var mqtt_status = 0 // 0 offline, 1 online, 2 connecting

func mqtt_reconnect(c mqtt.Client) {
	// 获取本机的MAC地址

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//fmt.Println("subscribe -- ")
	// if token := c.Subscribe("haworkshopyc1/sensor/w0004/state", 0, onMessageReceived); token.Wait() && token.Error() != nil {
	// 	panic(token.Error())
	// }
	// if token := c.Subscribe("haworkshopyc1/sensor/w0005/state", 0, onMessageReceived); token.Wait() && token.Error() != nil {
	// 	panic(token.Error())
	// }
	// if token := c.Subscribe("haworkshopyc1/sensor/w0007/state", 0, onMessageReceived); token.Wait() && token.Error() != nil {
	// 	panic(token.Error())
	// }

	// if token := c.Subscribe("haworkshopyc1/sensor/powermeteryc1/state", 0, onMessageReceived); token.Wait() && token.Error() != nil {
	// 	panic(token.Error())
	// }

	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			mqtt_cleanup(c)
		}
	}()

}

func mqtt_cleanup(c mqtt.Client) {

	fmt.Println("unsubscribe -- ")

	// if token := c.Unsubscribe("haworkshopyc1/sensor/w0004/state"); token.Wait() && token.Error() != nil {
	// 	panic(token.Error())
	// }
	// if token := c.Unsubscribe("haworkshopyc1/sensor/w0005/state"); token.Wait() && token.Error() != nil {
	// 	panic(token.Error())
	// }
	// if token := c.Unsubscribe("haworkshopyc1/sensor/w0007/state"); token.Wait() && token.Error() != nil {
	// 	panic(token.Error())
	// }

	// if token := c.Unsubscribe("haworkshopyc1/sensor/powermeteryc1/state"); token.Wait() && token.Error() != nil {
	// 	panic(token.Error())
	// }

	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
		c.Disconnect(250)
	}()
}

func enqueue(queue []int, element int) []int {
	queue = append(queue, element) // Simply append to enqueue.
	//fmt.Println("Enqueued:", element)
	return queue
}

func dequeue(queue []int) []int {
	//element := queue[0] // The first element is the one to be dequeued.
	//fmt.Println("Dequeued:", element)
	return queue[1:] // Slice off the element once it is dequeued.
}

var queue []int

var enery_start_today int
var last_day = 0

func paint_console() {
	fmt.Print("\033[u\033[K") // restore the cursor position and clear the line

	fmt.Print("\033[s") // save the cursor position
	fmt.Println("")

	var hash_4, hash_5, hash_7, pa, pb, pc, psolar, psolar_sungrow, psolar_ginlong, esolar string
	var high_temp_4, high_temp_5, high_temp_7 string
	var delay_4, delay_5, delay_7, delay_meter, delay_solar int64

	switch mqtt_status {
	case 0:
		fmt.Println("-- offline --")
	case 1:
		fmt.Println("-- online ---")
	case 2:
		fmt.Println("-- offline --")
	default:
		fmt.Println("-- unkown ---")
	}

	//for _, rig := range rigs {
	secNow := time.Now().Unix()
	// fmt.Printf("%s : %dM - %ds | %dW \n", w0004.ID, w0004.Hash/1000000, secNow-int64(w0004.Timestamp), meter.PA)
	// fmt.Printf("%s : %dM - %ds | %dW \n", w0005.ID, w0005.Hash/1000000, secNow-int64(w0005.Timestamp), meter.PB)
	// fmt.Printf("%s : %dM - %ds | %dW \n", w0007.ID, w0007.Hash/1000000, secNow-int64(w0007.Timestamp), meter.PC)

	if delay_4 = secNow - int64(w0004.Timestamp); delay_4 > 60 {
		//fmt.Printf("%s : offline \n", "w0004")
		hash_4 = "-"
		high_temp_4 = "-"
		//delay_4 = "-"
	} else {
		//fmt.Printf("%s : %ds ago \n", "w0004", delay)
		hash_4 = strconv.Itoa(w0004.MainCrypto.HashRate / 1000000)
		high_temp_4 = strconv.Itoa(w0004.HighTemp)
		//delay_4 = strconv.Itoa(int(delay))
	}
	if delay_5 = secNow - int64(w0005.Timestamp); delay_5 > 60 {
		//fmt.Printf("%s : offline \n", "w0005")
		hash_5 = "-"
		high_temp_5 = "-"
		//delay_5 = "-"
	} else {
		//fmt.Printf("%s : %ds ago \n", "w0005", delay)
		hash_5 = strconv.Itoa(w0005.MainCrypto.HashRate / 1000000)
		high_temp_5 = strconv.Itoa(w0005.HighTemp)
		//delay_5 = strconv.Itoa(int(delay))
	}
	if delay_7 = secNow - int64(w0007.Timestamp); delay_7 > 60 {
		//fmt.Printf("%s : offline \n", "w0007")
		hash_7 = "-"
		high_temp_7 = "-"
		//delay_7 = "-"
	} else {
		//fmt.Printf("%s : %ds ago \n", "w0007", delay)
		hash_7 = strconv.Itoa(w0007.MainCrypto.HashRate / 1000000)
		high_temp_7 = strconv.Itoa(w0007.HighTemp)
		//delay_7 = strconv.Itoa(int(delay))
	}
	//fmt.Printf("%s : %ds ago \n", "w0005", secNow-int64(w0005.Timestamp))
	//fmt.Printf("%s : %ds ago \n", "w0007", secNow-int64(w0007.Timestamp))
	//fmt.Printf("%s : %ds ago \n", "power", secNow-int64(meter.Timestamp))
	if delay_meter = secNow - int64(meter.Timestamp); delay_meter > 60 {
		//fmt.Printf("%s : offline \n", "meter")
		pa = "-"
		pb = "-"
		pc = "-"

	} else {
		//fmt.Printf("%s : %ds ago \n", "meter", delay)
		pa = strconv.Itoa(meter.PA)
		pb = strconv.Itoa(meter.PB) // + " (" + strconv.Itoa(int(0-delay)) + "s)"
		pc = strconv.Itoa(meter.PC) // + " (" + strconv.Itoa(int(0-delay)) + "s)"
		//delay_meter =  strconv.Itoa(int(delay))
	}

	if delay_solar = secNow - int64(solar.Timestamp); delay_solar > 300 {
		//fmt.Printf("%s : offline \n", "meter")
		psolar = "-"
		psolar_sungrow = "-"
		psolar_ginlong = "-"
	} else {
		//fmt.Printf("%s : %ds ago \n", "meter", delay)
		psolar = strconv.Itoa(int(solar.P))                //+ " (" + strconv.Itoa(int(0-delay)) + "s)"
		psolar_sungrow = strconv.Itoa(int(solar.PSUNGROW)) //+ " (" + strconv.Itoa(int(0-delay)) + "s)"
		psolar_ginlong = strconv.Itoa(int(solar.PGINLONG)) //+ " (" + strconv.Itoa(int(0-delay)) + "s)"

		if last_day != time.Now().Day() {
			enery_start_today = int(solar.E)
		}
		last_day = time.Now().Day()
		esolar = strconv.Itoa(int(solar.E - float64(enery_start_today)))
		//esolar_ginlong = strconv.Itoa(int(solar.EGINLONG))
		//esolar_sungrow = strconv.Itoa(int(solar.ESUNGROW))
	}

	//==========================================================================
	t := table.NewWriter()
	//==========================================================================

	//==========================================================================
	// Append a few rows and render to console
	//==========================================================================
	// a row need not be just strings
	// t.AppendHeader(table.Row{"#", "ID", "Hash   (MB)", "Power   (W)", "Solar      (W)"})
	// t.AppendRow(table.Row{"1", "w0004", hash_4, pa, psolar_ginlong})
	// t.AppendRow(table.Row{"2", "w0005", hash_5, pb, psolar_sungrow})
	// t.AppendRow(table.Row{"3", "w0007", hash_7, pc, ""})
	// t.AppendSeparator()
	// t.AppendFooter(table.Row{"", "Total", (w0004.Hash + w0005.Hash + w0007.Hash) / 1000000, meter.PA + meter.PB + meter.PC, psolar})

	t.AppendHeader(table.Row{"ID", "Hash(MB)", "TEMP", "Power(W)"})
	t.AppendRow(table.Row{"4", hash_4, high_temp_4, pa})
	t.AppendRow(table.Row{"5", hash_5, high_temp_5, pb})
	t.AppendRow(table.Row{"7", hash_7, high_temp_7, pc})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Total", (w0004.MainCrypto.HashRate + w0005.MainCrypto.HashRate + w0007.MainCrypto.HashRate) / 1000000, "", meter.PA + meter.PB + meter.PC})
	t.AppendSeparator()
	t.AppendRow(table.Row{"GLONG", "", "", psolar_ginlong})
	t.AppendRow(table.Row{"SUN", "", "", psolar_sungrow})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Total", "", esolar, psolar})

	//==========================================================================
	// ASCII is too simple for me.
	//==========================================================================

	t.SetStyle(table.StyleLight)
	t.SetStyle(table.StyleRounded)

	//t.SetCaption("Table using the style 'StyleLight'.\n")

	//colorBOnW := text.Colors{text.BgWhite, text.FgBlack}
	colorHeader := text.Colors{text.BgBlack, text.FgWhite}
	// set colors using Colors/ColorsHeader/ColorsFooter
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Colors: text.Colors{text.FgWhite, text.BgBlack}, ColorsHeader: colorHeader},
		{Number: 2, Align: text.AlignRight, Colors: text.Colors{text.FgGreen, text.BgBlack}, ColorsHeader: colorHeader},
		{Number: 3, Align: text.AlignRight, Colors: text.Colors{text.FgGreen, text.BgBlack}, ColorsHeader: colorHeader, ColorsFooter: colorHeader},
		{Number: 4, Align: text.AlignRight, Colors: text.Colors{text.FgGreen, text.BgBlack}, ColorsHeader: colorHeader, ColorsFooter: colorHeader},
		{Number: 5, Align: text.AlignRight, Colors: text.Colors{text.FgGreen, text.BgBlack}, ColorsHeader: colorHeader, ColorsFooter: colorHeader},
	})
	//t.SetAllowedRowLength(50)

	//CallClear()

	fmt.Println(t.Render())

	queue = dequeue(queue)
	queue = enqueue(queue, 0)

	log := ""
	for _, v := range queue {
		if v == 0 {
			log += "_"
		} else if v == 1 {
			log += "P"
		} else if v == 2 {
			log += "S"
		} else {
			log += strconv.Itoa(v)
		}
	}
	fmt.Println(log)
	//fmt.Println("GOOD")
}

func init_mqtt() mqtt.Client {
	// for i := 0; i < 5; i++ {
	// 	text := fmt.Sprintf("this is msg #%d!", i)
	// 	token := c.Publish("go-mqtt/sample", 0, false, text)
	// 	token.Wait()
	// }

	//fmt.Printf("wait -- ")
	hostname, _ := os.Hostname()

	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	// opts := mqtt.NewClientOptions().AddBroker("tcp://iot.wiin.win:11883").SetClientID("Status-" + hostname)
	// opts.SetKeepAlive(2 * time.Second)
	// //opts.SetDefaultPublishHandler(f)
	// opts.SetPingTimeout(1 * time.Second)
	// opts.SetUsername("mao")
	// opts.SetPassword("linmao8888")

	server := flag.String("server", "tls://w.wiin.win:1884", "The full url of the MQTT server to connect to ex: tcp://iot.wiin.win:11883")
	//topic := flag.String("topic", "#", "Topic to subscribe to")
	//qos := flag.Int("qos", 0, "The QoS to subscribe to messages at")
	//clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	clientid := flag.String("clientid", hostname, "A clientid for the connection")
	username := flag.String("username", "mao", "A username to authenticate to the MQTT server")
	password := flag.String("password", "linmao8888", "Password to match username")
	flag.Parse()

	connOpts := mqtt.NewClientOptions().AddBroker(*server).SetClientID(*clientid).SetCleanSession(true)
	if *username != "" {
		connOpts.SetUsername(*username)
		if *password != "" {
			connOpts.SetPassword(*password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c mqtt.Client) {
		topics := map[string]byte{"haworkshopyc1/sensor/w0004/state": 0,
			"haworkshopyc1/sensor/w0005/state":         0,
			"haworkshopyc1/sensor/w0007/state":         0,
			"haworkshopyc1/sensor/powermeteryc1/state": 0,
			"haworkshopyc1/sensor/solaryc1/state":      0}
		if token := c.SubscribeMultiple(topics, onMessageReceived); token.Wait() && token.Error() != nil {
			//if token := c.Subscribe(*topic, byte(*qos), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		mqtt_status = 1

		CallClear()
	}
	connOpts.OnConnectionLost = func(c mqtt.Client, err error) {
		mqtt_status = 0
	}
	connOpts.OnReconnecting = func(c mqtt.Client, co *mqtt.ClientOptions) {
		mqtt_status = 2
	}

	return mqtt.NewClient(connOpts)
}

func main() {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	go func() {
		for sig := range ch {
			// sig is a ^C,handle it
			fmt.Println("got signal:", sig)
			sig_exit = true
		}
	}()

	CallClear()

	mqtt_client := init_mqtt()

	for i := 0; i < 150; i++ {
		queue = enqueue(queue, 0)
	}
	//time.Sleep(30 * time.Second)
	for !sig_exit {
		//fmt.Printf("==================================\n")
		if !mqtt_client.IsConnected() {
			mqtt_reconnect(mqtt_client)
		}

		paint_console()

		//}
		time.Sleep(2000 * time.Millisecond)
	}

	mqtt_cleanup(mqtt_client)
}
