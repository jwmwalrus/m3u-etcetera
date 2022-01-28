package base

import (
	"path/filepath"

	"github.com/jwmwalrus/bnp/slice"
	"github.com/jwmwalrus/bnp/urlstr"
)

const (
	// SupportedFileExtensionMP3 supported mp3
	SupportedFileExtensionMP3 = ".mp3"

	// SupportedFileExtensionM4A supported m4a
	SupportedFileExtensionM4A = ".m4a"

	// SupportedFileExtensionOGG supported ogg
	SupportedFileExtensionOGG = ".ogg"

	// SupportedFileExtensionFLAC supported flac
	SupportedFileExtensionFLAC = ".flac"
)

var (
	// SupportedFileExtensions supported file extensons
	SupportedFileExtensions = []string{
		SupportedFileExtensionMP3,
		SupportedFileExtensionM4A,
		SupportedFileExtensionOGG,
		SupportedFileExtensionFLAC,
	}

	SupportedURISchemes = []string{
		"file",
		"http",
	}

	SupportedMIMETypes = []string{
		"audio/x-mp3",
		"application/x-id3",
		"audio/mpeg",
		"audio/x-mpeg",
		"audio/x-mpeg-3",
		"audio/mpeg3",
		"audio/mp3",
		"audio/x-m4a",
		"audio/mpc",
		"audio/x-mpc",
		"audio/mp",
		"audio/x-mp",
		"application/ogg",
		"application/x-ogg",
		"audio/vorbis",
		"audio/x-vorbis",
		"audio/ogg",
		"audio/x-ogg",
		"audio/x-flac",
		"application/x-flac",
		"audio/flac",
	}

	// IgnoredFileExtensions supported file extensons
	IgnoredFileExtensions = []string{
		".bmp",
		".db",
		".gif",
		".jpeg",
		".jpg",
		".png",
	}
)

// IsSupportedURL returns true if the path is supported
func IsSupportedURL(s string) bool {
	path, err := urlstr.URLToPath(s)
	if err != nil {
		return false
	}

	return IsSupportedFile(path)
}

// IsSupportedFile returns true if the path is supported
func IsSupportedFile(path string) bool {
	return slice.Contains(SupportedFileExtensions, filepath.Ext(path))
}

// IsIgnoredFile returns true if the path should be ignored
func IsIgnoredFile(path string) bool {
	return slice.Contains(IgnoredFileExtensions, filepath.Ext(path))
}
