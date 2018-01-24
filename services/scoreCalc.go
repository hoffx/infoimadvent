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
	storage.MemoryStorage
	UStorer *iiastorage.UserStorer
	QStorer *iiastorage.QuestStorer
	RStorer *iiastorage.RelationStorer
}

// TODO: find better solution for error handling

func NewDBStorage(uStorer *iiastorage.UserStorer, qStorer *iiastorage.QuestStorer, rStorer *iiastorage.RelationStorer) *DBStorage {
	return &DBStorage{*storage.NewMemoryStorage(), uStorer, qStorer, rStorer}
}

func (s *DBStorage) SetupRoutines() {
	scheduler := scheduler.New(s)
	loc := time.Now().Location()

	_, err := scheduler.RunAt(time.Date(time.Now().Year(), config.Config.Server.ResetMonth, 1, 0, 0, 0, 0, loc), s.SetupYearRoutine, scheduler, time.Now().Year(), loc)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *DBStorage) SetupYearRoutine(scheduler scheduler.Scheduler, year int, loc *time.Location) {
	err := iiastorage.ResetUsers(s.UStorer, s.RStorer)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: implement reward service

	s.SetupCalcRoutine(scheduler, year, loc)

	s.SetupYearRoutine(scheduler, year+1, loc)
}

func (s *DBStorage) SetupCalcRoutine(scheduler scheduler.Scheduler, year int, loc *time.Location) {
	for i := 1; i <= 24; i++ {
		_, err := scheduler.RunAt(time.Date(year, time.January, i, 3, 0, 0, 0, loc), s.CalcScores)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *DBStorage) CalcScores() {
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
		quest, err := s.QStorer.Get(map[string]interface{}{"day": d, "grade": u.Grade})
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
