package main

import (
	"crypto/tls"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	gaoyan "./gaoyan"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var rig_w0004 = gaoyan.RIG{ID: "w0004", IP: "192.168.0.204", Username: "user", Password: "1", ClaymorePort: 3334, ConfigSent: false}
var rig_w0005 = gaoyan.RIG{ID: "w0005", IP: "192.168.0.205", Username: "user", Password: "9UNXmhyV", ClaymorePort: 3334, ConfigSent: false}
var rig_w0007 = gaoyan.RIG{ID: "w0007", IP: "192.168.0.207", Username: "user", Password: "1", ClaymorePort: 3334, ConfigSent: false}

var rigs = [3]gaoyan.RIG{rig_w0004, rig_w0005, rig_w0007}

func do_job(c mqtt.Client) {
	log.Println("loop")

	for i, _ := range rigs {
		go func(i int) {
			json_string, err := rigs[i].GetStat()
			if err != nil {
				log.Println(err)
			} else {
				if false == rigs[i].ConfigSent {
					rigs[i].PublishConfig(c, json_string)
					rigs[i].ConfigSent = true
				}
				//log.Print(json_string)
				rigs[i].PublishData(c, json_string)
				log.Printf("%s", rigs[i].ID)
			}
		}(i)

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

var mqtt_status = 0

func init_mqtt() mqtt.Client {
	// for i := 0; i < 5; i++ {
	// 	text := fmt.Sprintf("this is msg #%d!", i)
	// 	token := c.Publish("go-mqtt/sample", 0, false, text)
	// 	token.Wait()
	// }

	//fmt.Printf("wait -- ")
	hostname, _ := os.Hostname()

	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//mqtt.ERROR = log.New(os.Stdout, "", 0)

	// opts := mqtt.NewClientOptions().AddBroker("tcp://iot.wiin.win:11883").SetClientID("Status-" + hostname)
	// opts.SetKeepAlive(2 * time.Second)
	// //opts.SetDefaultPublishHandler(f)
	// opts.SetPingTimeout(1 * time.Second)
	// opts.SetUsername("mao")
	// opts.SetPassword("linmao8888")

	server := flag.String("server", "tls://w.wiin.win:1884", "The full url of the MQTT server to connect to ex: tls://w.wiin.win:1884")
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
		log.Println("mqtt OnConnect")
		// topics := map[string]byte{"haworkshopyc1/sensor/w0004/state": 0,
		// 	"haworkshopyc1/sensor/w0005/state":         0,
		// 	"haworkshopyc1/sensor/w0007/state":         0,
		// 	"haworkshopyc1/sensor/powermeteryc1/state": 0,
		// 	"haworkshopyc1/sensor/solaryc1/state":      0}
		// if token := c.SubscribeMultiple(topics, onMessageReceived); token.Wait() && token.Error() != nil {
		// 	//if token := c.Subscribe(*topic, byte(*qos), onMessageReceived); token.Wait() && token.Error() != nil {
		// 	panic(token.Error())
		// }

		// for _, rig := range rigs {
		// 	json_string, err := rig.GetStat()
		// 	if err != nil {
		// 		log.Println(err)
		// 	} else {
		// 		rig.PublishConfig(c, json_string)
		// 	}
		// }

		mqtt_status = 1
	}
	connOpts.OnConnectionLost = func(c mqtt.Client, err error) {
		mqtt_status = 0
		log.Println("mqtt OnConnectionLost")
	}
	connOpts.OnReconnecting = func(c mqtt.Client, co *mqtt.ClientOptions) {
		mqtt_status = 2
		log.Println("mqtt OnReconnecting")
	}

	return mqtt.NewClient(connOpts)
}

func init() {
	log.SetFlags(0)

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

	mqtt_client := init_mqtt()

loop:
	for {
		//log.Printf("mqtt status:%d\n", mqtt_status)

		if !mqtt_client.IsConnected() {
			if token := mqtt_client.Connect(); token.Wait() && token.Error() != nil {
				log.Println(token.Error())
			}
		} else {
			do_job(mqtt_client)
		}

		select {
		case <-time.After(time.Millisecond * 10000):
		case <-done:
			log.Println("job stopped, exting...")
			break loop
		}
	}

	mqtt_client.Disconnect(250)

}
