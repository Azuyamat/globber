package glob

import (
	"io/fs"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
)

type fSMatcher struct {
	matcher *matcher
}

func FSMatcher(pattern string) *fSMatcher {
	return &fSMatcher{
		matcher: Matcher(pattern),
	}
}

func (fm *fSMatcher) Matches(path string) (bool, error) {
	return fm.matcher.Matches(path)
}

func (fm *fSMatcher) Walk(fsys fs.FS, fn func(path string, entry fs.DirEntry) error) error {
	return fm.walkSerial(fsys, ".", fn)
}

func (fm *fSMatcher) walkParallel(fsys fs.FS, root string, fn func(path string, entry fs.DirEntry) error) error {
	type match struct {
		path  string
		entry fs.DirEntry
	}

	var (
		wg        sync.WaitGroup
		errOnce   sync.Once
		firstErr  error
		sem       = make(chan struct{}, runtime.NumCPU()*4)
		stopped   atomic.Bool
		matchChan = make(chan match, 1000)
	)

	var processDir func(string)
	processDir = func(dir string) {
		defer wg.Done()
		defer func() { <-sem }()

		if stopped.Load() {
			return
		}

		entries, err := fs.ReadDir(fsys, dir)
		if err != nil {
			errOnce.Do(func() {
				firstErr = err
				stopped.Store(true)
			})
			return
		}

		for _, entry := range entries {
			if stopped.Load() {
				return
			}

			var path string
			if dir == "." {
				path = entry.Name()
			} else {
				path = dir + "/" + entry.Name()
			}

			matches, err := fm.Matches(path)
			if err != nil {
				errOnce.Do(func() {
					firstErr = err
					stopped.Store(true)
				})
				return
			}

			if matches {
				matchChan <- match{path: path, entry: entry}
			}

			if entry.IsDir() {
				sem <- struct{}{}
				wg.Add(1)
				go processDir(path)
			}
		}
	}

	go func() {
		wg.Wait()
		close(matchChan)
	}()

	wg.Add(1)
	sem <- struct{}{}
	go processDir(root)

	for m := range matchChan {
		if err := fn(m.path, m.entry); err != nil {
			stopped.Store(true)
			firstErr = err
			break
		}
	}

	wg.Wait()
	return firstErr
}

func (fm *fSMatcher) walkSerial(fsys fs.FS, dir string, fn func(path string, entry fs.DirEntry) error) error {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		var path string
		if dir == "." {
			path = entry.Name()
		} else {
			path = dir + "/" + entry.Name()
		}

		matches, err := fm.Matches(path)
		if err != nil {
			return err
		}

		if matches {
			if err := fn(path, entry); err != nil {
				return err
			}
		}

		if entry.IsDir() {
			if err := fm.walkSerial(fsys, path, fn); err != nil {
				return err
			}
		}
	}

	return nil
}

func (fm *fSMatcher) walkParallelV2(fsys fs.FS, root string, fn func(path string, entry fs.DirEntry) error) error {
	var (
		dirWg    sync.WaitGroup
		mu       sync.Mutex
		firstErr error
		sem      = make(chan struct{}, runtime.NumCPU()*2)
		matches  []struct {
			path  string
			entry fs.DirEntry
		}
	)

	var processDir func(string)
	processDir = func(dir string) {
		defer dirWg.Done()
		defer func() { <-sem }()

		entries, err := fs.ReadDir(fsys, dir)
		if err != nil {
			mu.Lock()
			if firstErr == nil {
				firstErr = err
			}
			mu.Unlock()
			return
		}

		for _, entry := range entries {
			var path string
			if dir == "." {
				path = entry.Name()
			} else {
				path = dir + "/" + entry.Name()
			}

			if entry.IsDir() {
				sem <- struct{}{}
				dirWg.Add(1)
				go processDir(path)
			} else {
				matched, err := fm.Matches(path)
				if err != nil {
					mu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					mu.Unlock()
					return
				}

				if matched {
					mu.Lock()
					matches = append(matches, struct {
						path  string
						entry fs.DirEntry
					}{path: path, entry: entry})
					mu.Unlock()
				}
			}
		}
	}

	dirWg.Add(1)
	sem <- struct{}{}
	go processDir(root)

	dirWg.Wait()

	if firstErr != nil {
		return firstErr
	}

	for _, m := range matches {
		if err := fn(m.path, m.entry); err != nil {
			return err
		}
	}

	return nil
}

func (fm *fSMatcher) WalkDirFS(rootPath string, fn func(path string, entry fs.DirEntry) error) error {
	return filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		matches, err := fm.Matches(filepath.ToSlash(relPath))
		if err != nil {
			return err
		}

		if matches {
			return fn(filepath.ToSlash(relPath), d)
		}

		return nil
	})
}
