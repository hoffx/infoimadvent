package routes

import (
	"log"

	"github.com/go-macaron/session"
	"github.com/hoffx/infoimadvent/config"
	"github.com/hoffx/infoimadvent/storage"
	"golang.org/x/crypto/bcrypt"
	macaron "gopkg.in/macaron.v1"
)

func Account(ctx *macaron.Context, log *log.Logger, sess session.Store, rStorer *storage.RelationStorer, uStorer *storage.UserStorer) {
	defer ctx.HTML(200, "account")

	type Config struct {
		MinPwL, MaxPwL uint
	}

	ctx.Data["Config"] = Config{config.Config.Auth.MinPwLength, config.Config.Auth.MaxPwLength}

	user := sess.Get("user").(storage.User)
	user, err := uStorer.Get(map[string]interface{}{"email": user.Email})
	if err != nil {
		ctx.Data["Error"] = ctx.Tr(ErrDB)
		log.Println(err)
		return
	}
	if user.Teacher {
		ctx.Data["IsTeacher"] = true
	}

	var mode string
	if modeSlice := ctx.Req.URL.Query()["mode"]; len(modeSlice) == 1 {
		mode = modeSlice[0]
	}
	if mode == "changepw" {
		// handle password change (POST) and/or display password-change form
		ctx.Data["ChangePw"] = true
		if ctx.Req.Method == "POST" {

			fPwOld := ctx.Req.FormValue("pwOld")
			fPw := ctx.Req.FormValue("pw1")
			fPw2 := ctx.Req.FormValue("pw2")

			ctx.Data["PwOld"] = fPwOld

			if fPw != fPw2 {
				ctx.Data["Error"] = ctx.Tr(ErrUnequalPasswords)
				return
			}

			ctx.Data["Pw"] = fPw

			err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(fPwOld))
			if err != nil {
				ctx.Data["Error"] = ctx.Tr(ErrWrongCredentials)
				return
			}

			pwBytes, err := bcrypt.GenerateFromPassword([]byte(fPw), bcrypt.DefaultCost)
			if err != nil {
				ctx.Data["Error"] = ctx.Tr(ErrUnexpected)
				return
			}
			user.Hash = string(pwBytes)

			err = uStorer.Put(user)
			if err != nil {
				ctx.Data["Error"] = ctx.Tr(ErrDB)
				log.Println(err)
				return
			}

			err = sess.Set("user", user)
			if err != nil {
				ctx.Data["Error"] = ctx.Tr(ErrUnexpected)
				log.Println(err)
				return
			}

			ctx.Data["Message"] = ctx.Tr(MessChangedPassword)
		}
	} else if mode == "relations" {
		// handle teacher-student connections
		ctx.Data["Relations"] = true
		ctx.Data["Confirm"] = ctx.Tr("confirm")
		ctx.Data["NotConfirmed"] = ctx.Tr("not_confirmed")
		ctx.Data["Remove"] = ctx.Tr("remove")

		// handle confirm and unconfirm/remove actions
		action := ctx.Req.URL.Query()["action"]
		email := ctx.Req.URL.Query()["email"]
		if len(action) == 1 && len(email) == 1 {
			switch action[0] {
			case "confirm":
				// avoid abuse by bad teachers
				if user.Teacher {
					ctx.Data["Error"] = ctx.Tr(ErrIllegalInput)
					return
				}
				// set connection's confirmed status to true
				err = rStorer.Put(storage.Relation{email[0], user.Email, true})
				if err == storage.ErrNoEffect {
					ctx.Data["Error"] = ctx.Tr(ErrUserNotFound)
					return
				} else if err != nil {
					ctx.Data["Error"] = ctx.Tr(ErrDB)
					log.Println(err)
					return
				}
				ctx.Data["Message"] = ctx.Tr(MessUserConfirmed)
			case "delete":
				// remove (for teachers) or disable (for students) connection
				if user.Teacher {
					err = rStorer.Delete(map[string]interface{}{"teacher": user.Email, "student": email[0]})
					if err != nil {
						ctx.Data["Error"] = ctx.Tr(ErrDB)
						log.Println(err)
						return
					}
					ctx.Data["Message"] = ctx.Tr(MessUserRemoved)
				} else {
					err = rStorer.Put(storage.Relation{email[0], user.Email, false})
					if err == storage.ErrNoEffect {
						ctx.Data["Error"] = ctx.Tr(ErrUserNotFound)
						return
					} else if err != nil {
						ctx.Data["Error"] = ctx.Tr(ErrDB)
						log.Println(err)
						return
					}
					ctx.Data["Message"] = ctx.Tr(MessUserUnconfirmed)
				}
			}
		}

		if user.Teacher {
			ctx.Data["IsTeacher"] = true
			// defer creation of connection-list
			defer func() {
				relations, err := rStorer.GetAll(map[string]interface{}{"teacher": user.Email})
				if err != nil {
					ctx.Data["Error"] = ctx.Tr(ErrDB)
					log.Println(err)
					return
				}

				// add student's scores to relations
				type TemplRelation struct {
					storage.Relation
					Score int
				}
				var tRelations []TemplRelation

				for _, r := range relations {
					var student storage.User
					if r.Confirmed {
						student, err = uStorer.Get(map[string]interface{}{"email": r.Student})
						if err != nil {
							ctx.Data["Error"] = ctx.Tr(ErrDB)
							log.Println(err)
							return
						}
					}
					tRelations = append(tRelations, TemplRelation{r, student.Score})
				}

				ctx.Data["RelationsList"] = tRelations
			}()
			if ctx.Req.Method == "POST" {
				// get requested user from db
				fEmail := ctx.Req.FormValue("email")
				reqestedU, err := uStorer.Get(map[string]interface{}{"email": fEmail})
				if err != nil {
					ctx.Data["Error"] = ctx.Tr(ErrDB)
					log.Println(err)
					return
				} else if reqestedU.Email == "" {
					ctx.Data["Error"] = ctx.Tr(ErrUserNotFound)
					return
				} else if reqestedU.Teacher {
					ctx.Data["Error"] = ctx.Tr(ErrNoStudent)
					return
				}
				relation, err := rStorer.Get(map[string]interface{}{"teacher": user.Email, "student": reqestedU.Email})
				if err != nil {
					ctx.Data["Error"] = ctx.Tr(ErrDB)
					log.Println(err)
					return
				} else if relation.Teacher != "" {
					ctx.Data["Error"] = ctx.Tr(ErrRelationExists)
					return
				}
				err = rStorer.Create(storage.Relation{user.Email, reqestedU.Email, false})
				if err != nil {
					ctx.Data["Error"] = ctx.Tr(ErrDB)
					log.Println(err)
					return
				}
				ctx.Data["Message"] = ctx.Tr(MessUserAssigned)
			}
		} else {
			relations, err := rStorer.GetAll(map[string]interface{}{"student": user.Email})
			if err != nil {
				ctx.Data["Error"] = ctx.Tr(ErrDB)
				log.Println(err)
				return
			}
			ctx.Data["RelationsList"] = relations
		}
	} else {
		// display user's score
		ctx.Data["Score"] = true
		ctx.Data["ScoreVal"] = user.Score
	}

	if certificateReady() {
		ctx.Data["Certificate"] = true
	}
}
