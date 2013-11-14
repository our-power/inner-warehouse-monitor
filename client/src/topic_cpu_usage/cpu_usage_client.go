package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"time"
	"util"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/bitly/go-nsq"
)

var (
	showVersion        = flag.Bool("version", false, "print version string")
	nsqdTCPAddrs     = util.StringArray{}
	lookupdHTTPAddrs = util.StringArray{}
	topic            = flag.String("topic", "", "nsq topic")
	channel          = flag.String("channel", "nsq_to_file", "nsq channel")
	maxInFlight        = flag.Int("max-in-flight", 200, "max number of messages to allow in flight")
	verbose            = flag.Bool("verbose", false, "enable verbose logging")
	maxBackoffDuration = flag.Duration("max-backoff-duration", 120*time.Second, "the maximum backoff duration")
	dbPath			= flag.String("dbPath", "./", "the path to store db file")
)

func init() {
	flag.Var(&nsqdTCPAddrs, "nsqd-tcp-address", "nsqd TCP address (may be given multiple times)")
	flag.Var(&lookupdHTTPAddrs, "lookupd-http-address", "lookupd HTTP address (may be given multiple times)")
}

type DBHandler struct {
	db sql.DB
}

func (h *DBHandler) HandleMessage(m *nsq.Message)(err error){
	/*
	实现队列消息处理功能
	*/
	fmt.Println(m.ID)
	fmt.Printf("%s\n", m.Body)
	fmt.Println(m.Timestamp)
	fmt.Println(m.Attempts)
	return nil
}

func NewDBHandler(dbDriver string, dbSourceName string)(dbHandler DBHandler, err error){
	notExist := false
	if _, e := os.Stat(dbSourceName); os.IsNotExist(e) {
		notExist = true
	}
	link, err := sql.Open(dbDriver, dbSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer link.Close()
	dbHandler := &DBHandler {
		db: link,
	}
	return dbHandler, err
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("nsq_to_http v%s\n", util.BINARY_VERSION)
		return
	}

	if *topic == "" || *channel == "" {
		log.Fatalf("--topic and --channel are required")
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

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)


	date := time.Now().Format("2006-01-02")
	if !strings.HasSuffix(dbPath, "/") {
		dbPath = dbPath + "/"
	}

	dbHandler, err := NewDBHandler("sqlite3", dbPath + date + ".db")
	if err != nil {
		fmt.Println(err)
		return
	}

	r, err := nsq.NewReader(*topic, *channel)
	if err != nil {
		log.Fatalf(err.Error())
	}
	r.SetMaxInFlight(*maxInFlight)
	r.SetMaxBackoffDuration(*maxBackoffDuration)
	r.VerboseLogging = *verbose
	r.AddHandler(dbHandler)

	for _, addrString := range nsqdTCPAddrs {
		err := r.ConnectToNSQ(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	for _, addrString := range lookupdHTTPAddrs {
		log.Printf("lookupd addr %s", addrString)
		err := r.ConnectToLookupd(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	for {
		select {
		case <-r.ExitChan:
			return
		case <-termChan:
			r.Stop()
		}
	}
}





