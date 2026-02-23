package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// Config mapping the TOML structure
type Config struct {
	Storage struct {
		Mode     string `toml:"mode"`
		Path     string `toml:"path"`
		PageSize int    `toml:"page_size"`
	} `toml:"storage"`
	Performance struct {
		SyncOnWrite bool `toml:"sync_on_write"`
	} `toml:"performance"`
}

// Engine Interface so we can swap RAM/Disk implementations easily
type Engine interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Close() error
}

// --- Implementation of the Factory ---

func LoadEngine() (Engine, error) {
	var conf Config
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	fmt.Printf("--- Initializing %s Mode ---\n", conf.Storage.Mode)

	switch conf.Storage.Mode {
	case "disk":
		return NewPersistentEngine(conf.Storage.Path)
	case "ram":
		// Blah
		return NewRamEngine(), nil
	default:
		return nil, fmt.Errorf("unknown storage mode: %s", conf.Storage.Mode)
	}
}

// --- Placeholder for the Engine Logic ---

type PersistentEngine struct {
	file *os.File
}

func NewPersistentEngine(path string) (*PersistentEngine, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	return &PersistentEngine{file: f}, err
}

func (e *PersistentEngine) Set(k, v string) error { 
	fmt.Printf("Writing %s to Disk...\n", k)
	return nil 
}
func (e *PersistentEngine) Get(k string) (string, error) { return "val", nil }
func (e *PersistentEngine) Close() error               { return e.file.Close() }

type RamEngine struct {
	data map[string]string
}

func NewRamEngine() *RamEngine { return &RamEngine{data: make(map[string]string)} }
func (e *RamEngine) Set(k, v string) error { 
	fmt.Printf("Writing %s to RAM...\n", k)
	e.data[k] = v
	return nil 
}
func (e *RamEngine) Get(k string) (string, error) { return e.data[k], nil }
func (e *RamEngine) Close() error                { return nil }

// --- Main Execution ---

func main() {
	db, err := LoadEngine()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Set("user:123", "Gopher")
}
