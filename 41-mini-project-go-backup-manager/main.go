package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// Config structure for backup jobs
type Config struct {
	BackupDir     string         `json:"backup_dir"`
	Sources       []SourceConfig `json:"sources"`
	RotationLimit int            `json:"rotation_limit"`
}

type SourceConfig struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// zipSource compresses a source directory (or file) into a zip file
func zipSource(sourcePath, targetZipPath string) error {
	zipFile, err := os.Create(targetZipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a header for the zip file
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Set the relative path inside the zip
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			relPath += "/"
		} else {
			header.Method = zip.Deflate
		}
		header.Name = relPath

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

func rotateBackups(backupDir, sourceName string, limit int) error {
	pattern := filepath.Join(backupDir, fmt.Sprintf("%s_*.zip", sourceName))
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(files) <= limit {
		return nil
	}

	// Sort by modification time (oldest first)
	sort.Slice(files, func(i, j int) bool {
		infoI, _ := os.Stat(files[i])
		infoJ, _ := os.Stat(files[j])
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	// Delete excessive files
	toDelete := len(files) - limit
	for i := 0; i < toDelete; i++ {
		fmt.Printf("🗑️  Rotating old backup: %s\n", filepath.Base(files[i]))
		if err := os.Remove(files[i]); err != nil {
			fmt.Printf("⚠️  Warning: failed to delete %s: %v\n", files[i], err)
		}
	}
	return nil
}

func performBackup(wg *sync.WaitGroup, src SourceConfig, backupDir string, limit int) {
	defer wg.Done()

	timestamp := time.Now().Format("20060102_150405")
	zipName := fmt.Sprintf("%s_%s.zip", src.Name, timestamp)
	targetPath := filepath.Join(backupDir, zipName)

	fmt.Printf("📦 Starting backup for: %s -> %s\n", src.Name, zipName)

	if err := zipSource(src.Path, targetPath); err != nil {
		log.Printf("❌ Failed to backup %s: %v\n", src.Name, err)
		return
	}

	fmt.Printf("✅ Backup complete: %s\n", zipName)

	if err := rotateBackups(backupDir, src.Name, limit); err != nil {
		log.Printf("⚠️  Failed to rotate backups for %s: %v\n", src.Name, err)
	}
}

func main() {
	cfgFile := "config.json"
	cfg, err := loadConfig(cfgFile)
	if err != nil {
		log.Fatalf("❌ Error loading config: %v", err)
	}

	// Ensure backup directory exists
	if err := os.MkdirAll(cfg.BackupDir, 0755); err != nil {
		log.Fatalf("❌ Failed to create backup dir: %v", err)
	}

	fmt.Printf("🚀 Starting Go-Backup-Manager (Sources: %d, Rotation: %d)\n", len(cfg.Sources), cfg.RotationLimit)
	fmt.Println("---------------------------------------------------------")

	var wg sync.WaitGroup
	for _, src := range cfg.Sources {
		wg.Add(1)
		go performBackup(&wg, src, cfg.BackupDir, cfg.RotationLimit)
	}

	wg.Wait()
	fmt.Println("---------------------------------------------------------")
	fmt.Println("🎉 All backup tasks finished!")
}
