package fexcel

import (
	"fmt"
	"net"
	"os"

	fanuc "github.com/onerobotics/go-fanuc"
)

type Target struct {
	client fanuc.Client

	Name     string
	Comments map[fanuc.Type]map[int]string
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

// TODO: maybe this should be pushed to go-fanuc package
func clientFor(path string, timeout int) (fanuc.Client, error) {
	switch {
	case isDir(path):
		client, err := fanuc.NewFileClient(path)
		if err != nil {
			return nil, err
		}
		return client, nil
	case isIP(path):
		client, err := fanuc.NewHTTPClient(path, timeout)
		if err != nil {
			return nil, err
		}
		return client, nil
	default:
		return nil, fmt.Errorf("%s is not a valid IP address or directory", path)
	}
}

func NewTarget(path string, timeout int) (*Target, error) {
	client, err := clientFor(path, timeout)
	if err != nil {
		return nil, err
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
