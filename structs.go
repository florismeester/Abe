package main

/* Structs for ABA
   @author Floris Meester floris.meester@gmail.com */

import (
        "github.com/jinzhu/gorm"
//	"time"
)



type Files struct {

	Path string
	Filter []string
}

type ServerType struct {

	Tablename string
	Dbname string
	Dbuser string
	Dbpass string
	Clientkeys []string
	Sender string
	Smtp string
	Smtpport int
	Subject string
}	

type Configuration struct {

	Sysloghost string
	Hostid string
	Syslogproto string
	Syslogport string
	Stdout bool
	Localonly bool
	Public string
	Secret string
	Messageport string
	Serveraddr string
	Domain string
	Network string
	Debug bool
	Destructive bool
	Notify string
	Servermode bool
	Logfiles []Files	
	Directories []string
	Suffixes []string
	Filters []string
	Server ServerType
	
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
