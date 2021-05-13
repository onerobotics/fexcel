package fexcel

import (
	"errors"
	"time"

	fanuc "github.com/onerobotics/go-fanuc"
)

var fanucType = [...]fanuc.Type{
	Numreg: fanuc.Numreg,
	Posreg: fanuc.Posreg,
	Ualm:   fanuc.Ualm,
	Ain:    fanuc.Ain,
	Aout:   fanuc.Aout,
	Din:    fanuc.Din,
	Dout:   fanuc.Dout,
	Gin:    fanuc.Gin,
	Gout:   fanuc.Gout,
	Rin:    fanuc.Rin,
	Rout:   fanuc.Rout,
	Sreg:   fanuc.Sreg,
	Flag:   fanuc.Flag,
}

type Target struct {
	client fanuc.Client

	Name     string
	Comments map[Type]map[int]string
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
	t.Comments = make(map[Type]map[int]string)

	return &t, nil
}

func (t *Target) GetComments(typ Type) error {
	t.Comments[typ] = make(map[int]string)

	switch typ {
	case Numreg:
		numregs, err := t.client.NumericRegisters()
		if err != nil {
			return err
		}
		for _, r := range numregs {
			t.Comments[typ][r.Id] = r.Comment
		}
	case Posreg:
		posregs, err := t.client.PositionRegisters()
		if err != nil {
			return err
		}
		for _, r := range posregs {
			t.Comments[typ][r.Id] = r.Comment
		}
	case Ain, Aout, Din, Dout, Flag, Gin, Gout, Rin, Rout:
		ports, err := t.client.IO(fanucType[typ])
		if err != nil {
			return err
		}
		for _, r := range ports {
			t.Comments[typ][r.Id] = r.Comment
		}
	}

	return nil
}

func (t *Target) SetComment(typ Type, id int, comment string) error {
	if c, ok := t.client.(*fanuc.HTTPClient); ok {
		return c.SetComment(fanucType[typ], id, comment)
	} else {
		return errors.New("need a fanuc.HTTPClient")
	}
}
