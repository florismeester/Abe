package main

/* Logfile reader for Abe
   @author Floris Meester floris.meester@gmail.com */

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"time"
	"encoding/json"
	//"encoding/hex"
	"strconv"
//	"encoding/binary"
)




var tcpports []string

func tcpportreader(configuration Configuration){

	checktcpport(true, configuration)

	for {
		checktcpport(false, configuration)
        	time.Sleep(time.Duration(configuration.Sleep)  * time.Second)
        }

}

func checktcpport(initial bool, configuration Configuration){
	var res uint64
	var lines []string
	tcpfile, err := os.Open("/proc/net/tcp")
	if err != nil {
		fatalerr(err)
	}

	defer tcpfile.Close()
	scanner := bufio.NewScanner(tcpfile)
    	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	
	}
	i := 0
	lines = append(lines[:i], lines[i+1:]...)
		
	for _,line := range lines {
		if  !strings.HasPrefix(strings.TrimSpace(line), "sl"){
			l := strings.Split(line, ":")
			remote := strings.Split(l[2], " ")
			if strings.TrimSpace(l[1]) != "0100007F" && strings.TrimSpace(remote[1]) == "00000000" {
				res, _= strconv.ParseUint(remote[0],16, 32)
				if initial == true {
					tcpports = append(tcpports,  fmt.Sprintf("%v",res))
				} else {
        				if tcpportcompare(fmt.Sprintf("%v",res), tcpports) {
            					stdoutlog(fmt.Sprintf("%v",res) + " exists", configuration)
        				}else {
						stdoutlog("New port detected: " + fmt.Sprintf("  %v",res), configuration)
                                		var hashdata Hashdata
                                		hashdata.Hostname = configuration.Hostid
                                		hashdata.Filename =  fmt.Sprintf("  %v",res)
                                		hashdata.Hash = ""
                                		hashdata.Index = false
                                		hashdata.Event  = "tcpport"
                                		hashdata.Notify = configuration.Notify
                                		hashjson, err :=  json.Marshal(hashdata)
                                		printerr(err)
                                		// Send data on the channel
                                		producerchannel <- hashjson
					}
				}	 
			}
		}
			
	}
	debugerr("TCP ports: ", tcpports, configuration)
}


func tcpportcompare(res string, tcpports []string) bool {
    for _, b := range tcpports {
        if b == res {
            return true
        }
    }
    return false
}
