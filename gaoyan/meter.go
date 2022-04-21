package gaoyan

import (
	"fmt"
	"log"

	modbusclient "github.com/dpapathanasiou/go-modbus"
)

//"encoding/json"
// "fmt"
// "log"

// claymore "../rpcclaymore"
// mqtt "github.com/eclipse/paho.mqtt.golang"

type METER struct {
	I       [3]float64
	ITotal  float64
	P       [3]float64
	PTotal  float64
	AP      [3]float64
	APTotal float64
	RP      [3]float64
	RPTotal float64
	F       [3]float64
	E       float64
}

func (m METER) Read(host string, port int, addr int) {
	//host := "127.0.0.1"
	//port := modbusclient.MODBUS_PORT

	// turn on the debug trace option, to see what is being transmitted
	trace := true

	conn, cerr := modbusclient.ConnectTCP(host, port)
	if cerr != nil {
		log.Println(fmt.Sprintf("Connection error: %s", cerr))
	} else {

		// attempt to read one (0x01) holding registers starting at address 200
		addr = 0x00
		readData := make([]byte, 3)
		readData[0] = byte(addr >> 8)   // (High Byte)
		readData[1] = byte(addr & 0xff) // (Low Byte)
		readData[2] = 0x01

		// count := 0x1E + 1
		// unit := 0x1
		//													# address count unit
		//            request = client.read_holding_registers(0, 0x1E+1,unit=self.UNIT)
		// make this read request transaction id 1, with a 300 millisecond tcp timeout
		readResult, readErr := modbusclient.TCPRead(conn, 300, 1, modbusclient.FUNCTION_READ_HOLDING_REGISTERS, false, 0x00, readData, trace)
		if readErr != nil {
			log.Println(readErr)
		}
		log.Println(readResult)

		// // attempt to write to a single coil at address 300
		// writeData := make([]byte, 3)
		// writeData[0] = byte(300 >> 8)   // (High Byte)
		// writeData[1] = byte(300 & 0xff) // (Low Byte)
		// writeData[2] = 0xff             // 0xff turns the coil on; 0x00 turns the coil off
		// // make this read request transaction id 2, with a 300 millisecond tcp timeout
		// writeResult, writeErr := modbusclient.TCPWrite(conn, 300, 2, modbusclient.FUNCTION_WRITE_SINGLE_COIL, false, 0x00, writeData, trace)
		// if writeErr != nil {
		// 	log.Println(writeErr)
		// }
		// log.Println(writeResult)

		modbusclient.DisconnectTCP(conn)
	}

}
