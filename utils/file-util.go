package utils

import (
	"log"
	"os"
	"path/filepath"
)

func CreateFileIfNotExist(dir string, fileName string) error {
	_, err := os.Stat(filepath.Join(dir, fileName))
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			log.Printf("mkdir failed![%v]\n", err)
			return err
		}
		file, err := os.OpenFile(filepath.Join(dir, fileName), os.O_CREATE, 0700)
		if err != nil {
			log.Printf("make file failed![%v]\n", err)
			return err
		}
		err = file.Sync()
		err = file.Close()
		if err != nil {
			log.Printf("close file failed![%v]\n", err)
			return err
		}
	}
	return err
}

func CreateDirIfNotExist(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			log.Printf("mkdir failed![%v]\n", err)
			return err
		}
	}
	return err
}
