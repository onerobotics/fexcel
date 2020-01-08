package fanuc

import (
	"fmt"
	"sync"
	"time"

	"github.com/onerobotics/comtool"
)

type Definition struct {
	DataType
	Id      int
	Comment string
}

type Updater interface {
	Update(Definition, string) error
}

type CommentToolUpdater struct {
	Timeout time.Duration
}

func (c *CommentToolUpdater) Update(d Definition, host string) error {
	fcode := codeForDataType(d.DataType)
	return comtool.Set(fcode, d.Id, d.Comment, host, c.Timeout)
}

type MultiUpdater struct {
	Hosts    map[string]bool
	Warnings []string
	Errors   map[string][]string // key is host
	Updater

	wMux sync.Mutex
	eMux sync.Mutex
}

func NewMultiUpdater(hosts []string, u Updater) *MultiUpdater {
	m := &MultiUpdater{
		Updater: u,
	}

	m.Hosts = make(map[string]bool)
	m.Errors = make(map[string][]string)

	for _, h := range hosts {
		m.Hosts[h] = true
	}

	return m
}

// translate between fanuc.DataType and comtool.FunctionCode
func codeForDataType(d DataType) comtool.FunctionCode {
	switch d {
	case Numreg:
		return comtool.NUMREG
	case Posreg:
		return comtool.POSREG
	case Ualm:
		return comtool.UALM
	case Rin:
		return comtool.RIN
	case Rout:
		return comtool.ROUT
	case Din:
		return comtool.DIN
	case Dout:
		return comtool.DOUT
	case Gin:
		return comtool.GIN
	case Gout:
		return comtool.GOUT
	case Ain:
		return comtool.AIN
	case Aout:
		return comtool.AOUT
	case Sreg:
		return comtool.SREG
	case Flag:
		return comtool.FLAG
	}

	return comtool.FunctionCode(0) // invalid
}

func maxLengthFor(dataType DataType) int {
	switch dataType {
	case Numreg, Posreg, Sreg:
		return 16
	case Ualm:
		return 29
	default:
		return 24
	}
}

func (c *MultiUpdater) warn(msg string) {
	c.wMux.Lock()
	c.Warnings = append(c.Warnings, msg)
	c.wMux.Unlock()
}

const maxErrors = 5

func (c *MultiUpdater) error(host string, msg string) {
	c.eMux.Lock()
	c.Errors[host] = append(c.Errors[host], msg)
	if len(c.Errors[host]) >= maxErrors {
		// disable host
		c.Hosts[host] = false
		c.Errors[host] = append(c.Errors[host], "Too many errors. Host disabled.")
	}
	c.eMux.Unlock()
}

func (c *MultiUpdater) Update(defs []Definition) error {
	for _, d := range defs {
		if maxLength := maxLengthFor(d.DataType); len(d.Comment) > maxLength {
			l := len(d.Comment)
			d.Comment = d.Comment[:maxLength]
			c.warn(fmt.Sprintf("comment for %s[%d] truncated to %q (%d > %d).", d.DataType, d.Id, d.Comment, l, maxLength))
		}
	}

	var wg sync.WaitGroup
	for host, _ := range c.Hosts {
		wg.Add(1)
		go func(c *MultiUpdater, host string, defs []Definition, wg *sync.WaitGroup) {
			defer wg.Done()
			for _, d := range defs {
				if !c.Hosts[host] {
					return
				}

				err := c.Updater.Update(d, host)
				if err != nil {
					c.error(host, fmt.Sprintf("Failed to update %s[%d].", d.DataType, d.Id))
				}
			}
		}(c, host, defs, &wg)
	}
	wg.Wait()

	return nil
}
