package storage

import (
	"os"

	"github.com/go-xorm/core"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type QuestStorer struct {
	Storer
}

type Quest struct {
	Path     string
	Grade    int
	Day      int
	Solution int
}

func NewQuestStorer(name, user, password string, doLog bool) (QuestStorer, error) {
	db, err := xorm.NewEngine("mysql", user+":"+password+"@/"+name+"?charset=utf8")
	if err != nil {
		return QuestStorer{}, err
	}
	if doLog {
		db.ShowSQL(true)
		db.Logger().SetLevel(core.LOG_DEBUG)
		f, err := os.Create("sql.log")
		if err != nil {
			return QuestStorer{}, err
		}
		db.SetLogger(xorm.NewSimpleLogger(f))
	}
	err = db.Sync(new(Quest))
	if err != nil {
		return QuestStorer{}, err
	}
	err = db.CreateTables(&Quest{})
	if err != nil {
		return QuestStorer{}, err
	}

	return QuestStorer{Storer{db, true}}, nil
}

func (s *QuestStorer) ResetDB() error {
	err := s.db.DropTables(Quest{})
	if err != nil {
		return err
	}
	err = s.db.CreateTables(Quest{})
	return err
}

func (s *QuestStorer) Create(quest Quest) error {
	_, err := s.db.Insert(quest)
	return err
}

func (s *QuestStorer) Put(quest Quest) error {
	oldQuest, err := s.Get(map[string]interface{}{"grade": quest.Grade, "day": quest.Day})
	if err != nil {
		return err
	}
	i, err := s.db.Delete(&oldQuest)
	if i == 0 {
		err = ErrNoEffect
	}
	if err != nil {
		return err
	}
	return s.Create(quest)
}

func (s *QuestStorer) Get(keys map[string]interface{}) (Quest, error) {
	if len(keys) == 0 {
		return Quest{}, ErrNoKey
	}
	var query string
	var values []interface{}
	first := true
	for k, v := range keys {
		values = append(values, v)
		if !first {
			query += " AND "
		} else {
			first = false
		}
		query += k + " = ?"
	}
	quest := Quest{}
	_, err := s.db.Table("quest").Where(query, values...).Get(&quest)
	return quest, err
}

func (s *QuestStorer) GetAll() (quests []Quest, err error) {
	err = s.db.Table("quest").Find(&quests)
	return
}
