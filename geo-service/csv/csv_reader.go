package csv

import (
	"encoding/csv"
	. "findhotel.com/geo-service/config"
	mongodb "findhotel.com/geo-service/db"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

//Importer is an interface for data importing
type Importer interface {
	//Import is a function to import a csv file into database. Before calling this method some consideration should be taken.
	//config.Config attributes must have valid amount (for database url, user and pass)
	Import() error
}

//ImporterImpl is responsible for import csv file into database. this the starting point of the process.
//It needs csv file address (FileAddress) and database url (MongodbUrl) and username (MongodbUsername) and password (MongodbPassword) if required
type ImporterImpl struct {
	client     mongodb.ClientWrapper
	dispatcher Dispatcher
}

//NewImporterImpl creates ImporterImpl, It needs mongodb.ClientWrapper for connecting to MongoDB and a Dispatcher to dispatch
//csv records between workers
func NewImporterImpl(client mongodb.ClientWrapper, dispatcher Dispatcher) ImporterImpl {
	return ImporterImpl{
		client:     client,
		dispatcher: dispatcher,
	}
}

func (r ImporterImpl) Import() error {

	//open csv file
	f, err := os.Open(C.FileAddress)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Warn().Err(err).Str("file", C.FileAddress).Msg("Error when closing file")
		}
	}(f)

	err = r.client.Connect() //Connect to database
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to database")
	}
	defer r.client.Close() //close after finishing importing

	recordChan := make(chan []string)
	finishedCh := make(chan bool)
	r.dispatcher.Dispatch(recordChan, finishedCh) //run a dispatcher and wait for records

	reader := csv.NewReader(f)
	for i := 1; true; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Err(err)
		}
		recordChan <- record //sending records to dispatcher
	}
	finishedCh <- true

	r.dispatcher.Wait()  //wait until workers finish their works
	r.dispatcher.Print() //Print results (time elapsed and number of successful and failed records)
	return nil
}
