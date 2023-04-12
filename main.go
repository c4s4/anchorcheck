package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var regexpHref = regexp.MustCompile(`href="#([^"]+)"`)

func processFile(file string) (bool, error) {
	fmt.Printf("Processing file '%s'...\n", file)
	ok := true
	bytes, err := os.ReadFile(file)
	if err != nil {
		return false, fmt.Errorf("reading file '%s': %v", file, err)
	}
	page := string(bytes)
	matches := regexpHref.FindAllStringSubmatch(page, -1)
	for _, match := range matches {
		anchor, err := url.QueryUnescape(match[1])
		if err != nil {
			return false, fmt.Errorf("unescaping anchor: %v", err)
		}
		if !strings.Contains(page, `id="`+anchor+`"`) {
			fmt.Fprintf(os.Stderr, "In file '%s' anchor link %s not found\n", file, anchor)
			ok = false
		}
	}
	return ok, nil
}

func listFilesInDir(dir string) ([]string, error) {
	var files []string
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			files = append(files, filepath.Join(dir, fileInfo.Name()))
		}
	}
	return files, nil
}

func processDir(dir string) (bool, error) {
	fmt.Printf("Processing directory '%s'...\n", dir)
	ok := true
	files, err := listFilesInDir(dir)
	if err != nil {
		return false, err
	}
	for _, file := range files {
		okFile, err := processFile(file)
		if err != nil {
			return false, err
		}
		ok = ok && okFile
	}
	return ok, nil
}

func process(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
		// handle the error and return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}
	if fileInfo.IsDir() {
		return processDir(path)
	} else {
		return processFile(path)
	}
}

func main() {
	okAll := true
	for _, path := range os.Args[1:] {
		ok, err := process(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR processing '%s': %v\n", path, err)
		}
		okAll = ok && okAll
	}
	if !okAll {
		os.Exit(1)
	}
}
