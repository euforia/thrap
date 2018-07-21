package analysis

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var extLanguageMap = map[string]string{
	".c":   "C",
	".cpp": "C++",
	".cs":  "C#",
	".go":  "Go",
	//"hcl":  "HCL",
	".java": "Java",
	".js":   "JavaScript",
	//"json": "JSON",
	//"md":   "Markdown",
	".py": "Python",
	".sh": "Shell",
	//"toml": "TOML",
	//"yaml": "YAML",
}

// FileAnalysis hold file type analysis
type FileAnalysis struct {
	// File extension
	Ext string
	// Programming language
	Language string
	// File extension occurrences
	Count int
	// Percent representation
	Percent float64
}

func (fa *FileAnalysis) String() string {

	return fmt.Sprintf(`FileAnalysis {
   Language = "%s"
   Ext      = "%s"
   Count    = %d
   Percent  = %f
}`, fa.Language, fa.Ext, fa.Count, fa.Percent)

}

// FileTypeSpread holds file types and associated count in a given analysis
// context
type FileTypeSpread struct {
	// Excludes all dirs, .git dir
	Total int
	// Keys sorted by count in ascending order
	keys   []string
	spread map[string]*FileAnalysis
}

// Spread returns an analysis spread in ascending order
func (s *FileTypeSpread) Spread() []*FileAnalysis {
	out := make([]*FileAnalysis, len(s.keys))
	for i, k := range s.keys {
		out[i] = s.spread[k]
	}
	return out
}

// Highest returns an analysis of the highest occuring file extension
func (s *FileTypeSpread) Highest() *FileAnalysis {
	if len(s.keys) == 0 {
		return nil
	}
	hkey := s.keys[len(s.keys)-1]
	return s.spread[hkey]
}

func (s *FileTypeSpread) Len() int {
	return len(s.spread)
}

func (s *FileTypeSpread) Less(i, j int) bool {
	return s.spread[s.keys[i]].Count < s.spread[s.keys[j]].Count
}

func (s *FileTypeSpread) Swap(i, j int) {
	s.keys[i], s.keys[j] = s.keys[j], s.keys[i]
}

// BuildFileTypeSpread returns a an analysis file type spread
func BuildFileTypeSpread(dirpath string) *FileTypeSpread {
	if dirpath[len(dirpath)-1] != '/' {
		dirpath += "/"
	}

	fts := new(FileTypeSpread)
	exts := make(map[string]int)

	filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		// Skip VCS files
		p := strings.TrimPrefix(path, dirpath)
		if len(p) > 3 {
			switch p[:4] {
			case ".git", ".svn":
				return nil
			}
		}

		// Update the total
		fts.Total++

		ext := filepath.Ext(info.Name())
		if len(ext) > 0 {
			exts[ext]++
		}

		return nil
	})

	fts.keys = make([]string, 0, len(exts))
	fts.spread = make(map[string]*FileAnalysis, len(exts))

	for k, v := range exts {
		fts.keys = append(fts.keys, k)

		fa := &FileAnalysis{
			Ext:     k,
			Count:   v,
			Percent: (float64(v) / float64(fts.Total)) * 100,
		}
		fa.Language, _ = extLanguageMap[fa.Ext]

		fts.spread[fa.Ext] = fa
	}

	sort.Sort(fts)

	return fts
}
