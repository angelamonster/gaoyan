package rpcclaymore

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"
	"strconv"
	"strings"
	"time"
)

const (
	methodGetInfo      = "miner_getstat1"
	methodRestartMiner = "miner_restart"
	methodReboot       = "miner_reboot"
)

var args = struct {
	id      string
	jsonrpc string
	psw     string
}{"0", "2.0", ""}

// Crypto Information about this concrete crypto-currency
type Crypto struct {
	HashRate       int `json:"hashrate"`
	Shares         int `json:"shares"`
	RejectedShares int `json:"rejected"`
	InvalidShares  int `json:"invalid"`
}

func (c Crypto) String() (s string) {
	if c.HashRate+c.RejectedShares+c.InvalidShares+c.Shares == 0 {
		return "Disabled\n"
	}
	s += fmt.Sprintf("HashRate:         %8d Mh/s\n", c.HashRate)
	s += fmt.Sprintf("Shares:           %8d\n", c.Shares)
	s += fmt.Sprintf("Rejected Shares:  %8d\n", c.RejectedShares)
	s += fmt.Sprintf("Invalid Shares:   %8d\n", c.InvalidShares)
	return s
}

// PoolInfo Information about the miner's connected pool
type PoolInfo struct {
	Address  string `json:"adress"`
	Switches int    `json:"switches"`
}

func (p PoolInfo) String() (s string) {
	if p.Address == "" {
		return "Disabled\n"
	}
	s += fmt.Sprintf("Address:   %s\n", p.Address)
	s += fmt.Sprintf("Switches:  %"+strconv.Itoa(len(p.Address))+"d\n", p.Switches)
	return s
}

// GPU Information about each concrete GPU
type GPU struct {
	HashRate    int `json:"hashrate"`
	AltHashRate int `json:"althashrate"`
	Temperature int `json:"temperature"`
	FanSpeed    int `json:"fanspeed"`
}

func (gpu GPU) String() (s string) {
	s += fmt.Sprintf("Hash Rate:     %8d Mh/s\n", gpu.HashRate)
	s += fmt.Sprintf("Alt Hash Rate: %8d Mh/s\n", gpu.AltHashRate)
	s += fmt.Sprintf("Temperature:   %8d º\n", gpu.Temperature)
	s += fmt.Sprintf("Fan Speed:     %8d %%\n", gpu.FanSpeed)
	return s
}

// IsStuck Return true if the GPU is not mining
func (gpu GPU) IsStuck() bool {
	return gpu.HashRate == 0
}

// MinerInfo Information about the miner
type MinerInfo struct {
	Version    string   `json:"version"`
	UpTime     int      `json:"uptime"`
	MainCrypto Crypto   `json:"maincrypto"`
	AltCrypto  Crypto   `json:"altcrypto"`
	MainPool   PoolInfo `json:"mainpool"`
	AltPool    PoolInfo `json:"altpool"`
	Timestamp  int64    `json:"timestamp"`
	GPUS       []GPU
}

// StuckGPUs Return the number of GPUs that are not mining
func (m MinerInfo) StuckGPUs() int {
	var total int
	for _, gpu := range m.GPUS {
		if gpu.IsStuck() {
			total++
		}
	}
	return total
}

func (m MinerInfo) String() string {
	var s string
	s += fmt.Sprintf("Version:   %10s\n", m.Version)
	s += fmt.Sprintf("Up Time:   %10d min\n", m.UpTime)
	s += "\n"
	s += fmt.Sprintf("Main Crypto\n%s\n", m.MainCrypto)
	s += fmt.Sprintf("Alt Crypto\n%s\n", m.AltCrypto)
	s += fmt.Sprintf("Main Pool\n%s\n", m.MainPool)
	s += fmt.Sprintf("Alt Pool\n%s\n", m.AltPool)
	for i, gpu := range m.GPUS {
		s += fmt.Sprintf("GPU %d\n%s\n", i, gpu)
	}
	return s
}

// Miner creates an instance to get info of a miner
type Miner struct {
	Address  string
	Password string
}

func (m Miner) String() (s string) {
	return fmt.Sprintf("Miner {Address: %s}\n", m.Address)
}

// Restart Stop and start the miner
func (m Miner) Restart() error {
	client, err := jsonrpc.Dial("tcp", m.Address)
	if err != nil {
		return err
	}
	defer client.Close()
	args.psw = m.Password
	return client.Call(methodRestartMiner, args, nil)
}

// Reboot Turn off and on again the computer
func (m Miner) Reboot() error {
	client, err := jsonrpc.Dial("tcp", m.Address)
	if err != nil {
		return err
	}
	defer client.Close()
	args.psw = m.Password
	return client.Call(methodReboot, args, nil)
}

// GetInfo return the status of the miner or throw and error if it is not reachable
func (m Miner) GetInfo() (MinerInfo, error) {
	var mi MinerInfo
	var reply []string
	client, err := jsonrpc.Dial("tcp", m.Address)
	if err != nil {
		return mi, err
	}
	defer client.Close()
	args.psw = m.Password
	err = client.Call(methodGetInfo, args, &reply)
	if err != nil {
		return mi, err
	}

	log.Println(reply)

	return parseResponse(reply), nil
}

func parseResponse(info []string) MinerInfo {
	var mi MinerInfo
	var group []string

	mi.Timestamp = time.Now().Unix()

	mi.Version = strings.Replace(info[0], " - ETH", "", 1)
	mi.UpTime = toInt(info[1])

	group = splitGroup(info[2])
	mi.MainCrypto.HashRate = toInt(group[0])
	mi.MainCrypto.Shares = toInt(group[1])
	mi.MainCrypto.RejectedShares = toInt(group[2])

	group = splitGroup(info[4])
	mi.AltCrypto.HashRate = toInt(group[0])
	mi.AltCrypto.Shares = toInt(group[1])
	mi.AltCrypto.RejectedShares = toInt(group[2])

	group = splitGroup(info[7])
	mi.MainPool.Address = group[0]
	if len(group) > 1 {
		mi.AltPool.Address = group[1]
	}

	group = splitGroup(info[8])
	mi.MainCrypto.InvalidShares = toInt(group[0])
	mi.MainPool.Switches = toInt(group[1])
	mi.AltCrypto.InvalidShares = toInt(group[2])
	mi.AltPool.Switches = toInt(group[3])

	for _, hashrate := range splitGroup(info[3]) {
		mi.GPUS = append(mi.GPUS, GPU{HashRate: toInt(hashrate)})
	}

	for i, val := range splitGroup(info[6]) {
		if i%2 == 0 {
			mi.GPUS[i/2].Temperature = toInt(val)
		} else {
			mi.GPUS[(i-1)/2].FanSpeed = toInt(val)
		}
	}

	if mi.AltPool.Address != "" {
		for i, val := range splitGroup(info[5]) {
			hashrate, err := strconv.Atoi(val)
			if err == nil {
				mi.GPUS[i].AltHashRate = hashrate
			}
		}
	}

	return mi
}

func splitGroup(s string) []string {
	return strings.Split(s, ";")
}

func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
