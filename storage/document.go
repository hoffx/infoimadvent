package storage

import (
	"os"

	"github.com/go-xorm/core"
	"github.com/hoffx/infoimadvent/config"

	// blank import required by xorm
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

// DocumentStorer is a Storer. The extra type is needed
// by macaron for the identification of the storer.
type DocumentStorer struct {
	Storer
	// Complete is true, when all required documents
	// are available
	Complete bool
}

// Document represents a markdown-document.
type Document struct {
	// Path to the .md file
	Path string
	// Grade this document is for
	Grade int
	// Day this document is for
	Day int
	// Solution of the document's question
	Solution int
	// Type of the document (Quest/About/ToS)
	Type int
}

// NewDocumentStorer generates a DocumentStorer. If doLog is true, it
//will create a log file "sql.log"
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

	qs := DocumentStorer{Storer{db}, false}
	qs.Complete = qs.isComplete()

	return qs, nil
}

// ResetDB deletes the document table
func (s *DocumentStorer) ResetDB() error {
	s.Complete = false
	err := s.db.DropTables(Document{})
	if err != nil {
		return err
	}
	err = s.db.CreateTables(Document{})
	return err
}

// Create creates a new entry
func (s *DocumentStorer) Create(document Document) error {
	_, err := s.db.Insert(document)
	if err != nil {
		return err
	}

	s.Complete = s.isComplete()

	return nil
}

// Put modifies a entry. The document is identified by its
// grade, day and type values
func (s *DocumentStorer) Put(document Document) error {
	oldDocument, err := s.Get(map[string]interface{}{"grade": document.Grade, "day": document.Day, "type": document.Type})
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

// Get returns one document defined by the provided keys
func (s *DocumentStorer) Get(keys map[string]interface{}) (Document, error) {
	if len(keys) == 0 {
		return Document{}, ErrNoKey
	}
	query, values := buildQuery(keys)
	document := Document{}
	_, err := s.db.Table("document").Where(query, values...).Get(&document)
	return document, err
}

// GetAll returns a slice of documents defined by the provided keys
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

// Delete deletes all entries defined by the provided keys
func (s *DocumentStorer) Delete(keys map[string]interface{}) error {
	if len(keys) == 0 {
		return ErrNoKey
	}
	query, values := buildQuery(keys)
	_, err := s.db.Table("document").Where(query, values...).Delete(Document{})
	s.Complete = s.isComplete()
	return err
}
