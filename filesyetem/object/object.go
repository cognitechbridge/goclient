package object

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Object struct {
	resolver  func(id string, writer io.Writer) (err error)
	readPath  string
	writePath string
}

func New(path string, resolver func(id string, writer io.Writer) (err error)) Object {
	return Object{
		readPath:  path,
		resolver:  resolver,
		writePath: filepath.Join(path, "Write"),
	}
}

func (o *Object) Move(oldId string, newId string) (err error) {
	//Move to Write cache path
	oldPath := filepath.Join(o.readPath, oldId)
	newPath := filepath.Join(o.writePath, newId)
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return
	}
	//Create link
	err = o.createWriteLink(newId)
	if err != nil {
		return
	}
	return nil
}

func (o *Object) Write(id string, buff []byte, ofst int64) (n int, err error) {
	p := filepath.Join(o.writePath, id)
	file, err := os.OpenFile(p, os.O_RDWR, 0666)
	if err != nil {
		return 0, fmt.Errorf("file is not in write cache: %v", err)
	}
	defer file.Close()
	if err != nil {
		return 0, err
	}
	n, err = file.WriteAt(buff, ofst)
	return
}

func (o *Object) Truncate(id string, size int64) (err error) {
	p := filepath.Join(o.writePath, id)
	file, err := os.OpenFile(p, os.O_RDWR, 0666)
	defer file.Close()
	err = file.Truncate(size)
	if err != nil {
		return err
	}
	return nil
}

func (o *Object) Create(id string) (err error) {
	objWritePath := filepath.Join(o.writePath, id)
	objFile, err := os.Create(objWritePath)
	objFile.Close()
	if err != nil {
		return
	}
	err = o.createWriteLink(id)
	if err != nil {
		return
	}
	return nil
}

func (o *Object) Read(id string, buff []byte, ofst int64) (n int, err error) {
	p := filepath.Join(o.readPath, id)

	if _, err := os.Stat(p); os.IsNotExist(err) {
		err = o.resolverFile(id)
		if err != nil {
			return 0, err
		}
	}

	file, err := os.OpenFile(p, os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		return 0, err
	}
	n, err = file.ReadAt(buff, ofst)
	return
}

func (o *Object) Flush(id string) (err error) {
	p := filepath.Join(o.writePath, id)
	err = os.Remove(p)
	return
}

func (o *Object) AsFile(id string) (file *os.File, err error) {
	p := filepath.Join(o.readPath, id)
	file, err = os.OpenFile(p, os.O_RDONLY, 0666)
	return file, err
}

func (o *Object) createWriteLink(id string) (err error) {
	objWritePath := filepath.Join(o.writePath, id)
	objFilePath := filepath.Join(o.readPath, id)
	err = os.Link(objWritePath, objFilePath)
	return err
}

func (o *Object) resolverFile(id string) (err error) {
	file, _ := os.Create(filepath.Join(o.readPath, id))
	defer file.Close()
	err = o.resolver(id, file)
	return
}