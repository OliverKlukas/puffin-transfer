package fileanalyzer

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"go-tui-file-project/internal/firestore"
	"os"
	"path/filepath"
	"sync"
)

type FileAnalyzerConfig struct {
	UseDuplicate bool
	UseSize      bool
	maxSize      int64
}

func walkDir(dir string, result chan<- FileInfo, quit <-chan struct{}, errc chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	visit := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != dir {
			wg.Add(1)
			go walkDir(path, result, quit, errc, wg)
			return filepath.SkipDir
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			content, err := os.ReadFile(path)
			result <- FileInfo{path, hash(content), info, err}
		}()

		select {
		case <-quit:
			return errors.New("traversal canceled")
		default:
			return nil
		}
	}

	err := filepath.Walk(dir, visit)
	if err != nil {
		errc <- err
	}
}

func monitorWorker(wg *sync.WaitGroup, result chan FileInfo, quit chan struct{}, errc chan error) {
	wg.Wait()
	close(result)
	close(quit)
}

func computeDuplicate(res FileInfo, hashMap map[string]string) *Result {
	existing, present := hashMap[res.hash]

	if !present {
		hashMap[res.hash] = res.path
		return nil
	}

	return &Result{res.path, fmt.Sprintf("Duplicate file: %s!", existing)}
}

func hash(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func SearchDir(dir string, config FileAnalyzerConfig, result chan<- Result, errc chan error) {
	quit := make(chan struct{})
	fileInfo := make(chan FileInfo)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go walkDir(dir, fileInfo, quit, errc, wg)
	go monitorWorker(wg, fileInfo, quit, errc)

	hashMap := make(map[string]string)

	for res := range fileInfo {
		if res.err != nil {
			errc <- res.err
		}

		// Duplicate
		if config.UseDuplicate {
			if duplicate := computeDuplicate(res, hashMap); duplicate != nil {
				result <- *duplicate
			}
		}

		// File size
		if config.UseSize && res.info.Size() >= config.maxSize {
			result <- Result{res.path, fmt.Sprintf("File too large: %s!", FormatSize(res.info.Size()))}
		}
	}

	close(result)
	close(errc)
}

func parseFileAnalyzerConfig(arguments []string) FileAnalyzerConfig {
	config := FileAnalyzerConfig{}

	for idx, arg := range arguments {
		switch arg {
		case "duplicate":
			config.UseDuplicate = true
		case "size":
			if len(arguments) > idx+1 {
				parsedSize, err := ParseSize(arguments[idx+1])
				if err != nil {
					fmt.Printf("Error parsing size: %v\n", err)
					continue
				}
				config.maxSize = parsedSize
			} else {
				config.maxSize = 10_000_000_000
			}
			config.UseSize = true
		}
	}

	return config
}

func Run(path string, arguments []string) {
	config := parseFileAnalyzerConfig(arguments)

	if !config.UseDuplicate && !config.UseSize {
		fmt.Println("Error no options passed!")
		return
	}

	result := make(chan Result)
	errc := make(chan error, 1)

	go SearchDir(path, config, result, errc)

	resultFound := false

	for res := range result {
		resultFound = true
		fmt.Printf("Transfering %s reason: %s\n", res.Path, res.Reason)
		err := firestore.Transfer(res.Path)
		if err != nil {
			fmt.Printf("Error transferring file: %v\n", err)
		}
	}

	if err := <-errc; err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if !resultFound {
		fmt.Printf("No results found!\n")
	} else {
		fmt.Printf("Done!\n")
	}
}
