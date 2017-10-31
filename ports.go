package main

/* Logfile reader for Abe
   @author Floris Meester floris.meester@gmail.com */

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	//"encoding/hex"
	"strconv"
//	"encoding/binary"
)




func portreader(configuration Configuration){

	tcpfile, err := os.Open("/proc/net/tcp")
	if err != nil {
		fatalerr(err)
	}
	
	udpfile, err := os.Open("/proc/net/udp")
	if err != nil {
		fatalerr(err)
	}
	
	defer tcpfile.Close()	
	defer udpfile.Close()	

	scanner := bufio.NewScanner(tcpfile)
	var lines []string
    	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		
	}
	i := 0
	lines = append(lines[:i], lines[i+1:]...)
	fmt.Println(lines)
	for _,line := range lines {
		l := strings.Split(line, ":")
		remote := strings.Split(l[2], " ")
		if strings.TrimSpace(l[1]) != "0100007F" && strings.TrimSpace(remote[1]) == "00000000" {
			fmt.Println(strconv.ParseUint(remote[0],16, 32))
		}
	}
		 

}



/*
                for _,item := range  item.Filter {
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
*/
