package main

/* Error handling functions for ABA
   @author Floris Meester floris.meester@gmail.com */


import (
	"fmt"
	"log"
)


// Print error to stdout and syslog
func printerr(err error){

	if err != nil {
		fmt.Println(err)
		log.Print(err)
	}
}

// In case of a fatal error log to syslog and bailout
func fatalerr(err error){

        if err != nil {
                log.Print(err)
		panic(err)
        }
}

// If debug is set produce more output 
func debugerr(comment string, data interface{}, configuration Configuration){

	if configuration.Debug {
                fmt.Println(comment, data)
                log.Print(comment, data)
	}
}

func stdoutlog(data interface{}, configuration Configuration){

	if configuration.Stdout{
		fmt.Println(data)
	}
}
