package fexcel

import (
	"time"

	fanuc "github.com/onerobotics/go-fanuc"
)

type Target struct {
	client fanuc.Client

	Name     string
	Comments map[fanuc.Type]map[int]string
}

func NewTarget(path string, timeout int) (*Target, error) {
	client, err := fanuc.NewClient(path)
	if err != nil {
		return nil, err
	}
	if c, ok := client.(*fanuc.HTTPClient); ok {
		c.SetTimeout(time.Duration(timeout) * time.Second)
	}

	var t Target
	t.client = client
	t.Name = path
	t.Comments = make(map[fanuc.Type]map[int]string)

	return &t, nil
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
