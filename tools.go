package main

/* Tools function for ABA
   @author Floris Meester floris.meester@gmail.com */

import (
	"strings"
	"fmt"
	"gopkg.in/gomail.v2"
	"crypto/tls"
)

// Test for suffixes or prefixes to ignore
func sanatize(filepath string, configuration Configuration) (bool) {
	
	flag := true
	
	// If it matches put the flag to false, since it's not sane to use
	for _,item := range  configuration.Ignoresuffix {
		if strings.HasSuffix(filepath, item){ 
		flag = false
		fmt.Println("suffix matches", item)
		}
	}
        for _,item := range  configuration.Ignoreprefix {
		if strings.HasPrefix(filepath, item){
		flag = false
		fmt.Println("prefix matches", item)
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
		event = "Removed"
	} else if  hashdata.Event == "notify.InAttrib" {
		event = "Atributes changed"
	} else if hashdata.Event == "notify.Create" {
		event  = "Created"
	}else if hashdata.Event == "notify.Rename" {
		event = "Renamed"
	}else if hashdata.Event == "notify.Write" {
		event = "Written"
	}
       
	
	mail.Sender = "abe@grid6.io"
	mail.To =  hashdata.Notify
	mail.Subject = "Abe filesystem warning"
	mail.Body = "This the Abe system to inform you about the following incident:\n" +
	"The system noticed that " + hashdata.Filename + " had the following event: " +
	event + "."

	message := gomail.NewMessage()
	message.SetHeader("From", mail.Sender)
	address := strings.Split(mail.To, ",")
	addr := make([]string, len(address))
	for i, to := range address {
		addr[i] = message.FormatAddress(to, "")
	} 
	message.SetHeader("To", addr...)
	message.SetHeader("Subject", mail.Subject)
	message.SetBody("text/html", mail.Body)
	//m.Attach("/home/Alex/lolcat.jpg")
	fmt.Println(message)
	dialer := gomail.Dialer{ Host: configuration.Smtp, Port: configuration.Smtpport, SSL: false }
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := dialer.DialAndSend(message); err != nil {
    		printerr(err)
	}

}








