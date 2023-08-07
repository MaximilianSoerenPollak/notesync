package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Step 1. We have to get the folder and check for updates on it?

var notesFolder = "../../Documents/Maximilians Brain"

// Step 2. Define file structure for the json we want to save

type fileStruct struct {
	FileName   string `json:"filename,omitempty"`
	ModifiedAt string `json:"modifiedAt,omitempty"`
}

func createFile(fileName string) error {
	_, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
func openFile(fileName string) ([]fileStruct, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	files := []fileStruct{}
	json.Unmarshal([]byte(byteValue), &files)
	return files, nil
}

func writeFileStructToFile(fs fileStruct, fileName string) error {
	// fmt.Println("This is a file struct", fs)
	files, err := openFile(fileName)
	if err != nil {
		if os.IsNotExist(err){
			fmt.Println("File does not exists, creating it now")
			err := createFile(fileName)
			if err != nil {
				fmt.Println("Error creating File")
				return err 
			}
			writeFileStructToFile(fs, fileName)
		}
		fmt.Println(err)
		return err 
	}
	if err != nil {
		fmt.Println("File potentially not yet created, creating File")
		err = createFile(fileName)
		if err != nil {
			fmt.Println(err)
		}
	}
	files = append(files, fs)
	data, err := json.MarshalIndent(files, "", " ")
	if err != nil {
		fmt.Println("Could not marshal data")
		fmt.Println(err)
		return err
	}
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		fmt.Println("Could write data to file")
		fmt.Println(err)
		return err
	}
	return nil
}

func findFiles(folder string) {
	err := filepath.Walk(folder,
		func(_ string, info os.FileInfo, err error) error {
			// We don't need to track obsidian config files
			if info.IsDir() && strings.HasPrefix(info.Name(), ".obsidian") {
				return filepath.SkipDir
			}
			fmt.Println("This is info.name", info.Name())
			fmt.Println("This is info.Modtime", info.ModTime())
			newFile := fileStruct{
				FileName:   info.Name(),
				ModifiedAt: info.ModTime().String(),
			}
			writeFileStructToFile(newFile, "temp.json")
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func compareFiles(){}


func main() {
	findFiles(notesFolder)
}
