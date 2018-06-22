package manifest

import (
	"fmt"
	"testing"

	"github.com/euforia/hclencoder"
	"github.com/stretchr/testify/assert"
)

var testHCLManifest = `
manifest "foo" {
    components {
      api {
        id   = "api"
        name = "${registry.image.addr}/testdir/api"
        type = "api"

        build {
          dockerfile = "api.dockerfile"
          context    = "."
        }

        secrets {
          destination = "secrets"
        }

        head = true
      }

      db {
        name    = "cockroachdb/cockroach"
        version = "v2.0.2"
        type    = "datastore"
      }

      www {
        name    = "nginx"
        version = "1.15-alpine"
        type    = "web"
      }
    }
}`

func Test_ParseBytes(t *testing.T) {
	mf, err := ParseYAML("../test-fixtures/thrap.yml")
	fatal(t, err)

	assert.Equal(t, 2, len(mf.Components))
	errs := mf.Validate()
	assert.Nil(t, errs)

	b, err := hclencoder.Encode(mf)
	fatal(t, err)

	fmt.Printf("\n%s\n", b)
}

func Test_parse_HCL(t *testing.T) {
	o, err := ParseHCLBytes([]byte(testHCLManifest))
	fatal(t, err)
	assert.Equal(t, "api.dockerfile", o.Components["api"].Build.Dockerfile)
}
