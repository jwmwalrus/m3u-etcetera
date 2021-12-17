package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/nightlyone/lockfile"
)

// Config Application's configuration
type Config struct {
	FirstRun bool `json:"firstRun"`
	Server   Server
	Query    Query
}

// Load Loads application's configuration
func (c *Config) Load(path, lockFile string) (err error) {
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		c.FirstRun = true
		if err = c.Save(path, lockFile); err != nil {
			return
		}
	}

	// var jsonFile *os.File
	f, err := os.Open(path)
	onerror.Panic(err)
	defer f.Close()

	bv, _ := ioutil.ReadAll(f)

	json.Unmarshal(bv, c)

	return
}

// Save Saves application's configuration
func (c *Config) Save(path, lockFile string) (err error) {
	c.SetDefaults()

	var lock lockfile.Lockfile
	lock, err = lockfile.New(lockFile)
	if err != nil {
		return
	}

	if err = lock.TryLock(); err != nil {
		return
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			fmt.Printf("Cannot unlock %q, reason: %v\n", lock, err)
		}
	}()

	var file []byte
	file, err = json.MarshalIndent(c, "", " ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path, file, 0644)

	return
}

// SetDefaults sets configuration defaults
func (c *Config) SetDefaults() {
	c.Server.SetDefaults()
	c.Query.SetDefaults()
}
