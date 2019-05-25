package main

import (
	"fmt"

	"crypto/md5"

	"github.com/ddosakura/ghost"
	"github.com/ddosakura/ghost/cmd"
	"github.com/ddosakura/ghost/cmd/proto/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// master data model
const (
	DeprecatedVersion = "0.0.0-2"
	CurrentVersion    = "0.0.0-3"

	cDBver    = "master.version" // proto/model Version
	cDBauth   = "master.auth"    // proto/model User
	cDBuser   = "master.user"    // proto/model User
	cDBdomain = "master.domain"  // proto/model Domain
	cDBconfig = "master.config"  // proto/model ServerConfig
	cDBinvite = "master.invite"  // proto/model InviteCode
)

type db struct {
	session *mgo.Session
	dbName  string
}

func md5c(s string) []byte {
	bs := md5.Sum([]byte(s))
	return bs[:]
}

func newRepo(session *mgo.Session) *db {
	return &db{
		session: session,
	}
}

func (m *db) init(dbName string) error {
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

func (m *db) _checkModelVersion() (e error) {
	defer func() {
		if err := recover(); err != nil {
			// e = err.(error)
			e = cmd.ErrModelVersion
		}
	}()

	db := m.session.DB(m.dbName)
	var ver model.Version
	if err := db.C(cDBver).Find(bson.M{}).One(&ver); err != nil {
		ghost.Error(err)
	}
	ghost.Info("Model Version:", ver.CurrentVersion)
	v := cmd.NewVer(ver.CurrentVersion)

	// <= DeprecatedVersion
	if v.Compare(DeprecatedVersion) < 1 {
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
	if v.Compare(CurrentVersion) == 1 {
		ghost.Error(fmt.Sprintf("current model version is v%s, but db model version is v%s", CurrentVersion, v))
	}

	return
}

func (m *db) _init() error {
	db := m.session.DB(m.dbName)

	ver := model.Version{
		CurrentVersion: CurrentVersion,
	}
	// "master.version"
	if err := db.C(cDBver).Insert(&ver); err != nil {
		return err
	}

	// "master.auth"
	var pass []byte
	if p, has := upServiceData.user.Password(); has {
		pass = md5c(p)
	}
	auth := model.User{
		User:   upServiceData.user.Username(),
		Pass:   pass,
		Domain: []string{upServiceData.host},
	}
	if err := db.C(cDBauth).Insert(&auth); err != nil {
		return err
	}

	// "master.domain"
	d0 := &model.Domain{
		Name: upServiceData.user.Username(),
		Jump: upServiceData.host,
	}
	d1 := &model.Domain{
		Name: upServiceData.host,
		IP:   "",
	}
	if err := db.C(cDBdomain).Insert(d0, d1); err != nil {
		return err
	}

	// "master.config"
	c := &model.ServerConfig{
		UserMode:         model.UserMode_REGIST_ADMIN,
		MaxUser:          10,
		MaxDomain:        100,
		MaxDomainPerUser: 10,
	}
	if err := db.C(cDBconfig).Insert(c); err != nil {
		return err
	}

	// "master.user"
	// "master.invite"

	return nil
}

func (m *db) conn(fn func(*mgo.Database)) {
	s := m.session.Clone()
	defer s.Close()
	fn(s.DB(m.dbName))
}
