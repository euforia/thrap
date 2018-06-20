package main

import (
	"io"
	"os"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap"
	"github.com/euforia/thrap/thrapb"
)

//
// func fileExists(p string) bool {
// 	_, err := os.Stat(p)
// 	return err == nil
// }

// func writeGitIgnoresFile(dpath string) error {
// 	if fileExists(dpath) {
// 		return os.ErrExist
// 	}
//
// 	ign := vcs.DefaultGitIgnores()
// 	ignores := []byte(strings.Join(ign, "\n"))
//
// 	return ioutil.WriteFile(dpath, ignores, 0644)
// }
//
// func writeSecretsFile(comp *thrapb.Component, dir string) error {
// 	name := filepath.Join(dir, comp.Secrets.Destination)
//
// 	if fileExists(name) {
// 		return os.ErrExist
// 	}
//
// 	fh, err := os.Create(name)
// 	if err == nil {
// 		err = fh.Close()
// 	}
// 	return err
// }

// func writeDockerfile(st *thrapb.Stack, comp *thrapb.Component, lang *languages.Language, dir string) error {
// 	name := filepath.Join(dir, comp.Build.Dockerfile)
//
// 	df, err := thrap.BuildDockerfile(st.ID, comp, lang)
// 	if err != nil {
// 		return err
// 	}
//
// 	if fileExists(name) {
// 		return os.ErrExist
// 	}
//
// 	b := []byte(df.String())
// 	return ioutil.WriteFile(name, b, 0644)
// }

// func writeReadme(st *thrapb.Stack, dir string) error {
// 	apath := filepath.Join(dir, "README.md")
// 	if fileExists(apath) {
// 		return nil
// 	}
//
// 	readme := builder.DefaultReadmeText(st.Name, st.Description)
// 	return ioutil.WriteFile(apath, []byte(readme), 0644)
// }

func writeHCLManifest(stack *thrapb.Stack, w io.Writer) error {

	key := `manifest "` + stack.ID + `"`
	out := map[string]interface{}{
		key: stack,
	}

	b, err := hclencoder.Encode(&out)
	if err == nil {
		w.Write([]byte("\n"))
		w.Write(b)
		w.Write([]byte("\n"))
	}

	return err
}

func writeHCLManifestFile(stack *thrapb.Stack, fpath string) error {
	if thrap.FileExists(fpath) {
		return os.ErrExist
	}

	fh, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer fh.Close()

	return writeHCLManifest(stack, fh)
}
