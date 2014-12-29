package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jonaz/mdns"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/logger"
	"github.com/ninjasphere/go-ninja/support"
	"github.com/ninjasphere/kodi_jsonrpc"
)

var info = ninja.LoadModuleInfo("./package.json")
var log = logger.GetLogger(info.Name)

type Driver struct {
	support.DriverSupport
}

func NewDriver() (*Driver, error) {

	driver := &Driver{}

	err := driver.Init(info)
	if err != nil {
		log.Fatalf("Failed to initialize driver: %s", err)
	}

	err = driver.Export(driver)
	if err != nil {
		log.Fatalf("Failed to export driver: %s", err)
	}

	return driver, nil
}

func (d *Driver) Start(_ interface{}) error {
	log.Infof("Driver Starting")

	castService := "_xbmc-jsonrpc._tcp"

	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {

			log.Debugf("Found mdns service: %v", entry)

			if !strings.Contains(entry.Name, castService) {
				return
			}

			log.Infof("Got new kodi: %v", entry)

			kodi, err := kodi_jsonrpc.New(fmt.Sprintf("%s:%d", entry.Addr, entry.Port), 15)

			if err != nil {
				log.Fatalf("Failed to connect to %s: %s", entry.Addr, err)
			}

			NewMediaPlayer(d, d.Conn, entry.Name, kodi)
		}
	}()

	go func() {
		// Start the lookup
		params := mdns.DefaultParams(castService)
		params.Entries = entriesCh
		params.Timeout = time.Second * 20
		mdns.Query(params)
	}()

	//spew.Dump(goupnp.DiscoverDevices("uuid:251ddbdc-effa-be48-4b16-fd46ad72bd73"))

	return nil
}
