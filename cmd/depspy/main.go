package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path"
	"strings"
)

func main() {
	cmd := exec.Command("go", "list", "-json", "all")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println("Error:", err)

	}

	var pkg build.Package

	decoder := json.NewDecoder(&out)
	for decoder.More() {
		err := decoder.Decode(&pkg)
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		fmt.Println("Package:", pkg.ImportPath)
		fmt.Println("Source files:", len(pkg.GoFiles))
		fmt.Println("Dir:", pkg.Dir)
		totalLoc := 0
		for _, f := range pkg.GoFiles {
			loc, _ := loc(path.Join(pkg.Dir, f))
			totalLoc += loc
		}
		fmt.Println("LoC:", totalLoc)
		fmt.Println()
	}

}

func loc(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer func() { _ = file.Close() }()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	loc := 0

	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			continue
		}
		loc++
	}

	return loc, scanner.Err()
}
