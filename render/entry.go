package render

import (
	"io/fs"
	"strings"
)

type Entry struct {
	Name       string
	IsDir      bool
	RenderName string
	IsDotEntry bool
}

func (e Entry) Init(dir_entry fs.DirEntry) (Entry, error) {
	info, err := dir_entry.Info()

	if err != nil {
		var ent Entry
		return ent, err
	}

	e.Name = info.Name()
	e.RenderName = info.Name()
	e.IsDir = info.IsDir()
	e.IsDotEntry = strings.HasPrefix(info.Name(), ".")

	if info.IsDir() {
		e.RenderName = "\033[38;5;112m" + info.Name() + "\033[38;5;231m"
	}

	return e, nil
}

func MakeEntry(dir_entry fs.DirEntry) (Entry, error) {
	e := new(Entry)

	entry, err := e.Init(dir_entry)

	if err != nil {
		return entry, err
	}

	return entry, nil
}
