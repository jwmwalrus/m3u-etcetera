package base

import (
	"path/filepath"

	"github.com/jwmwalrus/bnp/urlstr"
	"golang.org/x/exp/slices"
)

// Playlist extensions
const (
	SupportedFileExtensionMP3  = ".mp3"
	SupportedFileExtensionM4A  = ".m4a"
	SupportedFileExtensionOGG  = ".ogg"
	SupportedFileExtensionFLAC = ".flac"

	SupportedPlaylistExtensionM3U  = ".m3u"
	SupportedPlaylistExtensionM3U8 = ".m3u8"
	SupportedPlaylistExtensionPLS  = ".pls"

	SupportedURISchemeFile = "file"
	SupportedURISchemeHTTP = "http"
)

var (
	// SupportedFileExtensions -
	SupportedFileExtensions = []string{
		SupportedFileExtensionMP3,
		SupportedFileExtensionM4A,
		SupportedFileExtensionOGG,
		SupportedFileExtensionFLAC,
	}

	// SupportedPlaylistExtensions -
	SupportedPlaylistExtensions = []string{
		SupportedPlaylistExtensionM3U,
		SupportedPlaylistExtensionM3U8,
		SupportedPlaylistExtensionPLS,
	}

	// SupportedURISchemes -
	SupportedURISchemes = []string{
		SupportedURISchemeFile,
		SupportedURISchemeHTTP,
	}

	// SupportedMIMETypes -
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

	// IgnoredFileExtensions -
	IgnoredFileExtensions = []string{
		".bmp",
		".db",
		".gif",
		".jpeg",
		".jpg",
		".png",
	}
)

// CheckUnsupportedFiles Returns the unsupported files from a given list
func CheckUnsupportedFiles(files []string) (unsupp []string) {
	for _, f := range files {
		if !IsSupportedFile(f) {
			unsupp = append(unsupp, f)
		}
	}
	return
}

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
	return slices.Contains(SupportedFileExtensions, filepath.Ext(path))
}

// IsSupportedPlaylistURL returns true if the path is supported
func IsSupportedPlaylistURL(s string) bool {
	path, err := urlstr.URLToPath(s)
	if err != nil {
		return false
	}

	return IsSupportedPlaylist(path)
}

// IsSupportedPlaylist returns true if the path is supported
func IsSupportedPlaylist(path string) bool {
	return slices.Contains(SupportedPlaylistExtensions, filepath.Ext(path))
}

// IsIgnoredFile returns true if the path should be ignored
func IsIgnoredFile(path string) bool {
	return slices.Contains(IgnoredFileExtensions, filepath.Ext(path))
}
