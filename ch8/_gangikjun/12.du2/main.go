package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var verbose = flag.Bool("v", false, "show verbose progress messages")

func main() {
	// 최초 디렉토리 설정
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// 파일 트리 탐색
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()

	// 주기적으로 결과 출력
	var tick <-chan time.Time
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	if *verbose {
		tick = ticker.C
	}
	var nfiles, nbytes int64
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	// 1e9 == 1 * 10의9제곱' == 1000000000
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

// walkDir은 dir부터 파일 트리를 재귀적으로 탐색하며
// 찾은 파일의 크기를 fileSizes로 보낸다.
func walkDir(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

// dirents 는 dir 디렉토리의 항목을 반환
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}
