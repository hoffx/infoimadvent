package services

import (
	"log"
	"time"

	iiastorage "github.com/hoffx/infoimadvent/storage"
	"github.com/rakanalh/scheduler/storage"
)

type DBStorage struct {
	storage.MemoryStorage
	UStorer *iiastorage.UserStorer
	QStorer *iiastorage.QuestStorer
}

func NewDBStorage(uStorer *iiastorage.UserStorer, qStorer *iiastorage.QuestStorer) *DBStorage {
	return &DBStorage{*storage.NewMemoryStorage(), uStorer, qStorer}
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
