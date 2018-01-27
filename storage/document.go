package storage

import (
	"os"

	"github.com/go-xorm/core"
	"github.com/hoffx/infoimadvent/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type DocumentStorer struct {
	Storer
	Complete bool
}

type Document struct {
	Path     string
	Grade    int
	Day      int
	Solution int
	Type     int
}

func NewDocumentStorer(name, user, password string, doLog bool) (DocumentStorer, error) {
	db, err := xorm.NewEngine("mysql", user+":"+password+"@/"+name+"?charset=utf8")
	if err != nil {
		return DocumentStorer{}, err
	}
	if doLog {
		db.ShowSQL(true)
		db.Logger().SetLevel(core.LOG_DEBUG)
		f, err := os.Create("sql.log")
		if err != nil {
			return DocumentStorer{}, err
		}
		db.SetLogger(xorm.NewSimpleLogger(f))
	}
	err = db.Sync(new(Document))
	if err != nil {
		return DocumentStorer{}, err
	}
	err = db.CreateTables(&Document{})
	if err != nil {
		return DocumentStorer{}, err
	}

	qs := DocumentStorer{Storer{db, true}, false}
	qs.Complete = qs.isComplete()

	return qs, nil
}

func (s *DocumentStorer) ResetDB() error {
	s.Complete = false
	err := s.db.DropTables(Document{})
	if err != nil {
		return err
	}
	err = s.db.CreateTables(Document{})
	return err
}

func (s *DocumentStorer) Create(document Document) error {
	_, err := s.db.Insert(document)
	if err != nil {
		return err
	}

	s.Complete = s.isComplete()

	return nil
}

func (s *DocumentStorer) Put(document Document) error {
	oldDocument, err := s.Get(map[string]interface{}{"grade": document.Grade, "day": document.Day})
	if err != nil {
		return err
	}
	i, err := s.db.Delete(&oldDocument)
	if i == 0 {
		err = ErrNoEffect
	}
	if err != nil {
		return err
	}
	return s.Create(document)
}

func (s *DocumentStorer) Get(keys map[string]interface{}) (Document, error) {
	if len(keys) == 0 {
		return Document{}, ErrNoKey
	}
	query, values := buildQuery(keys)
	document := Document{}
	_, err := s.db.Table("document").Where(query, values...).Get(&document)
	return document, err
}

func (s *DocumentStorer) GetAll(keys map[string]interface{}) (documents []Document, err error) {
	query, values := buildQuery(keys)
	err = s.db.Table("document").Where(query, values...).Find(&documents)
	return
}

func (s *DocumentStorer) isComplete() bool {
	for day := 1; day <= 24; day++ {
		for grade := config.Config.Grades.Min; grade <= config.Config.Grades.Max; grade++ {
			q, err := s.Get(map[string]interface{}{"day": day, "grade": grade})
			if q.Path == "" || err != nil {
				return false
			}
		}
	}
	return true
}

func (s *DocumentStorer) Delete(keys map[string]interface{}) error {
	if len(keys) == 0 {
		return ErrNoKey
	}
	query, values := buildQuery(keys)
	_, err := s.db.Table("document").Where(query, values...).Delete(Document{})
	return err
}
