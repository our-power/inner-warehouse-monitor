package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	//"strings"
	"time"
	"util"
	"github.com/bitly/go-nsq"
	"github.com/influxdb/influxdb-go"
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

	influxdb_host     = flag.String("influxdb_host", "127.0.0.1:8086", "host of influxdb server")
	influxdb_user     = flag.String("influxdb_user", "root", "influxdb username")
	influxdb_passwd   = flag.String("influxdb_passwd", "root", "the passwd of influxdb user")
	influxdb_database = flag.String("influxdb_database", "", "the name of target database")

	termChan chan os.Signal
)

func init() {
	flag.Var(&nsqdTCPAddrs, "nsqd-tcp-address", "nsqd TCP address (may be given multiple times)")
	flag.Var(&lookupdHTTPAddrs, "lookupd-http-address", "lookupd HTTP address (may be given multiple times)")
}

func runCpuUsageClient(cuh *cpu_usage.CPUUsageHandler) (cuTodb *nsq.Reader, err error) {
	cuTodb, err = nsq.NewReader("cpu_usage", "todb")
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
	muTodb, err = nsq.NewReader("mem_usage", "todb")
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
	nfTodb, err = nsq.NewReader("net_flow", "todb")
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
	hbTodb, err = nsq.NewReader("heartbeat", "todb")
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
	aTodb, err = nsq.NewReader("accessibility", "todb")
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
	registerTodb, err = nsq.NewReader("register", "todb")
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

	config := influxdb.ClientConfig {
		Host: *influxdb_host,
		Username: *influxdb_user,
		Password: *influxdb_passwd,
		Database: *influxdb_database
	}

	client, err := influxdb.NewClient(&config)

	// 初始化各种指标的处理类
	cpuUsageHandler, err := cpu_usage.NewCPUUsageHandler(client)
	if err != nil {
		fmt.Println(err)
	}

    memUsageHandler, err := mem_usage.NewMemUsageHandler(client)
    if err != nil {
        fmt.Println(err)
    }

    netFlowHandler, err := net_flow.NewNetFlowHandler(client)
    if err != nil {
        fmt.Println(err)
    }

    heartBeatHandler, err := heartbeat.NewHeartBeatHandler(client)
    if err != nil {
        fmt.Println(err)
    }

    accessibilityToDBHandler, err := accessibility.NewAccessibilityToDBHandler(client)
    if err != nil {
        fmt.Println(err)
    }
    // 可达性异常检测处理类，无需读写DB
    accessibilityCheckHandler, err := accessibility.NewAccessibilityCheckHandler()
    if err != nil {
        fmt.Println(err)
    }

    registerToDBHandler, err := register.NewRegisterToDBHandler(client)
    if err != nil {
        fmt.Println(err)
    }

    /*
    心跳数据定期检测，根据检测的结果修改register数据表中机器（正常运行、不正常运行两类）的当前状态
    检测条件：3分钟内是否收到心跳数据
    */
    //heartBeatHandler.CheckPeriodically(register_db_link)

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
