package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Render will apply templating on all files in given directory
func Render(inPath, outPath string, valueFiles []string) {
	files := []string{}
	err := error(nil)

	logger := log.New(os.Stderr, "Error: ", 0)

	isStdout := outPath == "stdout" || outPath == ""

	if !isStdout {
		if result, err := validDirectory(outPath); err != nil {
			logger.Fatal(err)
			return
		} else if !result {
			logger.Fatal(outPath + " is not a valid directory")
			return
		}
	}

	if files, err = getFiles(inPath); err != nil {
		logger.Fatal(err)
		return
	}

	var values Values

	for _, valueFile := range valueFiles {
		var currentValues Values
		if currentValues, err = ReadValuesFile(valueFile); err != nil {
			logger.Fatal(err)
			return
		}
		values = CoalesceValues(currentValues, values)
	}

	for _, file := range files {
		var result string
		if text, err := ioutil.ReadFile(file); err != nil {
			logger.Fatal(err)
		} else if result, err = ExecuteTemplate(string(text), map[string]interface{}{"Values": values}); err != nil {
			logger.Fatal(err)
		}

		if isStdout {
			fmt.Println("---")
			fmt.Println(result)
		} else {
			writeFile(outPath, file, result)
		}
	}

	return
}

func validDirectory(path string) (result bool, err error) {
	path = filepath.FromSlash(path)
	//	strings.TrimRight(path, string(os.PathSeparator))
	fi, err := os.Stat(path)
	if err != nil {
		return
	}
	result = fi.Mode().IsDir()
	return
}

func writeFile(path, file, text string) (err error) {
	path = filepath.FromSlash(path)
	strings.TrimRight(path, string(os.PathSeparator))
	_, filename := filepath.Split(file)
	fileNameWithPath := path + string(os.PathSeparator) + filename
	err = ioutil.WriteFile(fileNameWithPath, []byte(text), 0644)
	fmt.Println("Written file " + fileNameWithPath)
	return
}

func getFiles(path string) (files []string, err error) {
	path = filepath.FromSlash(path)
	strings.TrimRight(path, string(os.PathSeparator))
	fi, err := os.Stat(path)
	if err != nil {
		return
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		if fileInfos, err := ioutil.ReadDir(path); err != nil {
			for _, fileInfo := range fileInfos {
				files = append(files, path+string(os.PathSeparator)+fileInfo.Name())
			}
		}
	case mode.IsRegular():
		// do file stuff
		files = append(files, path)
	}
	return
}
