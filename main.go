package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

func inputString() string {
	reader := bufio.NewReader(os.Stdin)
	//read until a newline character
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occurred while reading input:", err)
		return ""
	}

	// Remove the trailing newline character (\n or \r\n on Windows)
	input = strings.TrimSpace(input)

	return input
}

func splitExt(path string) (string, string) {
	// 1. Get the extension (includes the leading dot, e.g., ".txt")
	ext := filepath.Ext(path)

	// 2. Trim the extension from the original path to get the base name
	base := strings.TrimSuffix(path, ext)

	return base, ext
}

func getFilesOnly(dir string) ([]string, error) {
	// ReadDir returns entries sorted by filename
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, entry := range entries {
		// Filter out subdirectories
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())

			// Use filepath.Join(dir, entry.Name()) instead if you need full paths
		}
	}

	slices.Sort(fileNames)
	return fileNames, nil
}

func inputFileDirectory() string {
	fmt.Print("Enter the file directory: ")
	fileDir := inputString()
	return fileDir
}

func copyFile(src, dst string) error {
	// 1. Open the source file for reading
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	// 2. Create the destination file (overwrites if it already exists)
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	// 3. Stream the data from source to destination
	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	// 4. Commit the file contents to stable storage
	return destination.Sync()
}

func renameFilesInOrder(inpPath string) {
	fmt.Println("\n\nRename files")
	fmt.Println("Enter root name to use: ")
	rootName := inputString()

	fmt.Println("Beginning seq: ")
	begSeqStr := inputString()
	begSeq, err := strconv.Atoi(begSeqStr)
	if err != nil {
		log.Fatalf("Not a number: %s", err)
	}

	fmt.Println("Number of digits: ")
	numDigitsStr := inputString()
	numDigits, err := strconv.Atoi(numDigitsStr)
	if err != nil {
		log.Fatalf("Not a number: %s", err)
	}

	fmt.Println("\nWill do:")
	fmt.Printf("input path       : %s\n", inpPath)
	fmt.Printf("root name        : %s\n", rootName)
	fmt.Printf("begining sequence: %d\n", begSeq)
	fmt.Printf("number of digits : %d\n", numDigits)
	fmt.Println("\nEnter (C)ontinue or (Q)uit")
	commStr := inputString()
	commStr = strings.ToUpper(commStr)
	switch commStr {
	case "C":
		files, err := getFilesOnly(inpPath)
		if err != nil {
			log.Fatalf("can not get files: %s", err)
		}
		for index, fileName := range files {
			fileBase, fileExt := splitExt(fileName)
			currNum := begSeq + index
			numDigitsFstring := strconv.Itoa(numDigits)
			formatStr := "%0" + numDigitsFstring + "d"
			currNumStr := fmt.Sprintf(formatStr, currNum)
			fileNameStr := rootName + currNumStr + fileExt
			fmt.Printf("orig=%s, name=%s, new=%s\n", fileName, fileBase, fileNameStr)
			copyFile(fileName, fileNameStr)
		}
		return
	case "Q":
		return
	}
}

func main() {
	fmt.Println("File utility")
	fmt.Println("\nMenu")
	fmt.Println("1 = rename files")
	fmt.Println("0 = quit")
	commStr := inputString()
	commStr = strings.ToUpper(commStr)
	switch commStr {
	case "1":
		inpFileDir := inputFileDirectory()
		renameFilesInOrder(inpFileDir)
		return
	case "0":
		return
	}
}
