package config

import "github.com/nightlyone/lockfile"

// Config  implements the rtc.Config interface.
type Config struct {
	FirstRun bool   `json:"firstRun"`
	Server   Server `json:"server"`
	Task     Task   `json:"task"`
	GTK      GTK    `json:"gtk"`
}

func (c *Config) GetFirstRun() bool { return c.FirstRun }

func (c *Config) SetFirstRun(v bool) { c.FirstRun = v }

func (c *Config) SetLockFile(_ lockfile.Lockfile) {}

func (c *Config) SetDefaults() {
	c.Server.SetDefaults()
	c.Task.SetDefaults()
	c.GTK.SetDefaults()
}
