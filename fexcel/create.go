package fexcel

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type Creator struct {
	file   *File
	target *Target
}

func NewCreator(path string, cfg Config, targetPath string) (*Creator, error) {
	if filepath.Ext(path) != ".xlsx" {
		return nil, errors.New("File path must end in .xlsx")
	}

	if _, err := os.Stat(path); err == nil {
		return nil, errors.New("File already exists")
	} else {
		// ok if err is IsNotExist
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	f, err := NewFile(path, cfg)
	if err != nil {
		return nil, err
	}
	f.New()

	t, err := NewTarget(targetPath, 500) // TODO: push timeout to config
	if err != nil {
		return nil, err
	}

	return &Creator{file: f, target: t}, nil
}

func (c *Creator) Create(w io.Writer) error {
	fmt.Fprintf(w, "Creating file: %s\n", c.file.path)

	for typ, location := range c.file.Locations {
		fmt.Fprintf(w, "Reading target %s comments\n", typ.VerboseName())
		err := c.target.GetComments(typ)
		if err != nil {
			return err
		}

		// create sheet if necessary
		c.file.CreateSheet(location.Sheet)

		// get start position
		col, row, err := excelize.CellNameToCoordinates(location.Axis)
		if err != nil {
			return err
		}

		// maps are not ordered, so let's create an ids slice we can sort
		var ids []int
		for id, _ := range c.target.Comments[typ] {
			ids = append(ids, id)
		}
		sort.Ints(ids)

		fmt.Fprintf(w, "Writing %d %s comments\n", len(ids), typ.VerboseName())
		for _, id := range ids {
			comment := c.target.Comments[typ][id]
			err := c.file.SetValue(location.Sheet, col, row, id)
			if err != nil {
				return err
			}

			err = c.file.SetValue(location.Sheet, col+c.file.Config.Offset, row, comment)
			if err != nil {
				return err
			}

			row++
		}
	}

	fmt.Fprintln(w, "Saving file.")
	return c.file.Save()
}
