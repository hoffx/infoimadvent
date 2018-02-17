package storage

import (
	"os"

	"github.com/go-xorm/core"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type UserStorer struct {
	Storer
}

type User struct {
	Email             string
	Hash              string
	Grade             uint
	Active            bool
	Confirmed         bool
	ConfirmationToken string
	Teacher           bool
	Days              []int
	Score             int
	IsAdmin           bool
	Lang              string
}

func NewUserStorer(name, user, password string, doLog bool) (UserStorer, error) {
	db, err := xorm.NewEngine("mysql", user+":"+password+"@/"+name+"?charset=utf8")
	if err != nil {
		return UserStorer{}, err
	}
	if doLog {
		db.ShowSQL(true)
		db.Logger().SetLevel(core.LOG_DEBUG)
		f, err := os.Create("sql.log")
		if err != nil {
			return UserStorer{}, err
		}
		db.SetLogger(xorm.NewSimpleLogger(f))
	}
	err = db.Sync(new(User))
	if err != nil {
		return UserStorer{}, err
	}
	err = db.CreateTables(&User{})
	if err != nil {
		return UserStorer{}, err
	}

	return UserStorer{Storer{db}}, nil
}

func (s *UserStorer) ResetDB() error {
	err := s.db.DropTables(User{})
	if err != nil {
		return err
	}
	err = s.db.CreateTables(User{})
	return err
}

func (s *UserStorer) Create(user User) error {
	_, err := s.db.Insert(user)
	return err
}

func (s *UserStorer) Put(user User) error {
	oldUser, err := s.Get(map[string]interface{}{"email": user.Email})
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

func (s *UserStorer) Get(keys map[string]interface{}) (User, error) {
	if len(keys) == 0 {
		return User{}, ErrNoKey
	}
	query, values := buildQuery(keys)
	user := User{}
	_, err := s.db.Table("user").Where(query, values...).Get(&user)
	return user, err
}

func (s *UserStorer) GetAll(keys map[string]interface{}) (users []User, err error) {
	query, values := buildQuery(keys)
	err = s.db.Table("user").Where(query, values...).Find(&users)
	return
}

func (s *UserStorer) Delete(keys map[string]interface{}) error {
	if len(keys) == 0 {
		return ErrNoKey
	}
	query, values := buildQuery(keys)
	_, err := s.db.Table("user").Where(query, values...).Delete(User{})
	return err
}
