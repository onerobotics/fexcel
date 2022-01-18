package compile

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/onerobotics/fexcel/fexcel"
)

const CACHE_FILE_NAME = ".fexcelcache"

var ErrCacheInvalid = errors.New("invalid cache")

// A cache caches compiler definitions and constants to a local file.
type cache struct {
	Path        string
	FileConfig  fexcel.FileConfig
	Version     string
	ModifiedAt  time.Time
	Definitions map[string]map[string]int
	Constants   map[string]string
}

// save saves the cache to a local file
func (c *cache) save() error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(CACHE_FILE_NAME, bytes, 0644)
}

// refresh grabs the most recent data from the spreadsheet.
func (c *cache) refresh() error {
	spreadsheet, err := fexcel.OpenFile(c.Path, c.FileConfig)
	if err != nil {
		return err
	}

	allDefs, err := spreadsheet.AllDefinitions()
	if err != nil {
		return err
	}

	for _, l := range spreadsheet.Locations {
		switch l.Type {
		case fexcel.Constant:
			// noop
		default:
			c.Definitions[l.Type.String()] = make(map[string]int)
		}
	}

	for t, defs := range allDefs {
		for _, def := range defs {
			if def.Comment == "" {
				continue
			}

			c.Definitions[t.String()][def.Comment] = def.Id
		}
	}

	if spreadsheet.Locations[fexcel.Constant] != nil {
		c.Constants, err = spreadsheet.Constants()
		if err != nil {
			return err
		}
	}

	c.Definitions["UI"] = map[string]int{
		"IMSTP":      1,
		"Hold":       2,
		"SFSPD":      3,
		"CycleStop":  4,
		"FaultReset": 5,
		"Start":      6,
		"Home":       7,
		"Enable":     8,
		"ProdStart":  18,
	}
	c.Definitions["UO"] = map[string]int{
		"CmdEnabled":  1,
		"SystemReady": 2,
		"PrgRunning":  3,
		"PrgPaused":   4,
		"MotionHeld":  5,
		"Fault":       6,
		"AtPerch":     7,
		"TPEnabled":   8,
		"BattAlarm":   9,
		"Busy":        10,
	}
	c.Definitions["SI"] = map[string]int{
		"FaultReset": 1,
		"Remote":     2,
		"Hold":       3,
		"UserPB1":    4,
		"UserPB2":    5,
		"CycleStart": 6,
	}
	c.Definitions["SO"] = map[string]int{
		"RemoteLED":  0,
		"CycleStart": 1,
		"Hold":       2,
		"FaultLED":   3,
		"BattAlarm":  4,
		"UserLED1":   5,
		"UserLED2":   6,
		"TPEnabled":  7,
	}

	c.ModifiedAt = time.Now()

	return nil
}

// loadCached attempts to load and validate the local cache against
// the provided configuration.
func (c *cache) loadCached() error {
	buffer, err := os.ReadFile(CACHE_FILE_NAME)
	if err != nil {
		return err
	}

	var fileCache cache
	err = json.Unmarshal(buffer, &fileCache)
	if err != nil {
		return err
	}

	err = fileCache.validate(c.Path, c.FileConfig)
	if err != nil {
		return err
	}

	*c = fileCache

	return nil
}

// populate populates the cache from the local file or the spreadsheet.
func (c *cache) populate() error {
	err := c.loadCached()
	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist), errors.Is(err, ErrCacheInvalid):
			return c.refreshAndSave()
		default:
			return err
		}
	}

	return nil
}

// refreshAndSaves refreshes the cache and saves the cache file.
func (c *cache) refreshAndSave() error {
	err := c.refresh()
	if err != nil {
		return err
	}

	return c.save()
}

// validate validates the cache against the provided spreadsheet path and configuration.
func (c *cache) validate(fpath string, cfg fexcel.FileConfig) error {
	source, err := os.Stat(fpath)
	if err != nil {
		return err
	}

	if c.Path != fpath {
		return ErrCacheInvalid
	}
	if c.ModifiedAt.Before(source.ModTime()) {
		return ErrCacheInvalid
	}
	if c.Version != fexcel.Version {
		return ErrCacheInvalid
	}
	if c.FileConfig != cfg {
		return ErrCacheInvalid
	}

	return nil
}

// newCache returns a cache for the provided path and configuration.
func newCache(fpath string, cfg fexcel.FileConfig) (*cache, error) {
	var c cache
	c.Path = fpath
	c.FileConfig = cfg
	c.Version = fexcel.Version
	c.Definitions = make(map[string]map[string]int)
	c.Constants = make(map[string]string)

	err := c.populate()
	if err != nil {
		return nil, err
	}

	return &c, nil
}
