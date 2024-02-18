package repositories

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

type LinkRepository struct {
	rootPath string
	file     *os.File
}

const FsSizeOffset = 0
const FsIdOffset = 8

func NewLinkRepository(rootPath string) *LinkRepository {
	return &LinkRepository{
		file:     nil,
		rootPath: rootPath,
	}
}

// Create link file
func (c *LinkRepository) Create(path string, key string, size int64) error {
	absPath := filepath.Join(c.rootPath, path)
	err := os.MkdirAll(filepath.Dir(absPath), os.ModePerm)
	if err != nil {
		return err
	}
	err = c.open(path)
	if err != nil {
		return err
	}
	err = c.WriteSize(path, size)
	if err != nil {
		return err
	}
	err = c.WriteId(path, key)
	if err != nil {
		return err
	}
	return nil
}

func (c *LinkRepository) WriteSize(path string, size int64) (err error) {
	defer c.openClose(path)()
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, size)
	_, err = c.file.WriteAt(buf.Bytes(), FsSizeOffset)
	return
}

func (c *LinkRepository) ReadSize(path string) (size int64, err error) {
	defer c.openClose(path)()
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 8)
	_, err = c.file.ReadAt(buf, FsSizeOffset)
	err = binary.Read(bytes.NewReader(buf), binary.BigEndian, &size)
	if err != nil {
		return
	}
	return
}

func (c *LinkRepository) WriteId(path string, key string) (err error) {
	defer c.openClose(path)()
	if err != nil {
		return err
	}
	_, err = c.file.WriteAt([]byte(key), FsIdOffset)
	return
}

func (c *LinkRepository) ReadId(path string) (key string, err error) {
	defer c.openClose(path)()
	if err != nil {
		return "", err
	}
	id := make([]byte, 36)
	_, _ = c.file.ReadAt(id, FsIdOffset)
	return string(id), nil
}

func (c *LinkRepository) openClose(path string) (fu func()) {
	_ = c.open(path)
	return func() {
		_ = c.close()
	}
}

func (c *LinkRepository) open(path string) (err error) {
	if c.file != nil {
		return nil
	}
	p := filepath.Join(c.rootPath, path)
	c.file, err = os.OpenFile(p, os.O_RDWR|os.O_CREATE, 0666)
	return
}

func (c *LinkRepository) close() (err error) {
	err = c.file.Close()
	c.file = nil
	return
}

func (c *LinkRepository) ListIdsByRegex(regex string) ([]string, error) {
	var matchedIds []string
	list, err := c.listFilesByRegex(regex)
	if err != nil {
		return nil, err
	}
	for _, file := range list {
		id, err := c.ReadId(file)
		if err != nil {
			return nil, err
		}
		matchedIds = append(matchedIds, id)
	}
	return matchedIds, nil
}

// listFilesByRegex lists all files in the specified directory that match the given regex pattern.
// Returns a slice of strings containing the names of matching files.
func (c *LinkRepository) listFilesByRegex(pattern string) ([]string, error) {
	var matchedFiles []string

	dirPath := c.rootPath

	// Walk the directory tree
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Return the error to stop the walk
		}
		if info.IsDir() {
			return nil // Skip directories
		}

		// Generate the relative path
		relativePath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err // Return the error to stop the walk
		}

		// Check if the file matches the pattern
		match, err := filepath.Match(pattern, relativePath)
		if err != nil {
			return err // Return the error to stop the walk
		}
		if match {
			matchedFiles = append(matchedFiles, relativePath) // Add matching file to the slice
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the path %q: %w", dirPath, err)
	}

	return matchedFiles, nil // Return the slice of matched file paths
}
