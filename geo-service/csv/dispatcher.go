package csv

import (
	"errors"
	"findhotel.com/geo-service/config"
	mongodb "findhotel.com/geo-service/db"
	"findhotel.com/geo-service/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

type Dispatcher interface {
	Dispatch(csvRecordCh chan []string, finishCh chan bool)
	Print()
	Wait()
}

// DispatcherImpl Dispatch works between workers, caller cab wait on WG until workers finish their works
type DispatcherImpl struct {
	workers []worker
	wg      *sync.WaitGroup
	start   int64
	end     int64
}

//NewDispatcherImpl Create new DispatcherImpl. Workers are created and started to accept new insertJob
func NewDispatcherImpl(geoLocationDao mongodb.GeoLocationDao) *DispatcherImpl {
	workers := make([]worker, config.C.NumberOfWorkers)
	wg := sync.WaitGroup{}
	for i := 0; i < config.C.NumberOfWorkers; i++ {
		workers[i] = worker{
			recordCh:    make(chan []string),
			records:     make([][]string, 0),
			count:       0,
			failedCount: 0,
			finished:    make(chan bool),
			wg:          &wg,
		}
		workers[i].start(geoLocationDao)
		wg.Add(1)
	}
	return &DispatcherImpl{
		workers: workers,
		wg:      &wg,
	}
}

func (r *DispatcherImpl) Wait() {
	r.wg.Wait()
	r.end = time.Now().Unix()
}

//Print just gives some statistics about jobs and workers
func (r *DispatcherImpl) Print() {
	total := 0
	failed := 0
	for _, w := range r.workers {
		total += w.count
		failed += w.failedCount
	}
	log.Info().Msgf("Total time elapsed: %d seconds", r.end-r.start)
	log.Info().Msgf("Total number of successful records: %d", total-failed)
	log.Info().Msgf("Total number of failed records: %d", failed-1) //the first record is always headers and is not failed
}

//Dispatch dispatch tasks to worker(s). csvRecordCh param is a channel to receive records to persist in DB
// and finishCh channel is a channel to tell dispatcher that there is no more records
func (r *DispatcherImpl) Dispatch(csvRecordCh chan []string, finishCh chan bool) {
	r.start = time.Now().Unix()
	go func() {
		workerNumber := 0
		for {
			select {
			case record := <-csvRecordCh:
				r.workers[workerNumber].recordCh <- record                   //send job to worker
				workerNumber = (workerNumber + 1) % config.C.NumberOfWorkers //next worker to send a job
			case <-finishCh:
				for _, worker := range r.workers {
					worker.finished <- true
				}
			}
		}
	}()
}

//worker is a unit to wait
//jobCh is a insertJob Channel, a worker receive jobs from this channel
type worker struct {
	recordCh    chan []string   //Channel for receiving raw csv records
	records     [][]string      //We keep records here until it reach to Config.InsertBulkSize
	count       int             //total number of records this worker processed
	failedCount int             //total number of failed records
	finished    chan bool       //finish channel
	wg          *sync.WaitGroup //This is dispatcher WG, to let dispatcher know worker has finished it's jobs
}

//start start a worker. After starting, it is ready to accept insertJob
func (r *worker) start(geoLocationDao mongodb.GeoLocationDao) {
	go func() {
		for {
			select {
			case record := <-r.recordCh:
				r.records = append(r.records, record)
				r.count++
				if len(r.records) == config.C.InsertBulkSize {
					failCount := insertMany(r.records, geoLocationDao)
					r.failedCount += failCount
					r.records = r.records[:0] //clear records but we do not make it eligible to GC to reuse it later
				}
			case <-r.finished:
				//remaining records
				if len(r.records) > 0 {
					failCount := insertMany(r.records, geoLocationDao)
					r.failedCount += failCount
					r.records = r.records[:0]
				}
				r.wg.Done()
				return
			}
		}
	}()
}

func insertMany(records [][]string, dao mongodb.GeoLocationDao) (failCount int) {
	failCount = 0
	geos := make([]model.GeoLocation, 0)
	for _, record := range records {
		location, err := model.ParseGeoLocation(record)
		if err != nil {
			failCount++
		} else {
			geos = append(geos, location)
		}
	}
	if len(geos) > 0 {
		_, err := dao.InsertMany(geos)
		if err != nil && errors.As(err, &mongo.BulkWriteException{}) {
			failCount += len(err.(mongo.BulkWriteException).WriteErrors)
		}
	}
	return failCount
}
