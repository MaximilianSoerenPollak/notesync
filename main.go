package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// Where are all the notes we want to track
var notesFolder = "../../Documents/Maximilians Brain"

type fileStruct struct {
	FileName   string `json:"filename,omitempty"`
	ModifiedAt string `json:"modifiedAt,omitempty"`
}

func createFile(fileName string) error {
	_, err := os.OpenFile(fileName, os.O_CREATE, 0644)
	if err != nil {
		log.Println("Creating file: ", fileName)
		log.Println(err)
	}
	return nil
}

func openAndReadFile(fileName string) ([]fileStruct, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	files := []fileStruct{}
	json.Unmarshal([]byte(byteValue), &files)
	return files, nil
}

func checkIfFileExists(fileName string) ([]fileStruct, error) {
	files, err := openAndReadFile(fileName)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Println("File does not exists, creating it now")
			err := createFile(fileName)
			if err != nil {
				log.Println("Error creating File: ", fileName)
				return nil, err
			}
			files, err := checkIfFileExists(fileName)
			if err != nil {
				log.Println("Error checking File: ", fileName)
				return nil, err
			}
			return files, nil
		}
		log.Println("Error reading file: ", fileName)
		log.Println(err)
		return nil, err
	}
	return files, nil
}

func writeFileStructToFile(fs fileStruct, fileName string) error {
	files, err := checkIfFileExists(fileName)
	if err != nil {
		log.Println(err)
		return err
	}
	files = append(files, fs)
	data, err := json.MarshalIndent(files, "", " ")
	if err != nil {
		fmt.Println("Could not marshal data")
		log.Println(err)
		return err
	}
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		fmt.Println("Could not write data to file")
		log.Println(err)
		return err
	}
	return nil
}

func findFiles(folder string) error {
	err := filepath.Walk(folder,
		func(_ string, info os.FileInfo, err error) error {
			// We don't need to track obsidian config files
			if info.IsDir() && strings.HasPrefix(info.Name(), ".obsidian") {
				return filepath.SkipDir
			}
			newFile := fileStruct{
				FileName:   info.Name(),
				ModifiedAt: info.ModTime().String(),
			}
			writeFileStructToFile(newFile, "temp.json")
			if err != nil {
				log.Println(err)
				return err
			}
			return nil
		})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func checkIfFileIsEmpty(fileName string) (bool, error) {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		log.Println("Could not read filesize of file: ", fileName)
		log.Println(err)
		return false, err
	}
	if fileInfo.Size() == 0 {
		return true, nil
	} else {
		return false, nil
	}

}
func replaceTrackedNotesFile() error {
	err := os.Remove("trackednotes.json")
	if err != nil {
		log.Println("Could not remove trackednotes.json")
		log.Println(err)
		return err
	}
	os.Rename("temp.json", "trackednotes.json")
	log.Println("Renamed temp.json -> trackednotes.json")
	return nil
}

func compareFiles() error {
	// Temp.json should already exists due to the called functions before
	temp_files, err := openAndReadFile("temp.json")
	if err != nil {
		log.Println("could not open File: temp.json")
		return err
	}
	// Perm.json could be not created yet. Therefore we ned to check
	tracked_notes, err := checkIfFileExists("trackednotes.json")
	if err != nil {
		log.Println("could not open File: trackednotes.json")
		return err
	}
	log.Println("Checking if trackednotes.json is empty")
	fileEmpty, err := checkIfFileIsEmpty("trackednotes.json")
	if err != nil {
		log.Println("could not check File for size: trackednotes.json")
		return err
	}
	if fileEmpty || !reflect.DeepEqual(temp_files, tracked_notes) {
		err := replaceTrackedNotesFile()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
	log.Println("Files are equal. Removing unneeded temp.json")
	err = os.Remove("temp.json")
	if err != nil {
		log.Println("Something went wrong deleting temp.json .")
		log.Fatal(err)
	}
	return nil
}

func main() {
	LOG_FILE := "notesync.log"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Printf("=========EXECUTING============")
	err = createFile("temp.json")
	if err != nil {
		log.Println("Something went wrong creating temp.json .")
		log.Fatal(err)
	}
	err = findFiles(notesFolder)
	if err != nil {
		log.Println("Something went wrong finding Files.")
		log.Fatal(err)
	}
	err = compareFiles()
	if err != nil {
		log.Println("Something went wrong comparing Files.")
		log.Fatal(err)
	}
	log.Printf("=========STOPPING EXECUTION =========")
}
