package gaoyan

import (
	//modbus "github.com/goburrow/modbus"
	//modbusclient "github.com/dpapathanasiou/go-modbus"
	"encoding/json"
	"fmt"

	modbus "github.com/thinkgos/gomodbus"
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

func build_value(results []byte, pos int) float64 {

	var high byte = results[pos]
	var low byte = results[pos+1]

	var v uint16 = uint16(high)<<8 + uint16(low)
	var s uint32 = 0x10000

	if v < 0x8000 {
		return float64(v)
	} else {
		return float64(uint32(v) - s)
	}

}

func (m METER) Read(host string, port int) (json_string string, err error) {
	p := modbus.NewTCPClientProvider(fmt.Sprintf("%s:%d", host, port), modbus.WithTCPTimeout(5000)) //|
	client := modbus.NewClient(p)
	cerr := client.Connect()
	if cerr != nil {
		fmt.Println("connect failed, ", cerr)
		return "", err
	} else {
		fmt.Println("starting")

		//var tx modbus.ProtocolDataUnit =

		var length uint16 = 0x1E + 1
		//results, err := client.Send()(0x01, 0x00, length)
		results, err := client.ReadInputRegisters(0x01, 0x00, length)
		//_, err := client.ReadCoils(1, 0, 10)
		if err != nil {
			fmt.Println(err.Error())
			return "", err
		} else {
			m.V[0] = float64(results[0]) * 0.1
			m.V[1] = float64(results[1]) * 0.1
			m.V[2] = float64(results[2]) * 0.1
			m.I[0] = float64(results[0x03+0]) * 0.01
			m.I[1] = float64(results[0x03+1]) * 0.01
			m.I[2] = float64(results[0x03+2]) * 0.01
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
		}

	}
	defer client.Close()

	//	fmt.Printf("ReadDiscreteInputs %#v\r\n", results)
	json_bytes, err := json.Marshal(m)

	return string(json_bytes), nil

	// // turn on the debug trace option, to see what is being transmitted
	// trace := true

	// conn, cerr := modbusclient.ConnectTCP(host, port)
	// if cerr != nil {
	// 	log.Printf("Connection error: %s", cerr)
	// 	return "", cerr
	// } else {

	// 	// attempt to read one (0x01) holding registers starting at address 200
	// 	//var size int = 0x1E + 1
	// 	addr := 0x00
	// 	readData := make([]byte, 3)
	// 	readData[0] = byte(addr >> 8)   // (High Byte)
	// 	readData[1] = byte(addr & 0xff) // (Low Byte)
	// 	//[2] = 0x01
	// 	readData[2] = byte(0x1E + 1)
	// 	// count := 0x1E + 1
	// 	//var U_INT byte = 0x1
	// 	//													# address count unit
	// 	//            request = client.read_holding_registers(0, 0x1E+1,unit=self.UNIT)
	// 	// make this read request transaction id 1, with a 300 millisecond tcp timeout
	// 	readResult, readErr := modbusclient.TCPRead(conn, 3000, 1, modbusclient.FUNCTION_READ_INPUT_REGISTERS, false, 0x01, readData, trace)

	// 	//readResult, readErr := modbusclient.TCPRead(conn, 300, 1, modbusclient.FUNCTION_READ_HOLDING_REGISTERS, false, 0x00, readData, trace)
	// 	if readErr != nil {
	// 		log.Println(readErr)
	// 	}
	// 	log.Println("readResult")
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

	// 	return "", nil
	// }

	// client := modbus.TCPClient(fmt.Sprintf("%s:%d", host, port))
	// // Read input register 9
	// var count uint16 = 0x1E + 1
	// results, err := client.ReadInputRegisters(0, count)

	// if err != nil {
	// 	log.Printf("Connection error: %s", err)
	// 	return "", err
	// } else {
	// 	// log.Println(fmt.Sprintf("result: %s", results))

	// 	// for _, b := range results {
	// 	// 	log.Printf("%02x", b)
	// 	// }
	// 	// log.Println("")

	// 	m.V[0] = build_value(results, 0) * 0.1
	// 	m.V[1] = build_value(results, 2) * 0.1
	// 	m.V[2] = build_value(results, 4) * 0.1
	// 	m.I[0] = build_value(results, (0x03+0)*2) * 0.01
	// 	m.I[1] = build_value(results, (0x03+1)*2) * 0.01
	// 	m.I[2] = build_value(results, (0x03+2)*2) * 0.01
	// 	m.I[0] = build_value(results, (0x03+0)*2) * 0.01
	// 	m.I[1] = build_value(results, (0x03+1)*2) * 0.01
	// 	m.I[2] = build_value(results, (0x03+2)*2) * 0.01
	// 	m.PTotal = build_value(results, (0x07)*2)
	// 	m.P[0] = build_value(results, (0x08+0)*2)
	// 	m.P[1] = build_value(results, (0x08+1)*2)
	// 	m.P[2] = build_value(results, (0x08+2)*2)
	// 	m.RPTotal = build_value(results, (0x0b)*2)
	// 	m.RP[0] = build_value(results, (0x0c+0)*2)
	// 	m.RP[1] = build_value(results, (0x0c+1)*2)
	// 	m.RP[2] = build_value(results, (0x0c+2)*2)
	// 	m.APTotal = build_value(results, (0x0f)*2)
	// 	m.AP[0] = build_value(results, (0x10+0)*2)
	// 	m.AP[1] = build_value(results, (0x10+1)*2)
	// 	m.AP[2] = build_value(results, (0x10+2)*2)
	// 	m.F[0] = build_value(results, (0x1a+0)*2) * 0.01
	// 	m.F[1] = build_value(results, (0x1a+1)*2) * 0.01
	// 	m.F[2] = build_value(results, (0x1a+2)*2) * 0.01

	// 	json_bytes, err := json.Marshal(m)

	// 	if err != nil {
	// 		log.Printf("json error: %s", err)
	// 		return "", err
	// 	}
	// 	return string(json_bytes), nil

	// }

}
