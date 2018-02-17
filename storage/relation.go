package storage

import (
	"os"

	"github.com/go-xorm/core"

	// blank import required by xorm
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

// RelationStorer is a normal Storer. The extra type is needed
// by macaron for the identification of the storer
type RelationStorer struct {
	Storer
}

// Relation is a connection between a user (type teacher) and
// a user (type student)
type Relation struct {
	Teacher   string
	Student   string
	Confirmed bool
}

// NewRelationStorer generates a RelationStorer. If doLog is true, it
//will create a log file "sql.log"
func NewRelationStorer(name, user, password string, doLog bool) (RelationStorer, error) {
	db, err := xorm.NewEngine("mysql", user+":"+password+"@/"+name+"?charset=utf8")
	if err != nil {
		return RelationStorer{}, err
	}
	if doLog {
		db.ShowSQL(true)
		db.Logger().SetLevel(core.LOG_DEBUG)
		f, err := os.Create("sql.log")
		if err != nil {
			return RelationStorer{}, err
		}
		db.SetLogger(xorm.NewSimpleLogger(f))
	}
	err = db.Sync(new(Relation))
	if err != nil {
		return RelationStorer{}, err
	}
	err = db.CreateTables(&Relation{})
	if err != nil {
		return RelationStorer{}, err
	}

	return RelationStorer{Storer{db}}, nil
}

// ResetDB deletes the relation table
func (s *RelationStorer) ResetDB() error {
	err := s.db.DropTables(Relation{})
	if err != nil {
		return err
	}
	err = s.db.CreateTables(Relation{})
	return err
}

// Create creates a new entry
func (s *RelationStorer) Create(relation Relation) error {
	_, err := s.db.Insert(relation)
	return err
}

// Put modifies a entry. The relation is identified by its
// teacher and student values
func (s *RelationStorer) Put(relation Relation) error {
	oldRelation, err := s.Get(map[string]interface{}{"teacher": relation.Teacher, "student": relation.Student})
	if err != nil {
		return err
	}
	i, err := s.db.Delete(&oldRelation)
	if i == 0 {
		err = ErrNoEffect
	}
	if err != nil {
		return err
	}
	return s.Create(relation)
}

// Get returns one relation defined by the provided keys
func (s *RelationStorer) Get(keys map[string]interface{}) (Relation, error) {
	if len(keys) == 0 {
		return Relation{}, ErrNoKey
	}
	query, values := buildQuery(keys)
	relation := Relation{}
	_, err := s.db.Table("relation").Where(query, values...).Get(&relation)
	return relation, err
}

// GetAll returns a slice of relations defined by the provided keys
func (s *RelationStorer) GetAll(keys map[string]interface{}) (relations []Relation, err error) {
	query, values := buildQuery(keys)
	err = s.db.Table("relation").Where(query, values...).Find(&relations)
	return
}

// Delete deletes all entries defined by the provided keys
func (s *RelationStorer) Delete(keys map[string]interface{}) error {
	if len(keys) == 0 {
		return ErrNoKey
	}
	query, values := buildQuery(keys)
	_, err := s.db.Table("relation").Where(query, values...).Delete(Relation{})
	return err
}
