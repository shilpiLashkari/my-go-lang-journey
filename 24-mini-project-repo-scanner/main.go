package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// RepoStatus holds the scanned result for a single git repository.
type RepoStatus struct {
	Path   string
	Branch string
	Dirty  bool
	Error  string
}

func main() {
	scanPath := flag.String("path", ".", "The parent directory to scan for Git repositories")
	flag.Parse()

	absPath, err := filepath.Abs(*scanPath)
	if err != nil {
		fmt.Printf("❌ Invalid path: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("📂 Go-Repo-Scanner")
	fmt.Println("==================================================")
	fmt.Printf("🔍 Scanning: %s\n\n", absPath)

	// 1. Discover all git repos
	repos := findRepos(absPath)

	if len(repos) == 0 {
		fmt.Println("No Git repositories found.")
		return
	}

	fmt.Printf("Found %d repositories. Checking status concurrently...\n\n", len(repos))

	// 2. Concurrently check each repo's status
	results := make(chan RepoStatus, len(repos))
	var wg sync.WaitGroup

	for _, repo := range repos {
		wg.Add(1)
		go func(repoPath string) {
			defer wg.Done()
			results <- checkRepo(repoPath)
		}(repo)
	}

	// Close the channel once all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// 3. Collect and print results
	fmt.Printf("%-6s %-20s %-40s\n", "STATUS", "BRANCH", "PATH")
	fmt.Println(strings.Repeat("-", 70))

	for res := range results {
		if res.Error != "" {
			fmt.Printf("%-6s %-20s %s\n", "⚠️", "?", res.Path)
			continue
		}

		status := "✅"
		if res.Dirty {
			status = "🔴"
		}

		// Truncate long paths for display
		displayPath := res.Path
		if len(displayPath) > 38 {
			displayPath = "..." + displayPath[len(displayPath)-35:]
		}

		branch := res.Branch
		if len(branch) > 18 {
			branch = branch[:18] + ".."
		}

		fmt.Printf("%-6s %-20s %s\n", status, branch, displayPath)
	}

	fmt.Println(strings.Repeat("-", 70))
	fmt.Println("✅ = Clean | 🔴 = Dirty (uncommitted changes)")
}

// findRepos walks the root path and returns directories that contain a .git folder.
func findRepos(root string) []string {
	var repos []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip unreadable dirs
		}

		if !info.IsDir() {
			return nil
		}

		// Check if this directory has a .git subfolder
		gitDir := filepath.Join(path, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			repos = append(repos, path)
			// Don't descend into the repo itself to avoid finding nested repos
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		fmt.Printf("⚠️  Walk error: %v\n", err)
	}

	return repos
}

// checkRepo runs git commands to determine the branch and dirty status of a repo.
func checkRepo(repoPath string) RepoStatus {
	status := RepoStatus{Path: repoPath}

	// Get the current branch name
	branchCmd := exec.Command("git", "-C", repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, err := branchCmd.Output()
	if err != nil {
		status.Error = err.Error()
		return status
	}
	status.Branch = strings.TrimSpace(string(branchOut))

	// Get the short status (empty output means clean)
	statusCmd := exec.Command("git", "-C", repoPath, "status", "--short")
	statusOut, err := statusCmd.Output()
	if err != nil {
		status.Error = err.Error()
		return status
	}

	status.Dirty = len(strings.TrimSpace(string(statusOut))) > 0

	return status
}
