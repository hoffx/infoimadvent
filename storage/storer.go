package storage

import (
	"errors"
	"os"

	"github.com/go-xorm/core"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type User struct {
	Email              string
	Password           string
	Grade              uint
	Active             bool
	Confirmed          bool
	ConfirmationToken  string
	Teacher            bool
	ConfirmedRelations []string
	RequestedRelations []string
	// Calendar  CalendarInfo
}

type CalendarInfo struct {
}

type Storer struct {
	db     *xorm.Engine
	Active bool
}

var ErrNoEffect = errors.New("no_effect_error")

func NewStorer(name, user, password string, doLog bool) (Storer, error) {
	db, err := xorm.NewEngine("mysql", user+":"+password+"@/"+name+"?charset=utf8")
	if err != nil {
		return Storer{}, err
	}
	if doLog {
		db.ShowSQL(true)
		db.Logger().SetLevel(core.LOG_DEBUG)
		f, err := os.Create("sql.log")
		if err != nil {
			return Storer{}, err
		}
		db.SetLogger(xorm.NewSimpleLogger(f))
	}
	err = db.Sync(new(User))
	if err != nil {
		return Storer{}, err
	}
	err = db.CreateTables(&User{})
	if err != nil {
		return Storer{}, err
	}

	return Storer{db, true}, nil
}

func (s *Storer) ResetDB() error {
	err := s.db.DropTables(&User{})
	if err != nil {
		return err
	}
	err = s.db.CreateTables(&User{})
	return err
}

func (s *Storer) Create(user User) error {
	_, err := s.db.Insert(&user)
	return err
}

func (s *Storer) Put(user User) error {
	oldUser, err := s.Get(user.Email)
	if err != nil {
		return err
	}
	i, err := s.db.Delete(&oldUser)
	if i == 0 {
		err = ErrNoEffect
	}
	if err != nil {
		return err
	}
	return s.Create(user)
}

func (s *Storer) Get(key string) (User, error) {
	var user User
	_, err := s.db.Table("user").Where("email = ?", key).Get(&user)
	return user, err
}
