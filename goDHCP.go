package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	toml "github.com/BurntSushi/toml"
	utils "github.com/epyon0/goUtils"
	//toml "github.com/pelletier/go-toml/v2"
)

/* Example Go Route:

func testFunc(j string, wg *sync.WaitGroup) {
	fmt.Printf("TEST:%s\r\n", j)
defer wg.Done()
}

	var wg sync.WaitGroup
	wg.Add(len(ips))
	for i := 0; i < len(n); i++ {


		go testFunc(ip, &wg)

	}
	wg.Wait()
*/

var debug *bool
var configFile, start, stop *string
var sport, cport *uint

type serverConfig struct {
	Poolstart  string
	Poolend    string
	Serverport uint16
	Clientport uint16
}

type tomlConfig struct {
	//server serverConfig `toml:"server"`
	poolStart string `toml:"poolStart"`
}

var configData tomlConfig

func PrintData() {
	utils.Debug(fmt.Sprintf("Pool Start: %s", configData.poolStart), *debug)
	utils.Debug(fmt.Sprintf("%+v", configData), *debug)
}

func main() {
	filePath, err := os.Executable()
	utils.Er(err)

	configFile = flag.String("config", fmt.Sprintf("%s/config.toml", filepath.Dir(filePath)), "Path to configuration file")
	start = flag.String("start", "192.168.0.100", "IPv4 address of the start of the DHCP pool")
	stop = flag.String("stop", "192.168.0.200", "IPv4 address of the end of the DHCP pool")
	sport = flag.Uint("sport", 67, "Server port")
	cport = flag.Uint("cport", 68, "Client port")
	debug = flag.Bool("v", false, "Enable verbose output")
	flag.Parse()

	_, err = os.Stat(*configFile)
	utils.Er(err)
	//config, err := os.ReadFile(*configFile)
	//utils.Er(err)

	//	toml.Unmarshal(config, &configData)
	//toml.Decode(string(config), &configData)
	mData, err := toml.DecodeFile(*configFile, &configData)
	utils.Er(err)
	fmt.Println(fmt.Sprintf("%+v", mData.Keys()))

	/*
		args, err := utils.GetArgs(0)
		utils.Er(err)
			_, err = os.Stat(data.configFile)
			utils.Er(err)
			file, err := os.ReadFile(data.configFile)
			utils.Er(err)
			json.NewDecoder(bytes.NewBuffer(file)).Decode(&data)
	*/

	PrintData()
}
