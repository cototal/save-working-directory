package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"
)

type WorkingDirectory struct {
	Shortcut string `json:"shortcut"`
	Path     string `json:"path"`
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dataDir := filepath.Join(homeDir, ".tmp")
	if _, err := os.Stat(dataDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dataDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	dataFile := filepath.Join(dataDir, "save-working-directory.json")
	if _, err := os.Stat(dataFile); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(dataFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		_, err = file.WriteString("[]")
		if err != nil {
			log.Fatal(err)
		}
	}
	jsonBytes, err := os.ReadFile(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	var workingDirs []WorkingDirectory
	if err := json.Unmarshal(jsonBytes, &workingDirs); err != nil {
		log.Fatal(err)
	}
	sort.Slice(workingDirs, func(idx, jdx int) bool {
		return workingDirs[idx].Shortcut < workingDirs[jdx].Shortcut
	})

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage...")
		os.Exit(0)
	}

	switch args[0] {
	case "--list":
		writer := tabwriter.NewWriter(os.Stdout, 1, 1, 5, ' ', 0)
		fmt.Fprintln(writer, "SHORTCUT\tPATH")
		for _, wd := range workingDirs {
			fmt.Fprintf(writer, "%s\t%s\n", wd.Shortcut, wd.Path)
		}
		writer.Flush()
	case "-s":
		workDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		shortcut := "default"
		if len(args) > 1 {
			shortcut = args[1]
		}
		thisDir := WorkingDirectory{
			Shortcut: shortcut,
			Path:     workDir,
		}
		isReplaced := false
		for idx, wd := range workingDirs {
			if wd.Shortcut == shortcut {
				isReplaced = true
				workingDirs[idx].Path = workDir
				break
			}
		}
		if !isReplaced {
			workingDirs = append(workingDirs, thisDir)
		}
		jsonBytes, err = json.Marshal(workingDirs)
		if err != nil {
			log.Fatal(err)
		}
		os.WriteFile(dataFile, jsonBytes, os.ModePerm)
		fmt.Println("Working directory saved as:", shortcut)
	case "-l":
		shortcut := "default"
		if len(args) > 1 {
			shortcut = args[1]
		}
		path := ""
		for _, wd := range workingDirs {
			if wd.Shortcut == shortcut {
				path = wd.Path
				break
			}
		}
		if len(path) == 0 {
			log.Fatal("No path found for shortcut:", shortcut)
		}
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}
		fmt.Println(path)
	case "-d":
		if len(args) == 1 {
			log.Fatal("Please provide the shortcut you wish to delete")
		}
		shortcut := args[1]
		wdIdx := -1
		for idx, wd := range workingDirs {
			if wd.Shortcut == shortcut {
				wdIdx = idx
				break
			}
		}
		if wdIdx < 0 {
			log.Fatal("No path found for shortcut:", shortcut)
		}
		workingDirs = append(workingDirs[:wdIdx], workingDirs[wdIdx+1:]...)
		jsonBytes, err = json.Marshal(workingDirs)
		if err != nil {
			log.Fatal(err)
		}
		os.WriteFile(dataFile, jsonBytes, os.ModePerm)
		fmt.Println("Shortcut deleted:", shortcut)

	default:
		log.Fatal("Unrecognized option:", args[0])
	}
}
