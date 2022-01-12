package config

// GTK Gtk-related config
type GTK struct {
	Playback struct {
		CoverFilenames []string `json:"coverFilenames"`
	} `json:"playback"`
}

// SetDefaults provides default settings
func (g *GTK) SetDefaults() {
	if len(g.Playback.CoverFilenames) == 0 {
		g.Playback.CoverFilenames = []string{
			"cover",
			"album",
			"artwork",
			"image",
			"front",
		}

	}
}
