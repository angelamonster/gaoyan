//go:build ignore

package main

import (
	"crypto/tls"
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
	//log.Println("loop")

	for i, _ := range rigs {
		log.Printf("%s timestamp gap: %ds\n", rigs[i].ID, time.Now().Unix()-rigs[i].BusyTimeStamp)
		if rigs[i].Busy == false || time.Now().Unix()-rigs[i].BusyTimeStamp > 30 {
			go func(i int) {
				rigs[i].Busy = true
				rigs[i].BusyTimeStamp = time.Now().Unix()
				mi, err := rigs[i].GetStat()
				if err == nil {
					if false == rigs[i].ConfigSent {
						rigs[i].PublishConfig(c, mi)
						rigs[i].ConfigSent = true
					}
					//log.Print(json_string)
					rigs[i].PublishData(c, mi)
					log.Printf("%s - %dMH - %d\n", rigs[i].ID, mi.MainCrypto.HashRate, mi.HighTemp)
				} else {
					log.Printf("%s getstat error:%s\n", rigs[i].ID, err.Error())
				}

				rigs[i].Busy = false
			}(i)
		}
	}
}

var mqtt_status = 0

func init_mqtt() mqtt.Client {
	hostname, _ := os.Hostname()

	connOpts := mqtt.NewClientOptions()
	connOpts.AddBroker("tls://w.wiin.win:1884")
	connOpts.SetClientID(hostname + "-rig-reporter")
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
