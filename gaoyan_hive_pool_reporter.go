//go:build ignore

package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"./gaoyan"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var hive = gaoyan.HIVE{Address: "2c2c0bd495dffa493f8cd6ca71650aa6585a6207", ConfigSent: false}

func do_job(c mqtt.Client) {
	//log.Println("loop")
	go func() {

		info, err := hive.Read()
		if err != nil {
			log.Println(err)
		} else {
			json_bytes, err := json.Marshal(info)
			if err != nil {
				log.Println(err)
			} else {

				if false == hive.ConfigSent {
					log.Println("Publish Config")
					hive.PublishConfig(c)
					hive.ConfigSent = true
				}

				hive.PublishData(c, string(json_bytes))
			}
		}

	}()

}

var mqtt_status = 0

func init_mqtt() mqtt.Client {
	hostname, _ := os.Hostname()

	connOpts := mqtt.NewClientOptions()
	connOpts.AddBroker("tls://w.wiin.win:1884")
	connOpts.SetClientID(hostname + "-hive-pool-reporter")
	connOpts.SetCleanSession(true)
	connOpts.SetUsername("mao")
	connOpts.SetPassword("linmao8888")

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c mqtt.Client) {
		log.Println("mqtt OnConnect")
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
}

func main() {

	var cycle time.Duration
	cycle = 30 //s

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
		case <-time.After(time.Millisecond * 1000 * cycle):
		case <-done:
			log.Println("job stopped, exting...")
			break loop
		}
	}

	mqtt_client.Disconnect(250)

}
