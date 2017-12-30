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
//	"strconv"
//	"encoding/binary"
)




var lkms []string

func lkmreader(configuration Configuration){

	checklkm(true, configuration)

	for {
		checklkm(false, configuration)
        	time.Sleep(time.Duration(configuration.Sleep)  * time.Second)
        }

}

func checklkm(initial bool, configuration Configuration){
	var res string
	var lines []string
	lkmfile, err := os.Open("/proc/modules")
	if err != nil {
		fatalerr(err)
	}

	defer lkmfile.Close()
	scanner := bufio.NewScanner(lkmfile)
    	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	
	}
//	i := 0
//	lines = append(lines[:i], lines[i+1:]...)
		
	for _,line := range lines {
		l := strings.Split(line, " ")
		//if strings.TrimSpace(l[1]) != "0100007F" && strings.TrimSpace(remote[1]) == "00000000" {
		res = strings.TrimSpace(l[0])
		if initial == true {
			lkms = append(lkms,  fmt.Sprintf("%v",res))
		} else {
       			if lkmcompare(fmt.Sprintf("%v",res), lkms) {
       				stdoutlog(fmt.Sprintf("%v",res) + " exists", configuration)
       			}else {
				stdoutlog("New LKM detected: " + fmt.Sprintf("  %v",res), configuration)
                      		var hashdata Hashdata
                       		hashdata.Hostname = configuration.Hostid
                       		hashdata.Filename =  fmt.Sprintf("  %v",res)
                       		hashdata.Hash = ""
                       		hashdata.Index = false
                       		hashdata.Event  = "lkm"
                       		hashdata.Notify = configuration.Notify
                       		hashjson, err :=  json.Marshal(hashdata)
                       		printerr(err)
                       		// Send data on the channel
                       		producerchannel <- hashjson
			}
			//	}	 
		}
	//	}
			
	}
	debugerr("Loadable kernel modules: ", lkms, configuration)
}


func lkmcompare(res string, lkms []string) bool {
    for _, b := range lkms {
        if b == res {
            return true
        }
    }
    return false
}
