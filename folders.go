// folders project folders.go
package folders

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

//Возвращает текущую директорию выполнения скрипта
func GetExDir() string {
	return path.Dir(os.Args[0])
}

//Проверяет наличие папки, создаём и возвращает полный путь к ней, если её нет
func CheckReportDir(dir_name string) string {
	PATH := GetExDir() //Директория выполения
	REPORT_PATH := path.Join(PATH, dir_name)
	if _, err := os.Stat(REPORT_PATH); os.IsNotExist(err) {
		os.Mkdir(REPORT_PATH, os.ModePerm)
		fmt.Println(fmt.Sprintf("Директория %s создана", REPORT_PATH))
	}
	return REPORT_PATH
}

//elem - список директорий раделённый ","
//Возвращает массив, содержащий полные пути к файлам
func GetFilesFromDir(elem ...string) []string {
	var (
		arr  []string
		PATH string
	)
	for _, item := range elem {
		PATH = path.Join(PATH, item)
	}
	files, err := ioutil.ReadDir(PATH)
	if err == nil {
		for _, item := range files {
			if item.IsDir() == false {
				arr = append(arr, path.Join(PATH, item.Name()))
			}
		}
	}
	return arr
}

func UnZIPFile(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(target, os.ModePerm); err != nil {
		return err
	}
	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		dir, _ := filepath.Split(path)
		CheckReportDir(dir)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}

func UnZIPReader(buf *bytes.Reader, size int64, target string) (string, error) {
	reader, _ := zip.NewReader(buf, size)
	var oval_path string
	for _, file := range reader.File {
		oval_path = path.Join(target, file.Name)
		dir, _ := filepath.Split(oval_path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, os.ModePerm)
			fmt.Println(fmt.Sprintf("Директория %s создана", dir))
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(oval_path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return "", err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(oval_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return "", err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return "", err
		}
	}
	return oval_path, nil
}
