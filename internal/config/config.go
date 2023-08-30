package config

import "github.com/nightlyone/lockfile"

// Config Application's configuration.
type Config struct {
	FirstRun bool   `json:"firstRun"`
	Server   Server `json:"server"`
	Task     Task   `json:"task"`
	GTK      GTK    `json:"gtk"`
}

// GetFirstRun implements rtc.Config.
func (c *Config) GetFirstRun() bool { return c.FirstRun }

// SetFirstRun implements rtc.Config.
func (c *Config) SetFirstRun(v bool) { c.FirstRun = v }

// SetLockFile implements rtc.Config.
func (c *Config) SetLockFile(_ lockfile.Lockfile) {}

// SetDefaults implements rtc.Config.
func (c *Config) SetDefaults() {
	c.Server.SetDefaults()
	c.Task.SetDefaults()
	c.GTK.SetDefaults()
}
