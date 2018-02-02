package services

import (
	"log"
	"time"

	"github.com/hoffx/infoimadvent/config"
	iiastorage "github.com/hoffx/infoimadvent/storage"
	"github.com/rakanalh/scheduler"
	"github.com/rakanalh/scheduler/storage"
)

type DBStorage struct {
	storage.Sqlite3Storage
	UStorer *iiastorage.UserStorer
	DStorer *iiastorage.DocumentStorer
	RStorer *iiastorage.RelationStorer
}

// TODO: find better solution for error handling

func NewDBStorage(uStorer *iiastorage.UserStorer, dStorer *iiastorage.DocumentStorer, rStorer *iiastorage.RelationStorer) *DBStorage {
	storage := storage.NewSqlite3Storage(storage.Sqlite3Config{config.Config.Scheduler.StoragePath})
	if err := storage.Connect(); err != nil {
		log.Fatal(err)
	}

	if err := storage.Initialize(); err != nil {
		log.Fatal(err)
	}
	return &DBStorage{storage, uStorer, dStorer, rStorer}
}

func (s *DBStorage) SetupRoutines() {
	scheduler := scheduler.New(s)
	loc := time.Now().Location()

	// TODO: remove this test

	scheduler.RunAfter(5*time.Second, s.test)

	_, err := scheduler.RunAt(time.Date(time.Now().Year(), config.Config.Server.ResetMonth, 1, 0, 0, 0, 0, loc), s.SetupYearRoutine, scheduler, time.Now().Year(), loc)
	if err != nil {
		log.Fatal(err)
	}
	scheduler.
}

func (s *DBStorage) SetupYearRoutine(scheduler scheduler.Scheduler, year int, loc *time.Location) {
	err := iiastorage.ResetUsers(s.UStorer, s.RStorer)
	if err != nil {
		log.Fatal(err)
	}
	err = iiastorage.ResetDocuments(s.DStorer, true)
	if err != nil {
		log.Fatal(err)
	}

	s.SetupCalcRoutine(scheduler, year, loc)

	s.SetupYearRoutine(scheduler, year+1, loc)
}

func (s *DBStorage) SetupCalcRoutine(scheduler scheduler.Scheduler, year int, loc *time.Location) {
	// shift calculation day by one because it is done at 3am
	for i := 2; i <= 25; i++ {
		_, err := scheduler.RunAt(time.Date(year, time.January, i, 3, 0, 0, 0, loc), s.calcScores)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *DBStorage) calcScores() {
	users, err := s.UStorer.GetAll(map[string]interface{}{})
	if err != nil {
		log.Println(err)
		return
	}

	_, m, d := time.Now().Date()

	// TODO: change back to december after testing
	if m != time.January {
		return
	}

	for _, u := range users {
		// executed at 3:00 am at server-time => -1 for last day + -1 for slice index-shift
		quest, err := s.DStorer.Get(map[string]interface{}{"day": d, "grade": u.Grade})
		if err != nil {
			log.Println(err)
			continue
		}

		if u.Days[d-2] == quest.Solution {
			u.Score += iiastorage.Right
		} else if u.Days[d-2] == iiastorage.None {
			u.Score += iiastorage.Missing
		} else {
			u.Score += iiastorage.Wrong
		}
	}
}

func (s *DBStorage) test() {
	log.Println("ho")
}
