package main

/* Logfile reader for Abe
   @author Floris Meester floris.meester@gmail.com */

import (
	"github.com/hpcloud/tail"
	"os"
	"encoding/json"
	"regexp"
	"fmt"
	

)




// Get a database connection
func taillogs(logfile string, configuration Configuration){
	config := tail.Config{Follow: true}
	config.Location = &tail.SeekInfo{0, os.SEEK_END} //&tail.SeekInfo{os.SEEK_END}
	t, err := tail.TailFile(logfile , config) //tail.Config{Follow: true, Location: })
	if err != nil {
		fatalerr(err)
	}
	fmt.Println(t)	
	for line := range t.Lines {
		for _,item := range  configuration.Logfilter {
			matched, err := regexp.MatchString(".*" + item + ".*", line.Text)
			if err != nil {
				printerr(err)
			}else if matched == true {
				fmt.Println(item, line.Text)
               			// send on channel zmq 
               			var hashdata Hashdata
               			hashdata.Hostname = configuration.Hostid
               			hashdata.Filename = line.Text
               			hashdata.Hash = ""
          	 		hashdata.Index = false
                		hashdata.Event  = "syslog"
                		hashdata.Notify = configuration.Notify
                		hashjson, err :=  json.Marshal(hashdata)
                		printerr(err)
                		// Send data on the channel
                		producerchannel <- hashjson
			}
		}
	}
}	

