package glob

import (
	"io/fs"
	"path/filepath"
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
