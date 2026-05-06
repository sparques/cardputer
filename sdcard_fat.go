package cardputer

import (
	"os"
	"strings"

	"tinygo.org/x/tinyfs"
	"tinygo.org/x/tinyfs/fatfs"
)

// SDFS exposes a FAT filesystem layered on top of the built-in microSD slot.
// It is intended for direct filesystem access on TinyGo targets where SD-backed
// FAT is supported.
var SDFS = &sdfs{}

type sdfs struct {
	fs      *fatfs.FATFS
	mounted bool
}

// Init initializes the SD card block device and prepares the FAT filesystem wrapper.
func (s *sdfs) Init() error {
	if err := SDCard.Init(); err != nil {
		return err
	}
	if s.fs == nil {
		s.fs = fatfs.New(SDCard).Configure(&fatfs.Config{SectorSize: fatfs.SectorSize})
	}
	return nil
}

// Filesystem returns the underlying FAT filesystem wrapper after initialization.
func (s *sdfs) Filesystem() (*fatfs.FATFS, error) {
	if err := s.Init(); err != nil {
		return nil, err
	}
	return s.fs, nil
}

// Mount mounts the FAT filesystem on the SD card.
func (s *sdfs) Mount() error {
	if err := s.Init(); err != nil {
		return err
	}
	if s.mounted {
		return nil
	}
	if err := s.fs.Mount(); err != nil {
		return err
	}
	s.mounted = true
	return nil
}

// Unmount unmounts the FAT filesystem if it is currently mounted.
func (s *sdfs) Unmount() error {
	if s.fs == nil || !s.mounted {
		return nil
	}
	if err := s.fs.Unmount(); err != nil {
		return err
	}
	s.mounted = false
	return nil
}

// Format formats the SD card with a FAT filesystem.
func (s *sdfs) Format() error {
	if err := s.Init(); err != nil {
		return err
	}
	return s.fs.Format()
}

// Open opens a file or directory from the mounted FAT filesystem.
func (s *sdfs) Open(path string) (tinyfs.File, error) {
	if err := s.Mount(); err != nil {
		return nil, err
	}
	return s.fs.Open(cleanSDFSPath(path))
}

// OpenFile opens a file on the mounted FAT filesystem using os.O_* flags.
func (s *sdfs) OpenFile(path string, flags int) (tinyfs.File, error) {
	if err := s.Mount(); err != nil {
		return nil, err
	}
	return s.fs.OpenFile(cleanSDFSPath(path), flags)
}

// Mkdir creates a directory on the mounted FAT filesystem.
func (s *sdfs) Mkdir(path string, mode os.FileMode) error {
	if err := s.Mount(); err != nil {
		return err
	}
	return s.fs.Mkdir(cleanSDFSPath(path), mode)
}

// Remove deletes a file or empty directory from the mounted FAT filesystem.
func (s *sdfs) Remove(path string) error {
	if err := s.Mount(); err != nil {
		return err
	}
	return s.fs.Remove(cleanSDFSPath(path))
}

// Rename renames a file or directory on the mounted FAT filesystem.
func (s *sdfs) Rename(oldPath, newPath string) error {
	if err := s.Mount(); err != nil {
		return err
	}
	return s.fs.Rename(cleanSDFSPath(oldPath), cleanSDFSPath(newPath))
}

// Stat returns file metadata from the mounted FAT filesystem.
func (s *sdfs) Stat(path string) (os.FileInfo, error) {
	if err := s.Mount(); err != nil {
		return nil, err
	}
	return s.fs.Stat(cleanSDFSPath(path))
}

// Free returns the number of free bytes reported by the mounted FAT filesystem.
func (s *sdfs) Free() (int64, error) {
	if err := s.Mount(); err != nil {
		return 0, err
	}
	return s.fs.Free()
}

func cleanSDFSPath(path string) string {
	path = strings.TrimSpace(path)
	switch path {
	case "", "/":
		return "/"
	default:
		return strings.TrimPrefix(path, "/")
	}
}
