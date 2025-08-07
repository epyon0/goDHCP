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
	IPTTL  uint8
	PMAT   uint32
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
	opts := BuildOptions()
	utils.Debug(fmt.Sprintf("Options:\n%s", utils.WalkByteSlice(opts)), *debug)
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

	if len(configData.Options.SNM) != 0 {
		ip, err := utils.Ip2Uint32(configData.Options.SNM)
		utils.Er(err)
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

	if len(configData.Options.SS) != 0 {
		ip, err := utils.Ip2Uint32(configData.Options.SS)
		utils.Er(err)
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

	output = append(output, 19, 1)
	if configData.Options.IPFW {
		output = append(output, 1)
	} else {
		output = append(output, 0)
	}

	output = append(output, 20, 1)
	if configData.Options.NLSR {
		output = append(output, 1)
	} else {
		output = append(output, 0)
	}

	if len(configData.Options.PF) != 0 {
		output = append(output, 21, byte(8*len(configData.Options.PF)))
		for i := 0; i < len(configData.Options.PF); i++ {
			ip, err := utils.Ip2Uint32(configData.Options.PF[i][0])
			utils.Er(err)
			nm, err := utils.Ip2Uint32(configData.Options.PF[i][1])
			utils.Er(err)
			output = append(output, byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
			output = append(output, byte(nm>>24), byte(nm>>16), byte(nm>>8), byte(nm))
		}
	}

	if configData.Options.MDRS != 0 {
		output = append(output, 22, 2, byte(configData.Options.MDRS>>8), byte(configData.Options.MDRS))
	}

	if configData.Options.IPTTL != 0 {
		output = append(output, 23, 1, configData.Options.IPTTL)
	}

	if configData.Options.PMAT != 0 {
		output = append(output, 24, 4, byte(configData.Options.PMAT>>24), byte(configData.Options.PMAT>>16), byte(configData.Options.PMAT>>8), byte(configData.Options.PMAT))
	}

	if len(configData.Options.PMPT) != 0 {
		output = append(output, 25, byte(2*len(configData.Options.PMPT)))
		for i := 0; i < len(configData.Options.PMPT); i++ {
			tmp16 := configData.Options.PMPT[i]
			output = append(output, byte(tmp16>>8), byte(tmp16))
		}
	}

	if configData.Options.IMTU != 0 {
		output = append(output, 26, 2, byte(configData.Options.IMTU>>8), byte(configData.Options.IMTU))
	}

	output = append(output, 27, 1)
	if configData.Options.ASAL {
		output = append(output, 1)
	} else {
		output = append(output, 0)
	}

	if len(configData.Options.BCADDR) != 0 {
		ip, err := utils.Ip2Uint32(configData.Options.BCADDR)
		utils.Er(err)
		output = append(output, 28, 4, byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
	}

	output = append(output, 29, 1)
	if configData.Options.PMD {
		output = append(output, 1)
	} else {
		output = append(output, 0)
	}

	output = append(output, 30, 1)
	if configData.Options.MSUP {
		output = append(output, 1)
	} else {
		output = append(output, 0)
	}

	output = append(output, 31, 1)
	if configData.Options.PRD {
		output = append(output, 1)
	} else {
		output = append(output, 0)
	}

	if len(configData.Options.RSA) != 0 {
		ip, err := utils.Ip2Uint32(configData.Options.RSA)
		utils.Er(err)
		output = append(output, 32, 4, byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
	}

	if len(configData.Options.SRT) != 0 {
		output = append(output, 33, byte(8*len(configData.Options.SRT)))
		for i := 0; i < len(configData.Options.SRT); i++ {
			ip, err := utils.Ip2Uint32(configData.Options.SRT[i][0])
			utils.Er(err)
			nm, err := utils.Ip2Uint32(configData.Options.SRT[i][1])
			utils.Er(err)
			output = append(output, byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
			output = append(output, byte(nm>>24), byte(nm>>16), byte(nm>>8), byte(nm))
		}
	}

	output = append(output, 34, 1)
	if configData.Options.TENCP {
		output = append(output, 1)
	} else {
		output = append(output, 0)
	}

	if configData.Options.ACTIM != 0 {
		tmp32 := configData.Options.ACTIM
		output = append(output, 35, 4, byte(tmp32>>24), byte(tmp32>>16), byte(tmp32>>8), byte(tmp32))
	}

	output = append(output, 36, 1)
	if configData.Options.EENCP {
		output = append(output, 1)
	} else {
		output = append(output, 0)
	}

	if configData.Options.TDTTL != 0 {
		output = append(output, 37, 1, configData.Options.TDTTL)
	}

	if configData.Options.TKAI != 0 {
		tmp32 := configData.Options.TKAI
		output = append(output, 38, 4, byte(tmp32>>24), byte(tmp32>>16), byte(tmp32>>8), byte(tmp32))
	}

	output = append(output, 39, 1)
	if configData.Options.TKAG {
		output = append(output, 1)
	} else {
		output = append(output, 0)
	}

	if len(configData.Options.NISD) != 0 {
		output = append(output, 40, byte(len(configData.Options.NISD)))
		output = append(output, []byte(configData.Options.NISD)...)
	}

	if len(configData.Options.NISVR) != 0 {
		output = append(output, 41, byte(4*len(configData.Options.NISVR)))
		output = append(output, GetIpSlice(configData.Options.NISVR)...)
	}

	if len(configData.Options.NTPS) != 0 {
		output = append(output, 42, byte(4*len(configData.Options.NTPS)))
		output = append(output, GetIpSlice(configData.Options.NTPS)...)
	}

	if len(configData.Options.VSI) != 0 {
		output = append(output, 43, byte(len(configData.Options.VSI)))
		for i := 0; i < len(configData.Options.VSI); i++ {
			output = append(output, configData.Options.VSI[i])
		}
	}

	if len(configData.Options.NBNS) != 0 {
		output = append(output, 44, byte(4*len(configData.Options.NBNS)))
		output = append(output, GetIpSlice(configData.Options.NBNS)...)
	}

	if len(configData.Options.NBDDS) != 0 {
		output = append(output, 45, byte(4*len(configData.Options.NBDDS)))
		output = append(output, GetIpSlice(configData.Options.NBDDS)...)
	}

	if configData.Options.NBNT != 0 {
		output = append(output, 46, configData.Options.NBNT)
	}

	if len(configData.Options.NBS) != 1 {
		output = append(output, 47, byte(len(configData.Options.NBS)))
		output = append(output, []byte(configData.Options.NBS)...)
	}

	if len(configData.Options.XWSFS) != 0 {
		output = append(output, 48, byte(4*len(configData.Options.XWSFS)))
		output = append(output, GetIpSlice(configData.Options.XWSFS)...)
	}

	if len(configData.Options.XWSDM) != 0 {
		output = append(output, 49, byte(4*len(configData.Options.XWSDM)))
		output = append(output, GetIpSlice(configData.Options.XWSDM)...)
	}

	if len(configData.Options.NISPD) != 0 {
		output = append(output, 64, byte(len(configData.Options.NISPD)))
		output = append(output, []byte(configData.Options.NISPD)...)
	}

	if len(configData.Options.NISPS) != 0 {
		output = append(output, 65, byte(4*len(configData.Options.NISPS)))
		output = append(output, GetIpSlice(configData.Options.NISPS)...)
	}

	if len(configData.Options.MIPHA) != 0 {
		output = append(output, 68, byte(4*len(configData.Options.MIPHA)))
		output = append(output, GetIpSlice(configData.Options.MIPHA)...)
	}

	if len(configData.Options.SMTPS) != 0 {
		output = append(output, 69, byte(4*len(configData.Options.SMTPS)))
		output = append(output, GetIpSlice(configData.Options.SMTPS)...)
	}

	if len(configData.Options.POP3S) != 0 {
		output = append(output, 70, byte(4*len(configData.Options.POP3S)))
		output = append(output, GetIpSlice(configData.Options.POP3S)...)
	}

	if len(configData.Options.NNTPS) != 0 {
		output = append(output, 71, byte(4*len(configData.Options.NNTPS)))
		output = append(output, GetIpSlice(configData.Options.NNTPS)...)
	}

	if len(configData.Options.DWWWS) != 0 {
		output = append(output, 72, byte(4*len(configData.Options.DWWWS)))
		output = append(output, GetIpSlice(configData.Options.DWWWS)...)
	}

	if len(configData.Options.DFS) != 0 {
		output = append(output, 73, byte(4*len(configData.Options.DFS)))
		output = append(output, GetIpSlice(configData.Options.DFS)...)
	}

	if len(configData.Options.DIRCS) != 0 {
		output = append(output, 74, byte(4*len(configData.Options.DIRCS)))
		output = append(output, GetIpSlice(configData.Options.DIRCS)...)
	}

	if len(configData.Options.STS) != 0 {
		output = append(output, 75, byte(4*len(configData.Options.STS)))
		output = append(output, GetIpSlice(configData.Options.STS)...)
	}

	if len(configData.Options.STDAS) != 0 {
		output = append(output, 76, byte(4*len(configData.Options.STDAS)))
		output = append(output, GetIpSlice(configData.Options.STDAS)...)
	}

	if len(configData.Options.RIPA) != 0 {
		ip, err := utils.Ip2Uint32(configData.Options.RIPA)
		utils.Er(err)
		output = append(output, 50, 4, byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
	}

	if configData.Options.IPALT != 0 {
		tmp32 := configData.Options.IPALT
		output = append(output, 51, 4, byte(tmp32>>24), byte(tmp32>>16), byte(tmp32>>8), byte(tmp32))
	}

	if configData.Options.OPTOVR != 0 {
		output = append(output, 52, 1, configData.Options.OPTOVR)
	}

	if len(configData.Options.TFTPSN) != 0 {
		output = append(output, 66, byte(len(configData.Options.TFTPSN)))
		output = append(output, []byte(configData.Options.TFTPSN)...)
	}

	if len(configData.Options.BFNAME) != 0 {
		output = append(output, 67, byte(len(configData.Options.BFNAME)))
		output = append(output, []byte(configData.Options.BFNAME)...)
	}

	if len(configData.Options.SVRID) != 0 {
		ip, err := utils.Ip2Uint32(configData.Options.SVRID)
		utils.Er(err)
		output = append(output, 54, 4, byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
	}

	if len(configData.Options.PRL) != 0 {
		output = append(output, 55, byte(len(configData.Options.PRL)))
		for i := 0; i < len(configData.Options.PRL); i++ {
			output = append(output, configData.Options.PRL[i])
		}
	}

	if len(configData.Options.MSG) != 0 {
		output = append(output, 56, byte(len(configData.Options.MSG)))
		output = append(output, []byte(configData.Options.MSG)...)
	}

	if configData.Options.MAXMSG != 0 {
		tmp16 := configData.Options.MAXMSG
		output = append(output, 57, 2, byte(tmp16>>8), byte(tmp16))
	}

	if configData.Options.T1 != 0 {
		t1 := configData.Options.T1
		output = append(output, 58, 4, byte(t1>>24), byte(t1>>16), byte(t1>>8), byte(t1))
	}

	if configData.Options.T2 != 0 {
		t2 := configData.Options.T2
		output = append(output, 59, 4, byte(t2>>24), byte(t2>>16), byte(t2>>8), byte(t2))
	}

	if len(configData.Options.VCID) != 0 {
		output = append(output, 60, byte(len(configData.Options.VCID)))
		for i := 0; i < len(configData.Options.VCID); i++ {
			output = append(output, configData.Options.VCID[i])
		}
	}

	if len(configData.Options.CIDENT) != 0 {
		output = append(output, 61, byte(len(configData.Options.CIDENT)))
		for i := 0; i < len(configData.Options.CIDENT); i++ {
			output = append(output, configData.Options.CIDENT[i])
		}
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
