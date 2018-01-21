package storage

import (
	"errors"

	"github.com/go-xorm/xorm"
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
