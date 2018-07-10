package crt

// func Test_scopeVars(t *testing.T) {
// 	if !utils.FileExists("/var/run/docker.sock") {
// 		t.Skip("Skipping: docker file descriptor not found")
// 	}

// 	stack, err := manifest.LoadManifest("../test-fixtures/thrap.hcl")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	errs := stack.Validate()
// 	if len(errs) > 0 {
// 		t.Fatal(errs)
// 	}

// 	dr, err := NewDocker()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	for _, comp := range stack.Components {
// 		var (
// 			ic  *container.Config
// 			err error
// 		)

// 		if comp.IsBuildable() {
// 			ic, err = dr.ImageConfig(stack.ID+"/"+comp.Name, "latest")
// 		} else {
// 			ic, err = dr.ImageConfig(comp.Name, comp.Version)
// 		}

// 		if err == nil {
// 			comp.Ports = make(map[string]int32, len(ic.ExposedPorts))
// 			for k := range ic.ExposedPorts {
// 				comp.Ports[k.Port()] = int32(k.Int())
// 			}
// 		}

// 	}

// 	svars := stack.ScopeVars()
// 	b, _ := json.MarshalIndent(svars, "", "  ")
// 	fmt.Printf("\n%s\n", b)
// }

// func Test_Docker_Build2(t *testing.T) {
// 	bldr, err := NewDocker()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	st, err := manifest.LoadManifest("../test-fixtures/builder.yml")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	st.Validate()

// 	var bcomp *thrapb.Component
// 	for _, comp := range st.Components {
// 		if comp.IsBuildable() {
// 			bcomp = comp
// 			bcomp.Build.Context = "../test-fixtures"
// 			break
// 		}
// 	}
// }

// func Test_Docker_Build(t *testing.T) {
// 	bldr, err := NewDocker()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	st, err := manifest.LoadManifest("../test-fixtures/builder.yml")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	st.Validate()

// 	var bcomp *thrapb.Component
// 	for _, comp := range st.Components {
// 		if comp.IsBuildable() {
// 			bcomp = comp
// 			bcomp.Build.Context = "../test-fixtures"
// 			break
// 		}
// 	}
// 	ctx := context.Background()

// 	err = bldr.Build(ctx, "default", bcomp, RequestOptions{Output: os.Stdout})
// 	assert.Nil(t, err)

// }
