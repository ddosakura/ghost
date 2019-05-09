package sh

import (
	"os"
	"path"

	"github.com/ddosakura/ghost"

	"github.com/mitchellh/go-homedir"
)

// Data Dir
var (
	RootDir      string
	RootDirLog   string
	RootDirCache string
	RootDirTmp   string // Default Work Dir
)

func init() {
	home, err := homedir.Dir()
	if err != nil {
		ghost.Crash(-1, err)
	}

	RootDir = path.Join(home, ".ddosakura")
	DirMustExist(RootDir)
	RootDir = path.Join(RootDir, "ghost")
	DirMustExist(RootDir)

	RootDirLog = subDirMustExist("log")
	RootDirCache = subDirMustExist("cache")
	RootDirTmp = subDirMustExist("tmp")
}

func subDirMustExist(dir string) string {
	d := path.Join(RootDir, dir)
	DirMustExist(d)
	return d
}

// DirMustExist will exit when can't fix the dir
func DirMustExist(dir string) {
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err = os.Mkdir(dir, 0755); err != nil {
			ghost.Crash(-1, err)
		}
	} else if err != nil {
		ghost.Crash(-1, err)
	} else if !stat.IsDir() {
		ghost.Crash(-1, ErrDirIsFile, "@", dir)
	}
}
