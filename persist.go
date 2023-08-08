package persist

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
)

var cacheDir string

type PersisentVar[T any] struct {
	V    T      // The actual variable
	Name string // The name of the variable
}

func getCachePath(name string) string {
	return path.Join(cacheDir, name+".gob")
}

// InitWithDir initialises variable persistence with a given directory
func InitWithDir(dir string) bool {
	if dir != "" {
		cacheDir = dir
	}
	return true
}

// Init initialises variable persistence with the default directory (./cache/)
func Init() bool {
	return InitWithDir("./cache/")
}

// Clean removes all cache from previous runs and ensures the directory exists
func Clean() {
	os.RemoveAll(cacheDir)
	os.MkdirAll(cacheDir, 0777)
}

// Load creates a persistent variable which is stored under the name. The name should be unique
func Load[T any](name string) *PersisentVar[T] {
	path := getCachePath(name)
	var v T
	if f, err := os.Open(path); err == nil {
		defer f.Close()
		dec := gob.NewDecoder(f)
		err := dec.Decode(&v)
		if err != nil {
			fmt.Printf("Failed to read cache, so value has been reset. Have you remembered to write 'var _ = persist.Clean()' at the start of your file?. Error was %v", err)
		}
	}
	return &PersisentVar[T]{v, name} // This will just return the empty value if it fails or the cache is not there
}

// Save saves the state of the variable so that other cells can use it
func (v *PersisentVar[T]) Save() {
	path := getCachePath(v.Name)
	if f, err := os.Create(path); err == nil {
		defer f.Close()
		enc := gob.NewEncoder(f)
		enc.Encode(v.V)
	} else {
		panic(fmt.Sprintf("Failed to save %v: %v", v.Name, err))
	}
}
