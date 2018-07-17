package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
)

func writeJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Printf("%s\n", b)
}

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
	if utils.FileExists(fpath) {
		return os.ErrExist
	}

	fh, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer fh.Close()

	return writeHCLManifest(stack, fh)
}
