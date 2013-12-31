package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"util"
	_ "github.com/mattn/go-sqlite3"
	"github.com/bitly/go-nsq"
	"cpu_usage"
	"mem_usage"
	"net_flow"
	"heartbeat"
	"accessibility"
	"register"
)

var (
	showVersion        = flag.Bool("version", false, "print version string")
	nsqdTCPAddrs       = util.StringArray{}
	lookupdHTTPAddrs   = util.StringArray{}
	maxInFlight        = flag.Int("max-in-flight", 200, "max number of messages to allow in flight")
	verbose            = flag.Bool("verbose", false, "enable verbose logging")
	maxBackoffDuration = flag.Duration("max-backoff-duration", 120*time.Second, "the maximum backoff duration")
	termChan chan os.Signal
	dbPath = flag.String("db-path", "../db/", "absolute or relative path for db.")
)

func init() {
	flag.Var(&nsqdTCPAddrs, "nsqd-tcp-address", "nsqd TCP address (may be given multiple times)")
	flag.Var(&lookupdHTTPAddrs, "lookupd-http-address", "lookupd HTTP address (may be given multiple times)")
}


func runCpuUsageClient(cuh *cpu_usage.CPUUsageHandler) (cuTodb *nsq.Reader, err error) {
	cuTodb, err = nsq.NewReader("cpu_usage", "multidb")
	if err != nil {
		log.Fatalf(err.Error())
	}
	cuTodb.SetMaxInFlight(*maxInFlight)
	cuTodb.SetMaxBackoffDuration(*maxBackoffDuration)
	cuTodb.VerboseLogging = *verbose
	cuTodb.AddHandler(cuh)
	fmt.Println(nsqdTCPAddrs)
	for _, addrString := range nsqdTCPAddrs {
		err := cuTodb.ConnectToNSQ(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	for _, addrString := range lookupdHTTPAddrs {
		log.Printf("lookupd addr %s", addrString)
		err := cuTodb.ConnectToLookupd(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	return
}

func runMemUsageClient(muh *mem_usage.MemUsageHandler) (muTodb *nsq.Reader, err error) {
	muTodb, err = nsq.NewReader("mem_usage", "multidb")
	if err != nil {
		log.Fatalf(err.Error())
	}
	muTodb.SetMaxInFlight(*maxInFlight)
	muTodb.SetMaxBackoffDuration(*maxBackoffDuration)
	muTodb.VerboseLogging = *verbose
	muTodb.AddHandler(muh)

	fmt.Println(nsqdTCPAddrs)
	for _, addrString := range nsqdTCPAddrs {
		err := muTodb.ConnectToNSQ(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	for _, addrString := range lookupdHTTPAddrs {
		log.Printf("lookupd addr %s", addrString)
		err := muTodb.ConnectToLookupd(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	return
}

func runNetFlowClient(nfh *net_flow.NetFlowHandler) (nfTodb *nsq.Reader, err error) {
	nfTodb, err = nsq.NewReader("net_flow", "multidb")
	if err != nil {
		log.Fatalf(err.Error())
	}
	nfTodb.SetMaxInFlight(*maxInFlight)
	nfTodb.SetMaxBackoffDuration(*maxBackoffDuration)
	nfTodb.VerboseLogging = *verbose
	nfTodb.AddHandler(nfh)

	fmt.Println(nsqdTCPAddrs)
	for _, addrString := range nsqdTCPAddrs {
		err := nfTodb.ConnectToNSQ(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	for _, addrString := range lookupdHTTPAddrs {
		log.Printf("lookupd addr %s", addrString)
		err := nfTodb.ConnectToLookupd(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	return
}

func runHeartBeatClient(hbh *heartbeat.HeartBeatHandler) (hbTodb *nsq.Reader, err error) {
	hbTodb, err = nsq.NewReader("heartbeat", "multidb")
	if err != nil {
		log.Fatalf(err.Error())
	}
	hbTodb.SetMaxInFlight(*maxInFlight)
	hbTodb.SetMaxBackoffDuration(*maxBackoffDuration)
	hbTodb.VerboseLogging = *verbose
	hbTodb.AddHandler(hbh)

	fmt.Println(nsqdTCPAddrs)
	for _, addrString := range nsqdTCPAddrs {
		err := hbTodb.ConnectToNSQ(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	for _, addrString := range lookupdHTTPAddrs {
		log.Printf("lookupd addr %s", addrString)
		err := hbTodb.ConnectToLookupd(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	return
}

func runAccessibilityToDBClient(ath *accessibility.AccessibilityToDBHandler) (aTodb *nsq.Reader, err error) {
	aTodb, err = nsq.NewReader("accessibility", "multidb")
	if err != nil {
		log.Fatalf(err.Error())
	}
	aTodb.SetMaxInFlight(*maxInFlight)
	aTodb.SetMaxBackoffDuration(*maxBackoffDuration)
	aTodb.VerboseLogging = *verbose
	aTodb.AddHandler(ath)

	fmt.Println(nsqdTCPAddrs)
	for _, addrString := range nsqdTCPAddrs {
		err := aTodb.ConnectToNSQ(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	for _, addrString := range lookupdHTTPAddrs {
		log.Printf("lookupd addr %s", addrString)
		err := aTodb.ConnectToLookupd(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	return
}

func runAccessibilityCheckClient(ach *accessibility.AccessibilityCheckHandler) (aCheck *nsq.Reader, err error) {
	aCheck, err = nsq.NewReader("accessibility", "check_exception")
	if err != nil {
		log.Fatalf(err.Error())
	}
	aCheck.SetMaxInFlight(*maxInFlight)
	aCheck.SetMaxBackoffDuration(*maxBackoffDuration)
	aCheck.VerboseLogging = *verbose
	aCheck.AddHandler(ach)

	fmt.Println(nsqdTCPAddrs)
	for _, addrString := range nsqdTCPAddrs {
		err := aCheck.ConnectToNSQ(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	for _, addrString := range lookupdHTTPAddrs {
		log.Printf("lookupd addr %s", addrString)
		err := aCheck.ConnectToLookupd(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	return
}

func runRegisterToDBClient(rh *register.RegisterToDBHandler) (registerTodb *nsq.Reader, err error) {
	registerTodb, err = nsq.NewReader("register", "multidb")
	if err != nil {
		log.Fatalf(err.Error())
	}
	registerTodb.SetMaxInFlight(*maxInFlight)
	registerTodb.SetMaxBackoffDuration(*maxBackoffDuration)
	registerTodb.VerboseLogging = *verbose
	registerTodb.AddHandler(rh)

	fmt.Println(nsqdTCPAddrs)
	for _, addrString := range nsqdTCPAddrs {
		err := registerTodb.ConnectToNSQ(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	for _, addrString := range lookupdHTTPAddrs {
		log.Printf("lookupd addr %s", addrString)
		err := registerTodb.ConnectToLookupd(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	return
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("nsq_to_http v%s\n", util.BINARY_VERSION)
		return
	}

	if *maxInFlight <= 0 {
		log.Fatalf("--max-in-flight must be > 0")
	}

	if len(nsqdTCPAddrs) == 0 && len(lookupdHTTPAddrs) == 0 {
		log.Fatalf("--nsqd-tcp-address or --lookupd-http-address required")
	}
	if len(nsqdTCPAddrs) > 0 && len(lookupdHTTPAddrs) > 0 {
		log.Fatalf("use --nsqd-tcp-address or --lookupd-http-address not both")
	}

	termChan = make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)


	date := time.Now().Format("20060102")
	// 初始化各种指标的处理类
	cpu_usage_db_link := util.NewDbLink(*dbPath, date)
	cpuUsageHandler, err := cpu_usage.NewCPUUsageHandler(cpu_usage_db_link)
	if err != nil {
		fmt.Println(err)
	}

	mem_usage_db_link := util.NewDbLink(*dbPath, date)
	memUsageHandler, err := mem_usage.NewMemUsageHandler(mem_usage_db_link)
	if err != nil {
		fmt.Println(err)
	}

	net_flow_db_link := util.NewDbLink(*dbPath, date)
	netFlowHandler, err := net_flow.NewNetFlowHandler(net_flow_db_link)
	if err != nil {
		fmt.Println(err)
	}

	heartbeat_db_link := util.NewDbLink(*dbPath, date)
	heartBeatHandler, err := heartbeat.NewHeartBeatHandler(heartbeat_db_link)
	if err != nil {
		fmt.Println(err)
	}

	accessibility_db_link := util.NewDbLink(*dbPath, date)
	accessibilityToDBHandler, err := accessibility.NewAccessibilityToDBHandler(accessibility_db_link)
	if err != nil {
		fmt.Println(err)
	}

	// 可达性异常检测处理类，无需读写DB
	accessibilityCheckHandler, err := accessibility.NewAccessibilityCheckHandler()
	if err != nil {
		fmt.Println(err)
	}

	register_db_link, err := sql.Open("sqlite3", "../db/register.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	util.CreateTable("register", register_db_link)
	registerToDBHandler, err := register.NewRegisterToDBHandler(register_db_link)
	if err != nil {
		fmt.Println(err)
	}

	/*
	心跳数据定期检测，根据检测的结果修改register数据表中机器（正常运行、不正常运行两类）的当前状态
	检测条件：3分钟内是否收到心跳数据
	*/
	heartBeatHandler.CheckPeriodically(register_db_link)

	// 注册各种指标的处理类，各自连接到NSQ的某个channel
	cuTodb, err := runCpuUsageClient(cpuUsageHandler)
	if err != nil {
		fmt.Println(err)
	}

	muTodb, err := runMemUsageClient(memUsageHandler)
	if err != nil {
		fmt.Println(err)
	}

	nfTodb, err := runNetFlowClient(netFlowHandler)
	if err != nil {
		fmt.Println(err)
	}

	hbTodb, err := runHeartBeatClient(heartBeatHandler)
	if err != nil {
		fmt.Println(err)
	}


	aTodb, err := runAccessibilityToDBClient(accessibilityToDBHandler)
	if err != nil {
		fmt.Println(err)
	}

	aCheck, err := runAccessibilityCheckClient(accessibilityCheckHandler)
	if err != nil {
		fmt.Println(err)
	}

	rTodb, err := runRegisterToDBClient(registerToDBHandler)
	if err != nil {
		fmt.Println(err)
	}

	for {
		select {
		case <-muTodb.ExitChan:
			return
		case <-cuTodb.ExitChan:
			return
		case <-nfTodb.ExitChan:
			return
		case <-hbTodb.ExitChan:
			return
		case <-aTodb.ExitChan:
			return
		case <-aCheck.ExitChan:
			return
		case <-rTodb.ExitChan:
			return

		case <-termChan:
			cuTodb.Stop()
			muTodb.Stop()
			nfTodb.Stop()
			hbTodb.Stop()
			aTodb.Stop()
			aCheck.Stop()
			rTodb.Stop()
		}
	}
}
