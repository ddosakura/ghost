package main

import (
	"fmt"

	"github.com/ddosakura/ghost"
	"github.com/ddosakura/ghost/cmd"
	"github.com/kr/pretty"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// master data model
const (
	currentVersion = "0.0.0-1"

	cDBver  = "master.version"
	cDBauth = "master.auth"
	cDBuser = "master.user"
)

type model struct {
	session *mgo.Session
	dbName  string
}

// use this to check model version, so it can't be changed!
type cDBverT struct {
	currentVersion string
}

type cDBauthT struct {
	rootUser string
	rootPass string
}

type cDBuserT struct {
	user string
	pass string
}

func newModel(session *mgo.Session) *model {
	return &model{
		session: session,
	}
}

func (m *model) init(dbName string) error {
	ds, e := m.session.DatabaseNames()
	//fmt.Println(ds)
	if e != nil {
		return e
	}
	m.dbName = dbName

	hasInit := false
	for _, dn := range ds {
		if dn == m.dbName {
			hasInit = true
			if e = m._checkModelVersion(); e != nil {
				return e
			}
		}
	}
	if hasInit {
		return nil
	}

	return m._init()
}

func (m *model) _checkModelVersion() (e error) {
	defer func() {
		if err := recover(); err != nil {
			// e = err.(error)
			e = cmd.ErrModelVersion
		}
	}()

	db := m.session.DB(m.dbName)
	var ver cDBverT
	if err := db.C(cDBver).Find(bson.M{}).One(&ver); err != nil {
		ghost.Error(err)
	}
	ghost.Info("Model Version:", ver.currentVersion)
	v := cmd.NewVer(ver.currentVersion)

	// <= v0.0.0-0
	if v.Compare("0") < 1 {
		ghost.Warn(fmt.Sprintf("v%s is no longer supported", v))
		fmt.Print("Drop it, and re-init?(y/N)")
		var action string
		fmt.Scanln(&action)
		switch action {
		case "y", "Y", "yes", "YES":
			if err := db.DropDatabase(); err != nil {
				ghost.Crash(-1, err)
			}
			if err := m._init(); err != nil {
				ghost.Crash(-1, err)
			}
		default:
			ghost.Error("Please choose other version!")
		}
	}

	// == currentVersion
	if v.Compare(currentVersion) == 1 {
		ghost.Error(fmt.Sprintf("current model version is v%s, but db model version is v%s", currentVersion, v))
	}

	return
}

func (m *model) _init() error {
	db := m.session.DB(m.dbName)
	if err := db.C(cDBver).Insert(&cDBverT{currentVersion}); err != nil {
		return err
	}

	var is []interface{}
	db.C(cDBver).Find(bson.M{}).All(&is)
	pretty.Println(is)

	// TODO: ...
	return nil
}

func (m *model) repo(fn func(*mgo.Database)) {
	s := m.session.Clone()
	defer s.Close()
	fn(s.DB(m.dbName))
}