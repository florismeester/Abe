package main

/* Database actions for ABA
   @author Floris Meester floris.meester@gmail.com */

import (
	"fmt"
        _ "github.com/lib/pq"
	"github.com/jinzhu/gorm"
	json "encoding/json"
	//"encoding/hex"
	

)




// Get a database connection
func getdbconn(configuration Configuration) (db *gorm.DB) {

        db, err := gorm.Open("postgres", "user=" + configuration.Dbuser +  " password=" + configuration.Dbpass +  " DB.name=" +
	configuration.Dbname +  " sslmode=disable")
	fatalerr(err)
	
	// Run automigrate before returning the connection in case the model has changed
        // fixme db.Debug().AutoMigrate(&Hashdata{})
	db.AutoMigrate(&Hashdata{})
	
	return db
}


func comparemessage(configuration Configuration, hashdata []byte, db  *gorm.DB) (Hashdata, error){

        // Create an instance of a Hashdata struct for unmarshalling data
        p := Hashdata{}
	var result Hashdata
        err := json.Unmarshal(hashdata, &p)
        printerr(err)
	err = db.Where("hash = ? AND filename = ? AND hostname = ?", p.Hash, p.Filename, p.Hostname).Last(&result).Error
	printerr(err)
	
	return result, err

}

// If index is true get file on filename and update if needed
func updater(configuration Configuration, p Hashdata, db  *gorm.DB){
	
	var result Hashdata		
	err := db.Where("filename = ? AND hostname = ?", p.Filename, p.Hostname).Last(&result).Error
	if err != nil {
		db.Save(&p)
	}else{
		result.Hash = p.Hash
		db.Save(&result)
	}
}


func create(configuration Configuration, p Hashdata, db  *gorm.DB){
	
	var result Hashdata		
	err := db.Where("filename = ? AND hostname = ?", p.Filename, p.Hostname).Last(&result).Error
	if err != nil {
		// record doesn't exist, so save it
		db.Save(&p)
	}else{
		// record exists, so noting
		debugerr("record exists:", result, configuration)
	}
}

// Commit a message
func commitmessage(configuration Configuration, hashdata []byte, db *gorm.DB){

	
	// Create an instance of a Hashdata struct for unmarshalling data
	p := Hashdata{}
	err := json.Unmarshal(hashdata, &p)
	printerr(err)
	
	// if index is true update the record or create if it doesn't exist
	if p.Index == true {
		updater(configuration, p, db)	

	} else {
		if configuration.Destructive {	

			// update the record if exist otherwise create
			updater(configuration, p, db)

		} else {

			// only save if not exists
			create(configuration, p, db)

		}
	}	
	if configuration.Debug == true {
		var result Hashdata
		db.Last(&result)
		fmt.Println("Last db commit:", result)
	}
	//isit, err := comparemessage(configuration, hashdata, db)
	//fmt.Println("ISIT:", isit, err)
}

