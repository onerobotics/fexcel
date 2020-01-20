package fexcel

import (
	"fmt"
	"sync"

	fanuc "github.com/onerobotics/go-fanuc"
)

type SetCommand struct {
	fpath   string
	file    *File
	targets []*Target

	Definitions map[fanuc.Type][]Definition
	Errors      map[string]*errorList
}

func NewSetCommand(fpath string, cfg Config, targets ...string) (*SetCommand, error) {
	if len(targets) == 0 {
		return nil, fmt.Errorf("Need at least one target")
	}

	err := cfg.FileConfig.Validate()
	if err != nil {
		return nil, err
	}

	s := SetCommand{fpath: fpath}

	for _, path := range targets {
		t, err := NewTarget(path, cfg.Timeout)
		if err != nil {
			return nil, err
		}

		if _, ok := t.client.(*fanuc.HTTPClient); !ok {
			return nil, fmt.Errorf("%q is not a valid remote host", path)
		}

		s.targets = append(s.targets, t)
	}

	f, err := OpenFile(fpath, cfg.FileConfig)
	if err != nil {
		return nil, err
	}
	s.file = f

	s.Definitions, err = s.file.AllDefinitions()
	if err != nil {
		return nil, err
	}

	s.Errors = make(map[string]*errorList)
	for _, host := range s.Hosts() {
		s.Errors[host] = &errorList{}
	}

	return &s, nil
}

// host -> type -> count
type setResult struct {
	Counts map[string]map[fanuc.Type]int
	mux    sync.Mutex
}

func (s *setResult) Inc(host string, t fanuc.Type) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.Counts[host][t]++
}

func newSetResult(targets []*Target) *setResult {
	var result setResult
	result.Counts = make(map[string]map[fanuc.Type]int)
	for _, t := range targets {
		result.Counts[t.Name] = make(map[fanuc.Type]int)
	}
	return &result
}

func (s *SetCommand) Hosts() []string {
	var hosts []string
	for _, t := range s.targets {
		hosts = append(hosts, t.Name)
	}
	return hosts
}

func (s *SetCommand) Set(wg *sync.WaitGroup, target *Target, result *setResult) {
	defer wg.Done()

	for typ, defs := range s.Definitions {
		err := target.GetComments(typ)
		if err != nil {
			s.Errors[target.Name].Add(err)
			return
		}

		for _, def := range defs {
			want := Truncated(def.Comment, typ)
			if got, ok := target.Comments[typ][def.Id]; ok && got == want {
				continue
			}

			err := target.SetComment(typ, def.Id, want)
			if err != nil {
				s.Errors[target.Name].Add(err)
				return
			} else {
				result.Inc(target.Name, typ)
			}
		}
	}
}

func (s *SetCommand) Execute() (*setResult, error) {
	result := newSetResult(s.targets)

	var wg sync.WaitGroup
	for _, target := range s.targets {
		wg.Add(1)
		go s.Set(&wg, target, result)
	}
	wg.Wait()

	return result, s.Err()
}

func (s SetCommand) Error() string {
	var str string
	for host, err := range s.Errors {
		if err.Err() != nil {
			str += fmt.Sprintf("%s: %s\n", host, err.Error())
		}
	}
	return str
}

func (s SetCommand) Err() error {
	if s.Error() == "" {
		return nil
	}
	return s
}
