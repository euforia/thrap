package registry

import (
	"fmt"
	"testing"
)

func Test_localDocker(t *testing.T) {
	ld := &DockerRuntime{}
	ld.Init(Config{})

	c, err := ld.ImageConfig("vault", "0.10.3")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(c.Env)
	for k := range c.ExposedPorts {
		fmt.Println(k.Int(), k.Port(), k.Proto())
	}
}
