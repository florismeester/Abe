package main

/* Tools function for ABA
   @author Floris Meester floris.meester@gmail.com */

import (
	"strings"
//	"fmt"
	"gopkg.in/gomail.v2"
	"crypto/tls"
	"regexp"
)

// Test for suffixes or filters to ignore
func sanatize(filepath string, configuration Configuration) (bool) {
	
	flag := true
	
	// If it matches put the flag to false, since it's not allowed  to use
	// fixme: this should be replaced by a regex solution
	for _,item := range  configuration.Suffixes {
		if strings.HasSuffix(filepath, item){ 
			flag = false
			debugerr("suffix matches", item, configuration)
		}
	}
        for _,item := range  configuration.Filters {
		 matched, err := regexp.MatchString(".*" + item + ".*", filepath)
		 if err != nil{
			printerr(err)
		 }
		 if matched == true {
			flag = false
			debugerr("filter matches", item, configuration)
		 }
	}
	return flag
	
}


// Mail stuff
func sendemail(configuration Configuration, hashdata Hashdata) {

	// create the message
	var mail Mail
	var event string
	if hashdata.Event == "notify.Remove" {
		event = "Remove"
	} else if  hashdata.Event == "notify.InAttrib" {
		event = "Atributes changed"
	} else if hashdata.Event == "notify.Create" {
		event  = "Create"
	} else if hashdata.Event == "notify.Rename" {
		event = "Rename"
	} else if hashdata.Event == "notify.Write" {
		event = "Write"
	} else if hashdata.Event == "syslog" {
		event = "Log"
	}
       
	
	mail.Sender = configuration.Server.Sender
	mail.To =  hashdata.Notify
	mail.Subject = configuration.Server.Subject
	mail.Body = "This the Abe system to inform you about the following incident:\n" +
	"The system noticed a " + event + " event.\n" + "Output: " + hashdata.Filename

	message := gomail.NewMessage()
	message.SetHeader("From", mail.Sender)
	address := strings.Split(mail.To, ",")
	addr := make([]string, len(address))
	for i, to := range address {
		addr[i] = message.FormatAddress(to, "")
	} 
	message.SetHeader("To", addr...)
	message.SetHeader("Subject", mail.Subject)
	message.SetBody("text/plain", mail.Body)
	debugerr("Email: ", message, configuration)
	dialer := gomail.Dialer{ Host: configuration.Server.Smtp, Port: configuration.Server.Smtpport, SSL: false }
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := dialer.DialAndSend(message); err != nil {
    		printerr(err)
	}

}








