package fexcel

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/olekukonko/tablewriter"
	fanuc "github.com/onerobotics/go-fanuc"
)

type DiffCommand struct {
	fpath string

	f       *File
	targets []Target
}

type Target struct {
	client fanuc.Client

	Name     string
	Comments map[fanuc.Type]map[int]string
}

func (t *Target) GetComments(typ fanuc.Type) error {
	t.Comments[typ] = make(map[int]string)

	switch typ {
	case fanuc.Numreg:
		numregs, err := t.client.NumericRegisters()
		if err != nil {
			return err
		}
		for _, r := range numregs {
			t.Comments[typ][r.Id] = r.Comment
		}
	case fanuc.Posreg:
		posregs, err := t.client.PositionRegisters()
		if err != nil {
			return err
		}
		for _, r := range posregs {
			t.Comments[typ][r.Id] = r.Comment
		}
	case fanuc.Ain, fanuc.Aout, fanuc.Din, fanuc.Dout, fanuc.Flag, fanuc.Gin, fanuc.Gout, fanuc.Rin, fanuc.Rout:
		ports, err := t.client.IO(typ)
		if err != nil {
			return err
		}
		for _, r := range ports {
			t.Comments[typ][r.Id] = r.Comment
		}
	}

	return nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return true
	}

	return false
}

func isIP(path string) bool {
	ip := net.ParseIP(path)
	if ip == nil {
		return false
	}

	return true
}

func NewTarget(path string, timeout int) (t Target, err error) {
	var client fanuc.Client

	switch {
	case isDir(path):
		client, err = fanuc.NewFileClient(path)
		if err != nil {
			return
		}
		t.Name = path
	case isIP(path):
		client, err = fanuc.NewHTTPClient(path, timeout)
		if err != nil {
			return
		}
		t.Name = path // TODO: filepath.Base()?
	default:
		err = fmt.Errorf("%s is not a valid IP address or directory", path)
		return
	}
	t.client = client

	t.Comments = make(map[fanuc.Type]map[int]string)

	return
}

func NewDiffCommand(fpath string, fileConfig Config, timeout int, targetPaths ...string) (*DiffCommand, error) {
	if len(targetPaths) == 0 {
		return nil, fmt.Errorf("Need at least one target")
	}

	d := DiffCommand{fpath: fpath}

	for _, path := range targetPaths {
		t, err := NewTarget(path, timeout)
		if err != nil {
			return nil, err
		}
		d.targets = append(d.targets, t)
	}

	f, err := NewFile(fpath, fileConfig)
	if err != nil {
		return nil, err
	}
	err = f.Open()
	if err != nil {
		return nil, err
	}
	d.f = f

	return &d, nil
}

type Comparison struct {
	Id   int
	Want string
	Got  map[string]string
}

func (c Comparison) Equal() bool {
	for _, got := range c.Got {
		if got != c.Want {
			return false
		}
	}

	return true
}

func (c Comparison) row() []string {
	diff := " "
	if !c.Equal() {
		diff = "X"
	}

	row := []string{strconv.Itoa(c.Id), diff, c.Want}
	for _, s := range c.Got {
		row = append(row, s)
	}

	return row
}

func (d *DiffCommand) Compare(t fanuc.Type) (comparisons []Comparison, err error) {
	definitions, err := d.f.Definitions(t)
	if err != nil {
		return
	}
	if len(definitions) == 0 {
		return
	}

	for _, target := range d.targets {
		err = target.GetComments(t)
		if err != nil {
			return
		}
	}

	// let's only diff the ones defined in the spreadsheet
	for _, def := range definitions {
		c := Comparison{Id: def.Id, Want: def.Comment}
		c.Got = make(map[string]string)

		for _, target := range d.targets {
			got := "undefined"
			if comment, ok := target.Comments[t][def.Id]; ok {
				got = comment
			}
			c.Got[target.Name] = got

		}

		comparisons = append(comparisons, c)
	}

	return
}

func (d *DiffCommand) FprintTable(w io.Writer, t fanuc.Type) error {
	comparisons, err := d.Compare(t)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "%ss\n", t.VerboseName())
	table := tablewriter.NewWriter(w)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)

	header := []string{"Id", "Diff", filepath.Base(d.fpath)}
	for _, target := range d.targets {
		header = append(header, target.Name)
	}
	table.SetHeader(header)

	for _, c := range comparisons {
		table.Append(c.row())
	}

	table.Render()

	return nil
}

func (d *DiffCommand) Execute(w io.Writer) error {
	if len(d.f.Locations) == 0 {
		fmt.Fprintln(w, "No cell location defined")
		return nil
	}

	for dataType, _ := range d.f.Locations {
		err := d.FprintTable(w, dataType)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, "")
	}

	return nil
}
