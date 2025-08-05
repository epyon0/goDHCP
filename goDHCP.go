package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	//toml "github.com/BurntSushi/toml"
	utils "github.com/epyon0/goUtils"
	toml "github.com/pelletier/go-toml/v2"
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
	PoolStart  string
	PoolEnd    string
	ServerPort uint16
	ClientPort uint16
}

type fieldsConfig struct {
	OP         byte
	HTYPE      byte
	HLEN       byte
	HOPS       byte
	XID        byte
	FLAGS      byte
	CIADDR     string
	YIADDR     string
	SIADDR     string
	GIADDR     string
	CHADDR     string
	CHADDR_RAW [16]byte
	SNAME      string
	FILE       string
}

type optionsConfig struct {
	PAD    bool
	END    bool
	SNM    string
	TOFF   uint32
	RTR    []string
	TSVR   []string
	NS     []string
	DNS    []string
	LSVR   []string
	CSVR   []string
	LPRS   []string
	ISVR   []string
	RLS    []string
	HN     string
	BFS    uint16
	MDF    string
	DN     string
	SS     string
	RP     string
	EP     string
	IPFW   bool
	NLSR   bool
	PF     [][2]string
	MDRS   uint16
	IPTTL  uint16
	PMAT   uint16
	PMPT   []uint16
	IMTU   uint16
	ASAL   bool
	BCADDR string
	PMD    bool
	MSUP   bool
	PRD    bool
	RSA    string
	SRT    [][2]string
	TENCP  bool
	ACTIM  uint32
	EENCP  bool
	TDTTL  uint8
	TKAI   uint32
	TKAG   bool
	NISD   string
	NISVR  []string
	NTPS   []string
	VSI    []byte
	NBNS   []string
	NBDDS  []string
	NBNT   uint8
	NBS    string
	XWSFS  []string
	XWSDM  []string
	NISPD  string
	NISPS  []string
	MIPHA  []string
	SMTPS  []string
	POP3S  []string
	NNTPS  []string
	DWWWS  []string
	DFS    []string
	DIRCS  []string
	STS    []string
	STDAS  []string
	RIPA   string
	IPALT  uint32
	OPTOVR uint8
	TFTPSN string
	BFNAME string
	MSGTYP uint8
	SVRID  string
	PRL    []byte
	MSG    string
	MAXMSG uint16
	T1     uint32
	T2     uint32
	VCID   []byte
	CIDENT []byte
}

type tomlConfig struct {
	Server  serverConfig  `toml:"server"`
	Fields  fieldsConfig  `toml:"fields"`
	Options optionsConfig `toml:"options"`
}

var configData tomlConfig

func PrintData() {
	utils.Debug(fmt.Sprintf("Configuration:\n%+v", configData), *debug)
	utils.Debug(fmt.Sprintf("Options:\n%s", utils.WalkByteSlice(BuildOptions())), *debug)
}

func GetIpSlice(ips []string) []byte {
	var output []byte
	for i := 0; i < len(ips); i++ {
		ip, err := utils.Ip2Uint32(ips[i])
		utils.Er(err)
		output = append(output, byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
	}
	return output
}

func BuildOptions() []byte {
	var output []byte

	ip, err := utils.Ip2Uint32(configData.Options.SNM)
	utils.Er(err)
	if ip != 0 {
		output = append(output, 1, 4, byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
	}

	if configData.Options.TOFF != 0 {
		output = append(output, 2, 4, byte(configData.Options.TOFF>>24), byte(configData.Options.TOFF>>16), byte(configData.Options.TOFF>>8), byte(configData.Options.TOFF))
	}

	if len(configData.Options.RTR) != 0 {
		output = append(output, 3, byte(4*len(configData.Options.RTR)))
		output = append(output, GetIpSlice(configData.Options.RTR)...)
	}

	if len(configData.Options.TSVR) != 0 {
		output = append(output, 4, byte(4*len(configData.Options.TSVR)))
		output = append(output, GetIpSlice(configData.Options.TSVR)...)
	}

	if len(configData.Options.NS) != 0 {
		output = append(output, 5, byte(4*len(configData.Options.NS)))
		output = append(output, GetIpSlice(configData.Options.NS)...)
	}

	if len(configData.Options.DNS) != 0 {
		output = append(output, 6, byte(4*len(configData.Options.DNS)))
		output = append(output, GetIpSlice(configData.Options.DNS)...)
	}

	if len(configData.Options.LSVR) != 0 {
		output = append(output, 7, byte(4*len(configData.Options.LSVR)))
		output = append(output, GetIpSlice(configData.Options.LSVR)...)
	}

	if len(configData.Options.CSVR) != 0 {
		output = append(output, 8, byte(4*len(configData.Options.CSVR)))
		output = append(output, GetIpSlice(configData.Options.CSVR)...)
	}

	if len(configData.Options.LPRS) != 0 {
		output = append(output, 9, byte(4*len(configData.Options.LPRS)))
		output = append(output, GetIpSlice(configData.Options.LPRS)...)
	}

	if len(configData.Options.ISVR) != 0 {
		output = append(output, 10, byte(4*len(configData.Options.ISVR)))
		output = append(output, GetIpSlice(configData.Options.ISVR)...)
	}

	if len(configData.Options.RLS) != 0 {
		output = append(output, 11, byte(4*len(configData.Options.RLS)))
		output = append(output, GetIpSlice(configData.Options.RLS)...)
	}

	if len(configData.Options.HN) != 0 {
		output = append(output, 12, byte(len(configData.Options.HN)))
		output = append(output, []byte(configData.Options.HN)...)
	}

	if configData.Options.BFS != 0 {
		output = append(output, 13, 2, byte(configData.Options.BFS>>8), byte(configData.Options.BFS))
	}

	if len(configData.Options.MDF) != 0 {
		output = append(output, 14, byte(len(configData.Options.MDF)))
		output = append(output, []byte(configData.Options.MDF)...)
	}

	if len(configData.Options.DN) != 0 {
		output = append(output, 15, byte(len(configData.Options.DN)))
		output = append(output, []byte(configData.Options.DN)...)
	}

	ip, err = utils.Ip2Uint32(configData.Options.SS)
	utils.Er(err)
	if ip != 0 {
		output = append(output, 16, 4, byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
	}

	if len(configData.Options.RP) != 0 {
		output = append(output, 17, byte(len(configData.Options.RP)))
		output = append(output, []byte(configData.Options.RP)...)
	}

	if len(configData.Options.EP) != 0 {
		output = append(output, 18, byte(len(configData.Options.EP)))
		output = append(output, []byte(configData.Options.EP)...)
	}

	output = append(output, 19, 1) // always create values for bools?
	if configData.Options.IPFW {
		output = append(output, 1)
	} else {
		output = append(output, 0)
	}

	return output
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
	config, err := os.ReadFile(*configFile)
	utils.Er(err)

	err = toml.Unmarshal(config, &configData)
	utils.Er(err)

	PrintData()
}
