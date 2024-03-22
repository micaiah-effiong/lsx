package terminal

import (
	"errors"
	"io/fs"
	"log"
	"os"

	"github.com/micaiah-effiong/lsx/render"
)

func GetPathEntries(path string) ([]render.Entry, error) {
	file_system := os.DirFS(path)
	ls, err := fs.ReadDir(file_system, ".")

	if err != nil {
		print(err)
		log.Fatal("No such file or directory")

		return nil, errors.New("No such file or directory")
	}

	var render_entries []render.Entry

	for _, dir_list_item := range ls {
		entry, err := render.MakeEntry(dir_list_item)

		if err != nil || entry.IsDotEntry {
			continue
		}

		render_entries = append(render_entries, entry)
	}

	return render_entries, nil
}
