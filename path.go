package jazz

import (
	"fmt"
	"os"
)

// CreateDirIfNotExist takes a path, check the state of that path
// if the path exists do nothing;
// else create a dir of that path
func (j *Jazz) CreateDirIfNotExist(dir string) error {
	const mode = 0755
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.Mkdir(dir, mode)
		if err != nil {
			ErrPathFail.Path = dir
			ErrPathFail.Cause = err
			return ErrPathFail
		}
	}
	return nil
}

func (j *Jazz) CreateDotEnvIfNotExists(path string) error {

	err := j.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		ErrEnvFileFail.Path = path + "/.env"
		ErrEnvFileFail.Cause = err
		return ErrEnvFileFail
	}
	return nil
}

func (j *Jazz) CreateFileIfNotExists(file string) error {
	_, errF := os.Stat(file)
	if os.IsNotExist(errF) {
		var f, err = os.Create(file)
		if err != nil {
			ErrFileFail.Path = file
			ErrFileFail.Cause = err
			return ErrFileFail
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
	}
	return nil
}
