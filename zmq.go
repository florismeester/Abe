package main

/* ZMQ api for ABA
   @author Floris Meester floris.meester@gmail.com */

import (
        "fmt"
        zmq "github.com/pebbe/zmq4"
	"gopkg.in/validator.v2"
	"encoding/json"
)


// Generate a Curve keypair
func generatekeypair(){

	client_public, client_secret, err := zmq.NewCurveKeypair()
	printerr(err)
	fmt.Println("Generated the following CURVE pair:")
	fmt.Println("Secret:", client_secret)
        fmt.Println("Public:", client_public)
}

// Start the authentication engine
func startauth(configuration Configuration){
	
        zmq.AuthSetVerbose(true)
        zmq.AuthStart()
        zmq.AuthAllow(configuration.Domain, configuration.Network)
}

// Publisher code when in client mode (server=false)
func  publisher(configuration Configuration){
	
        publisher, _ := zmq.NewSocket(zmq.PUB)
        publisher.ClientAuthCurve(configuration.Public, configuration.Public, configuration.Secret)
        defer publisher.Close()
        err := publisher.Connect("tcp://"+ configuration.Serveraddr + ":" + configuration.Messageport)
	printerr(err)
        for{
                val := <- producerchannel
		debugerr("publisher: ", val, configuration)
                _, err := publisher.SendBytes(val, 0)
		printerr(err)
        }
}

// Subscriber code that binds and listens for json messages (server=true)
func subscriber (configuration Configuration){

	// fixme: Add the public keys of clients, this should be fixed to a better solution
	for _,pub := range configuration.Server.Clientkeys{
		zmq.AuthCurveAdd(configuration.Domain, pub)
	}
	
	// Create a new subscriber socket
        subscriber, _ := zmq.NewSocket(zmq.SUB)

	// for now a dummy configuartion option
	subscriber.ServerAuthCurve(configuration.Domain, configuration.Secret)

	// Bind and listen to a socker
        err := subscriber.Bind("tcp://*" + ":" + configuration.Messageport)
	printerr(err)

	// Set a general subscription
        subscriber.SetSubscribe("")

	// Should never get there
        defer subscriber.Close()
	
        // Create a database connection
        db := getdbconn(configuration)
	
	// Create an instance of a message struct
	var hashdata Hashdata

        for{
		// Start receiving 
                mess, err := subscriber.RecvBytes(0)
                printerr(err)
		
                debugerr("receiver:", string(mess), configuration)
		
		// Unmarshall the message in a struct 
        	err = json.Unmarshal(mess, &hashdata)
                if err := validator.Validate(hashdata); err != nil {
                        printerr(err)
                } else if (hashdata.Event == "notify.InAttrib" ||  hashdata.Event == "notify.Remove" || hashdata.Event == "syslog" || hashdata.Hash == "") {
			if configuration.Debug == true {
				fmt.Println("receiver Event", hashdata.Event, hashdata.Filename)
			}
			// just send a notification, no persistence
			fmt.Println("attrib or log:" + hashdata.Filename)
			sendemail(configuration, hashdata)
		}else {
                       	// If validated marshal and commit
                       	hashjson, err := json.Marshal(hashdata)
                       	printerr(err)
			sendemail(configuration, hashdata)
			commitmessage(configuration, hashjson, db)
		}
        }
}
