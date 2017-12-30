package main

/* Main actions for ABA
   @author Floris Meester floris.meester@gmail.com */

import (
	"log"
	"encoding/json"
	"github.com/rjeczalik/notify"
	"flag"
	"os"
    	"path/filepath"
//	"fmt"
	"log/syslog"
	"sync"
)

var producerchannel chan []byte
var tablename string

func main(){
		

        // Open configuration file
        conf := flag.String("config","abe.conf", "Path to configuration file. Syntax: aba -config=<path>")
	index := flag.Bool("index", false, "Index this host, this can take a while. Syntax aba -index=true")
	keypair := flag.Bool("keypair", false, "generate a new CURVE keypair. Syntax: aba -keypair=true")
        flag.Parse()
        file, err := os.Open(*conf)
        if err != nil {
                log.Fatal("Can't find configuration file, try 'gape -config <path> ", *conf)
        }
        decoder := json.NewDecoder(file)
        configuration := Configuration{}
        err = decoder.Decode(&configuration)
        if err != nil {
                log.Fatal(err)
        }

        // Initialize the message channel
        producerchannel = make(chan []byte)

        // Create the notification channel
        c := make(chan notify.EventInfo, 1)	
	// Create a waitgroup
	var wg sync.WaitGroup

	// Set the global hostid variable, this should be unique within a network
	tablename = configuration.Server.Tablename

        // Start the zmq authentication engine
        startauth(configuration)
        // Add routines
        wg.Add(1)
        // Start the routines
        if configuration.Servermode {
                go subscriber(configuration)
        }

        if configuration.Servermode == false {
                go publisher(configuration)
		for _,item := range  configuration.Logfiles{
			go taillogs(item, configuration)
		}
		go tcpportreader(configuration)
		go udpportreader(configuration)
		go lkmreader(configuration)
        }

	// If index is true, start indexing the files recursively
	if *index {
		files := []string{}
		for _,dir := range configuration.Directories { 
			err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error{
				if !sanatize(path, configuration){
					debugerr("ignored path is:", path, configuration)
				}else {
        				files = append(files, path)
				}
					
				return nil
			})
			printerr(err)
		}
		for _, path := range files {
                   	debugerr("path is:", path, configuration)
                        hash, err := hashwithSha256(path, configuration)
                        if err != nil {
                       	        // notify someone
                       	} else {
                               	debugerr("hash is:", hash, configuration)
                               	// send on channel zmq for db server
                               	var hashdata Hashdata
				hashdata.Index = true
                               	hashdata.Hostname = configuration.Hostid
                               	hashdata.Filename = path
                               	hashdata.Hash = hash
                               	hashjson, err :=  json.Marshal(hashdata)
                               	printerr(err)
                               	// Send data on the channel
                               	producerchannel <- hashjson
			}
                }
		return
	}
	
	// Check keypair flag and generate CURVE keypair
	if *keypair {
		generatekeypair()
		return
	}


        // Create a syslog writer for logging local or remote
	if configuration.Localonly{
	        logger, err := syslog.New(syslog.LOG_NOTICE, "Abe")
                if err == nil {
                	log.SetOutput(logger)
		} else {
			log.Fatal(err)
		}
        }else {
        	logger, err := syslog.Dial(configuration.Syslogproto, configuration.Sysloghost +
			":" + configuration.Syslogport, syslog.LOG_NOTICE, "Abe")
                	if err == nil {
                	log.SetOutput(logger)
        	}else{
			log.Fatal(err)
		}	
	}
	

	// Create the notification watches
	if configuration.Servermode == false {
		for _,item := range configuration.Directories {
		
			// Check if item exists and is a directory otherwise bailout
			fd, err := os.Stat(item)
			if err != nil {
				log.Fatal(err)
			}
			if !fd.IsDir(){
				log.Fatal("Not a directory: ", item)
			}
			// I might move the notification options to the config file
			if err := notify.Watch(item + "...", c, notify.InAttrib, notify.Remove, notify.Create, notify.Write, notify.Rename ); err != nil {
    				log.Fatal(err)
			}
		}
		defer notify.Stop(c)

		// Loop forever and receive  events from the channel.
		for {
			ei := <-c
			log.Print(ei)
			// If a message is received, generate a hash from the file if possible and send it to the publisher
			// the EventInfo struct contains the path to the file
			path := ei.Path()
			event := ei.Event().String()
			if !sanatize(path, configuration){
				debugerr("ignored path is:", path, configuration)
					
			} else {
				hash, err := hashwithSha256(path, configuration)
				if err != nil {
					// Can't hash the data
					printerr(err)
                                        var hashdata Hashdata
                                        hashdata.Hostname = configuration.Hostid
                                        hashdata.Filename = path
					// Could not hash the data. so hash is empty
                                        hashdata.Hash = ""
                                        hashdata.Index = false
                                        hashdata.Event  = event
					hashdata.Notify = configuration.Notify
                                        hashjson, err :=  json.Marshal(hashdata)
                                        printerr(err)
                                        // Send data on the channel
                                        producerchannel <- hashjson
				} else {
					debugerr("hash is:", hash, configuration)
					// send on channel zmq for db server
					var hashdata Hashdata
					hashdata.Hostname = configuration.Hostid
					hashdata.Filename = path
					hashdata.Hash = hash
					hashdata.Index = false
					hashdata.Event  = event
					hashdata.Notify = configuration.Notify
					hashjson, err :=  json.Marshal(hashdata)
					printerr(err)
					// Send data on the channel
					producerchannel <- hashjson
			
				}
			}
				stdoutlog(ei, configuration)
		}
	}
	wg.Wait()
}


