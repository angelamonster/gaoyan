package gaoyan

import (
	"encoding/json"
	"fmt"
	"log"

	modbus "github.com/goburrow/modbus"
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
	V       [3]float64
}

func (m METER) Read(host string, port int, addr int) (json_string string, err error) {
	client := modbus.TCPClient(fmt.Sprintf("%s:%d", host, port))
	// Read input register 9
	var count uint16 = 0x1E + 1
	results, err := client.ReadInputRegisters(0, count)
	if err != nil {
		log.Println(fmt.Sprintf("Connection error: %s", err))
		return "", err
	} else {
		// log.Println(fmt.Sprintf("result: %s", results))

		// for _, b := range results {
		// 	log.Printf("%02x", b)
		// }
		// log.Println("")
		var v uint16 = uint16(results[0])<<8 + uint16(results[1])
		m.V[0] = float64(v) * 0.1
		m.V[1] = float64(results[1]) * 0.1
		m.V[2] = float64(results[2]) * 0.1
		m.I[0] = float64(results[0x03+0]) * 0.01
		m.I[1] = float64(results[0x03+1]) * 0.01
		m.I[2] = float64(results[0x03+2]) * 0.01
		m.PTotal = float64(results[0x07])
		m.P[0] = float64(results[0x08+0])
		m.P[1] = float64(results[0x08+1])
		m.P[2] = float64(results[0x08+2])
		m.RPTotal = float64(results[0x0b])
		m.RP[0] = float64(results[0x0c+0])
		m.RP[1] = float64(results[0x0c+1])
		m.RP[2] = float64(results[0x0c+2])
		m.APTotal = float64(results[0x0f])
		m.AP[0] = float64(results[0x10+0])
		m.AP[1] = float64(results[0x10+1])
		m.AP[2] = float64(results[0x10+2])
		m.F[0] = float64(results[0x1a+0]) * 0.01
		m.F[1] = float64(results[0x1a+1]) * 0.01
		m.F[2] = float64(results[0x1a+2]) * 0.01

		json_bytes, err := json.Marshal(m)

		if err != nil {
			log.Printf("json error: %s", err)
			return "", err
		}
		return string(json_bytes), nil

	}

	// //host := "127.0.0.1"
	// //port := modbusclient.MODBUS_PORT

	// // turn on the debug trace option, to see what is being transmitted
	// trace := true

	// conn, cerr := modbusclient.ConnectTCP(host, port)
	// if cerr != nil {
	// 	log.Println(fmt.Sprintf("Connection error: %s", cerr))
	// } else {

	// 	// attempt to read one (0x01) holding registers starting at address 200
	// 	var size int = 0x1E + 1
	// 	addr = 0x00
	// 	readData := make([]byte, size)
	// 	readData[0] = byte(addr >> 8)   // (High Byte)
	// 	readData[1] = byte(addr & 0xff) // (Low Byte)
	// 	//[2] = 0x01
	// 	readData[2] = 0x1E + 1
	// 	// count := 0x1E + 1
	// 	//var U_INT byte = 0x1
	// 	//													# address count unit
	// 	//            request = client.read_holding_registers(0, 0x1E+1,unit=self.UNIT)
	// 	// make this read request transaction id 1, with a 300 millisecond tcp timeout
	// 	readResult, readErr := modbusclient.TCPRead(conn, 300, 1, modbusclient.FUNCTION_READ_HOLDING_REGISTERS, false, 0x1, readData, trace)

	// 	//readResult, readErr := modbusclient.TCPRead(conn, 300, 1, modbusclient.FUNCTION_READ_HOLDING_REGISTERS, false, 0x00, readData, trace)
	// 	if readErr != nil {
	// 		log.Println(readErr)
	// 	}
	// 	log.Println(readResult)

	// 	// // attempt to write to a single coil at address 300
	// 	// writeData := make([]byte, 3)
	// 	// writeData[0] = byte(300 >> 8)   // (High Byte)
	// 	// writeData[1] = byte(300 & 0xff) // (Low Byte)
	// 	// writeData[2] = 0xff             // 0xff turns the coil on; 0x00 turns the coil off
	// 	// // make this read request transaction id 2, with a 300 millisecond tcp timeout
	// 	// writeResult, writeErr := modbusclient.TCPWrite(conn, 300, 2, modbusclient.FUNCTION_WRITE_SINGLE_COIL, false, 0x00, writeData, trace)
	// 	// if writeErr != nil {
	// 	// 	log.Println(writeErr)
	// 	// }
	// 	// log.Println(writeResult)

	// 	modbusclient.DisconnectTCP(conn)
	// }

}
