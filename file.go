package main

import "os"

func FileExists(path string) bool {
	fileStat, err := os.Stat(path)
	return (err == nil || os.IsExist(err)) && !fileStat.IsDir()
}
