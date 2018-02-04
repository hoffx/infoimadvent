package storage

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/hoffx/infoimadvent/config"
	macaron "gopkg.in/macaron.v1"
)

// answers
const (
	None = iota
	A
	B
	C
	D
)

// scores
const (
	// calculation
	Right   = 10
	Wrong   = 0
	Missing = 2
	// grading
	Full       = 240
	Insane     = 220
	Incredible = 200
	Good       = 160
	Ok         = 120
)

// types
const (
	Quest = iota
	About
	ToS
)

var ErrNoEffect = errors.New("no_effect_error")
var ErrNoKey = errors.New("no_key_error")

type Storer struct {
	db     *xorm.Engine
	Active bool
}

func buildQuery(keys map[string]interface{}) (query string, values []interface{}) {
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
	return
}

func ResetUsers(uStorer *UserStorer, rStorer *RelationStorer) (err error) {
	err = os.RemoveAll(config.Config.Sessioner.StoragePath)
	if err != nil {
		return
	}
	err = os.Mkdir(config.Config.Sessioner.StoragePath, os.ModePerm)
	if err != nil {
		return
	}
	_, err = os.Create(config.Config.Sessioner.StoragePath + "/keep.me")
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		return
	}
	err = uStorer.ResetDB()
	if err != nil {
		return
	}

	// write admin to db
	err = uStorer.Create(User{config.Config.Auth.AdminMail, config.Config.Auth.AdminHash, config.Config.Grades.Max, true, true, "", true, make([]int, 24), 0, true, "en-US"})
	if err != nil {
		return
	}

	err = rStorer.ResetDB()
	return
}

func ResetDocuments(dStorer *DocumentStorer, questsOnly bool) (err error) {
	if questsOnly {
		var docs []Document
		docs, err = dStorer.GetAll(map[string]interface{}{"type": Quest})
		if err != nil {
			return
		}
		files := make(map[string]bool, 0)
		for _, d := range docs {
			files[d.Path] = true
		}
		for k := range files {
			err = os.Remove(k)
			if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
				return
			}
			err = os.RemoveAll(config.Config.FileSystem.AssetsStoragePath + "/" + path.Base(k))
			if err != nil {
				return
			}
		}
		err = dStorer.Delete(map[string]interface{}{"type": Quest})
	} else {
		err = os.RemoveAll(config.Config.FileSystem.MDStoragePath)
		if err != nil {
			return
		}
		err = os.Mkdir(config.Config.FileSystem.MDStoragePath, os.ModePerm)
		if err != nil {
			return
		}
		_, err = os.Create(config.Config.FileSystem.MDStoragePath + "/keep.me")
		if err != nil && !strings.Contains(err.Error(), "file exists") {
			return
		}
		err = os.RemoveAll(config.Config.FileSystem.AssetsStoragePath)
		if err != nil {
			return
		}
		err = os.Mkdir(config.Config.FileSystem.AssetsStoragePath, os.ModePerm)
		if err != nil {
			return
		}
		_, err = os.Create(config.Config.FileSystem.AssetsStoragePath + "/keep.me")
		if err != nil && !strings.Contains(err.Error(), "file exists") {
			return
		}
		err = dStorer.ResetDB()
	}
	return
}

func InitStorers() (u UserStorer, d DocumentStorer, r RelationStorer, err error) {
	u, err = NewUserStorer(config.Config.DB.Name, config.Config.DB.User, config.Config.DB.Password, macaron.Env == macaron.DEV)
	if err != nil {
		return
	}
	d, err = NewDocumentStorer(config.Config.DB.Name, config.Config.DB.User, config.Config.DB.Password, macaron.Env == macaron.DEV)
	if err != nil {
		return
	}
	r, err = NewRelationStorer(config.Config.DB.Name, config.Config.DB.User, config.Config.DB.Password, macaron.Env == macaron.DEV)
	if err != nil {
		return
	}
	return
}
