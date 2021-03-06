package fexcel

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type Creator struct {
	file    *File
	target  *Target
	headers bool
}

func NewCreator(path string, cfg Config, headers bool, targetPath string) (*Creator, error) {
	if filepath.Ext(path) != ".xlsx" {
		return nil, errors.New("File path must end in .xlsx")
	}

	hasOverlaps, err := cfg.HasOverlaps()
	if err != nil {
		return nil, err
	}
	if hasOverlaps {
		return nil, errors.New("configuration has overlapping columns")
	}

	if headers {
		err = cfg.CheckHeaders()
		if err != nil {
			return nil, err
		}
	}

	f, err := NewFile(path, cfg.FileConfig)
	if err != nil {
		return nil, err
	}

	t, err := NewTarget(targetPath, cfg.Timeout)
	if err != nil {
		return nil, err
	}

	return &Creator{file: f, target: t, headers: headers}, nil
}

func (c *Creator) Create(w io.Writer) error {
	fmt.Fprintf(w, "Creating file: %s\n", c.file.path)

	for t, location := range c.file.Locations {
		fmt.Fprintf(w, "Reading target %s comments\n", t)
		err := c.target.GetComments(t)
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

		// write header
		if c.headers {
			err = c.file.SetValue(location.Sheet, col, row-1, t.String()+"s")
			if err != nil {
				return err
			}
		}

		// maps are not ordered, so let's create an ids slice we can sort
		var ids []int
		for id, _ := range c.target.Comments[t] {
			ids = append(ids, id)
		}
		sort.Ints(ids)

		fmt.Fprintf(w, "Writing %d %s comments\n", len(ids), t)
		for _, id := range ids {
			comment := c.target.Comments[t][id]
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
