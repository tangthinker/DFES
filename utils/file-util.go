package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func CreateFileIfNotExist(dir string, fileName string) error {
	_, err := os.Stat(filepath.Join(dir, fileName))
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		chErr := make(chan error)
		go func(ch chan error) {
			log.Println("create file", dir, fileName)
			err := os.MkdirAll(dir, 0700)
			if err != nil {
				log.Printf("mkdir failed![%v]\n", err)
				chErr <- err
			}
			file, err := os.OpenFile(filepath.Join(dir, fileName), os.O_CREATE|os.O_WRONLY, 0700)
			if err != nil {
				log.Printf("make file failed![%v]\n", err)
				chErr <- err
			}
			err = file.Sync()
			if err != nil {
				log.Printf("sync file failed![%v]\n", err)
				chErr <- err
			}
			//fd := file.Fd()
			//if err := syscall.Fsync(int(fd)); err != nil {
			//	log.Println("sync file to disk error:", err)
			//	chErr <- err
			//}
			err = file.Close()
			if err != nil {
				log.Println("close file error:", err)
				chErr <- err
			}
			//err = syscall.Sync()
			//if err != nil {
			//	log.Println("syscall sync error")
			//	chErr <- err
			//}
			chErr <- nil
		}(chErr)
		select {
		case err := <-chErr:
			if err != nil {
				return err
			}
		}
		log.Println("create file success!")
	}
	return nil
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
		time.Sleep(time.Second)
	}
	return err
}
