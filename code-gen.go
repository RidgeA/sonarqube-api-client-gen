package main

import (
	"github.com/pkg/errors"
	"os"
)

const (
	targetDirPermission = 0755
)

func checkOutput(out string) error {
	dir, err := os.Open(out)
	if err != nil {
		return errors.Wrapf(err, "failed to open target dir (%s)", out)
	}
	dirInfo, err := dir.Stat()
	if err != nil {
		return errors.Wrapf(err, "failed to open target dir (%s)", out)
	}
	if !dirInfo.IsDir() {
		return errors.Errorf("target output is not a directory (%s)", dirInfo.Mode().String())
	}
	return nil
}

func getFileWriter(out, name string) (*os.File, error) {
	path := out + "/" + name
	file, err := os.Create(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create file")
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
		return errors.Wrap(err, "cant create destination directory")
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
