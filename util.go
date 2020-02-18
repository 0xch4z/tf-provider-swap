package swap

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	// ErrNotTerraformWorkspace denotes that the current working directory does
	// not have any context of a Terraform workspace.
	ErrNotTerraformWorkspace = "This is not a Terraform workspace directory"
)

// getFileSha256 generates a Sha256 hash for a file given a path.
func getFileSha256OrDie(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf(`Error getting opening file "%s"\n`, err.Error())
	}
	defer f.Close()

	hasher := sha256.New()
	io.Copy(hasher, f)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// getDirectories accumulates all the directories in a given path.
func getDirectories(path string) ([]string, error) {
	dirs := []string{}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return dirs, err
	}
	for _, file := range files {
		if file.IsDir() {
			dirPath := filepath.Join(path, file.Name())
			dirs = append(dirs, dirPath)
		}
	}
	return dirs, nil
}

// copyFile copies a given file to a new location.
func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	newFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, srcFile)
	return err
}

// AssertInTerraformWorkspace panics if the current directory does not contain a Terraform
// workspace.
func AssertInTerraformWorkspace() {
	if stat, err := os.Stat(TerraformWorkspaceDirectory); err != nil && err != os.ErrNotExist {
		log.Fatalf("Error asserting Terraform workspace directory: %s\n", err.Error())
	} else if err == os.ErrNotExist || !stat.IsDir() {
		log.Fatal(ErrNotTerraformWorkspace)
	}
}
