package swap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// ProviderLock represents the state of Terraform Provider locks. This maps
// provider names to provider binary hashes.
type ProviderLock map[string]string

// GetLockFileOrDie reads a lockfile from a path or panics.
func GetLockFileOrDie(path string) ProviderLock {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf(`Error reading lock file "%s": %s\n`, path, err.Error())
	}
	defer f.Close()

	b, _ := ioutil.ReadAll(f)
	lock := ProviderLock{}

	if err := json.Unmarshal([]byte(b), &lock); err != nil {
		log.Fatalf(`Error parsing lock file "%s": %s\n`, path, err.Error())
	}
	return lock
}

// WriteLockFileOrDie writes a lockfile to a path or panics.
func WriteLockFileOrDie(path string, l ProviderLock) {
	b, err := json.Marshal(l)
	if err != nil {
		log.Fatalf(`Error marshaling lock file "%s": %s\n`, path, err.Error())
	}

	if err := ioutil.WriteFile(path, b, os.ModePerm); err != nil {
		log.Fatalf(`Error writing lock file "%s": %s\n`, path, err.Error())
	}
}

// UpdateProvider updates the lock file to use the new given binary for the given
// provider.
func UpdateProvider(name, bin string) error {
	platPath := GetDefaultTerraformPlatformPath()
	binPath := GetTerraformProviderArtifactPath(name, platPath)

	lockPath := filepath.Join(platPath, "lock.json")
	lock := GetLockFileOrDie(lockPath)

	if _, ok := lock[name]; !ok {
		return fmt.Errorf(`Plugin "%s" does not exist in lock file "%s"`, name, lockPath)
	}

	hash := getFileSha256OrDie(bin)
	lock[name] = hash

	if err := copyFile(bin, binPath); err != nil {
		return err
	}

	WriteLockFileOrDie(lockPath, lock)
	color.Cyan("Successfully wrote lockfile...")
	color.Green(`[*] Patched %s to "%s" [%s]!`, name, bin, hash)

	return nil
}
