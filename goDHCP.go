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
	config, err := os.ReadFile(*configFile)
	utils.Er(err)

	err = toml.Unmarshal(config, &configData)
	utils.Er(err)

	PrintData()
}
