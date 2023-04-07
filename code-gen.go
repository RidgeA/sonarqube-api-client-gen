package main

import (
	"fmt"
	"os"
)

const (
	targetDirPermission = 0755
)

func checkOutput(out string) error {
	dir, err := os.Open(out)
	if err != nil {
		return fmt.Errorf("failed to open target dir (%s)：%w", out, err)
	}
	dirInfo, err := dir.Stat()
	if err != nil {
		return fmt.Errorf("failed to open target dir (%s)：%w", out, err)
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("target output is not a directory (%s)：%w", dirInfo.Mode().String(), err)
	}
	return nil
}

func getFileWriter(out, name string) (*os.File, error) {
	path := out + "/" + name
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file：%w", err)
	}
	return file, nil
}

func generateCode(def *apiDefinition, out string) error {

	if err := checkOutput(out); err != nil {
		return err
	}

	path := out + "/" + def.PackageName

	err := os.MkdirAll(path, targetDirPermission)
	if err != nil {
		return fmt.Errorf("cant create destination directory：%w", err)
	}

	//create files for service
	for _, service := range def.WebServices {
		if err := generateService(path, service); err != nil {
			return err
		}
	}

	// create main client file
	file, err := getFileWriter(path, clientFileName)
	defer file.Close()
	if err != nil {
		return err
	}
	err = renderClient(file, def)
	if err != nil {
		return err
	}

	return nil
}
func generateService(path string, service *webService) error {
	file, err := getFileWriter(path, service.fileName())
	if err != nil {
		return err
	}
	defer file.Close()
	err = renderService(file, service)
	if err != nil {
		return err
	}
	return nil
}
