package main

/* Structs for ABA
   @author Floris Meester floris.meester@gmail.com */

import (
        "github.com/jinzhu/gorm"
//	"time"
)




type Configuration struct {

	Sysloghost string
	Tablename string
	Hostid string
	Syslogproto string
	Syslogport string
	Stdout bool
	Localonly bool
	Paths []string
	Ignoresuffix []string
	Ignoreprefix []string
	Public string
	Secret string
	Clientkeys []string
	Messageport string
	Serveraddr string
	Domain string
	Network string
	Dbname string
	Dbuser string
	Dbpass string
	Debug bool
	Destructive bool
	Notify string
	Smtp string
	Smtpport int
	Server bool
}


type Hashdata struct {

	gorm.Model
	Hostname string `validate:"nonzero"`
	Filename string `validate:"nonzero"`
	Hash	string 
	Event	string  
	Index	bool  
	Notify  string
}

func (c Hashdata) TableName() string {
    return tablename
}


type Mail struct {
	Sender string
	To    string
	Subject  string
	Body     string
}
