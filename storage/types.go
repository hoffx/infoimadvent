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
	Days               []int
	Score              int
}

type Quest struct {
	Path     string
	Grade    int
	Day      int
	Solution int
}

type Storer struct {
	db     *xorm.Engine
	Active bool
}

type UserStorer struct {
	Storer
}

type QuestStorer struct {
	Storer
}

var ErrNoEffect = errors.New("no_effect_error")
var ErrNoKey = errors.New("no_key_error")
