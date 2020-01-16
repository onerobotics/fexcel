package fexcel

import (
	"fmt"
	"sync"
	"time"

	"github.com/onerobotics/comtool"
	fanuc "github.com/onerobotics/go-fanuc"
)

type Setter interface {
	Set(Definition, string) error
}

type CommentToolSetter struct {
	Timeout time.Duration
}

func (c *CommentToolSetter) Set(d Definition, host string) error {
	fcode := codeForType(d.Type)
	return comtool.Set(fcode, d.Id, d.Comment, host, c.Timeout)
}

type MultiSetter struct {
	Hosts    map[string]bool
	Warnings []string
	Errors   map[string][]string // key is host
	Setter

	wMux sync.Mutex
	eMux sync.Mutex
}

func NewMultiSetter(hosts []string, u Setter) *MultiSetter {
	m := &MultiSetter{
		Setter: u,
	}

	m.Hosts = make(map[string]bool)
	m.Errors = make(map[string][]string)

	for _, h := range hosts {
		m.Hosts[h] = true
	}

	return m
}

// translate between fanuc.Type and comtool.FunctionCode
func codeForType(d fanuc.Type) comtool.FunctionCode {
	switch d {
	case fanuc.Numreg:
		return comtool.NUMREG
	case fanuc.Posreg:
		return comtool.POSREG
	case fanuc.Ualm:
		return comtool.UALM
	case fanuc.Rin:
		return comtool.RIN
	case fanuc.Rout:
		return comtool.ROUT
	case fanuc.Din:
		return comtool.DIN
	case fanuc.Dout:
		return comtool.DOUT
	case fanuc.Gin:
		return comtool.GIN
	case fanuc.Gout:
		return comtool.GOUT
	case fanuc.Ain:
		return comtool.AIN
	case fanuc.Aout:
		return comtool.AOUT
	case fanuc.Sreg:
		return comtool.SREG
	case fanuc.Flag:
		return comtool.FLAG
	}

	return comtool.FunctionCode(0) // invalid
}

func maxLengthFor(t fanuc.Type) int {
	switch t {
	case fanuc.Numreg, fanuc.Posreg, fanuc.Sreg:
		return 16
	case fanuc.Ualm:
		return 29
	default:
		return 24
	}
}

func (c *MultiSetter) warn(msg string) {
	c.wMux.Lock()
	c.Warnings = append(c.Warnings, msg)
	c.wMux.Unlock()
}

const maxErrors = 5

func (c *MultiSetter) error(host string, msg string) {
	c.eMux.Lock()
	c.Errors[host] = append(c.Errors[host], msg)
	if len(c.Errors[host]) >= maxErrors {
		// disable host
		c.Hosts[host] = false
		c.Errors[host] = append(c.Errors[host], "Too many errors. Host disabled.")
	}
	c.eMux.Unlock()
}

func (c *MultiSetter) Set(defs []Definition) error {
	for _, d := range defs {
		if maxLength := maxLengthFor(d.Type); len(d.Comment) > maxLength {
			l := len(d.Comment)
			d.Comment = d.Comment[:maxLength]
			c.warn(fmt.Sprintf("comment for %s[%d] truncated to %q (%d > %d).", d.Type, d.Id, d.Comment, l, maxLength))
		}
	}

	var wg sync.WaitGroup
	for host, _ := range c.Hosts {
		wg.Add(1)
		go func(c *MultiSetter, host string, defs []Definition, wg *sync.WaitGroup) {
			defer wg.Done()
			for _, d := range defs {
				if !c.Hosts[host] {
					return
				}

				err := c.Setter.Set(d, host)
				if err != nil {
					c.error(host, fmt.Sprintf("Failed to update %s[%d].", d.Type, d.Id))
				}
			}
		}(c, host, defs, &wg)
	}
	wg.Wait()

	return nil
}
