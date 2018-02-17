package storage

import (
	"os"

	"github.com/go-xorm/core"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type RelationStorer struct {
	Storer
}

type Relation struct {
	Teacher   string
	Student   string
	Confirmed bool
}

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

func (s *RelationStorer) ResetDB() error {
	err := s.db.DropTables(Relation{})
	if err != nil {
		return err
	}
	err = s.db.CreateTables(Relation{})
	return err
}

func (s *RelationStorer) Create(relation Relation) error {
	_, err := s.db.Insert(relation)
	return err
}

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

func (s *RelationStorer) Get(keys map[string]interface{}) (Relation, error) {
	if len(keys) == 0 {
		return Relation{}, ErrNoKey
	}
	query, values := buildQuery(keys)
	relation := Relation{}
	_, err := s.db.Table("relation").Where(query, values...).Get(&relation)
	return relation, err
}

func (s *RelationStorer) GetAll(keys map[string]interface{}) (relations []Relation, err error) {
	query, values := buildQuery(keys)
	err = s.db.Table("relation").Where(query, values...).Find(&relations)
	return
}

func (s *RelationStorer) Delete(keys map[string]interface{}) error {
	if len(keys) == 0 {
		return ErrNoKey
	}
	query, values := buildQuery(keys)
	_, err := s.db.Table("relation").Where(query, values...).Delete(Relation{})
	return err
}
