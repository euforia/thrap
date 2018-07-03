package crt

// func Test_scopeVars(t *testing.T) {
// 	stack, err := manifest.LoadManifest("../test-fixtures/thrap.hcl")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	errs := stack.Validate()
// 	if len(errs) > 0 {
// 		t.Fatal(errs)
// 	}

// 	dr := registry.DockerRuntime{}
// 	dr.Init(registry.Config{})
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
// 			// 		fmt.Println(ic.ExposedPorts)
// 			var i int
// 			comp.Ports = make(map[string]int32, len(ic.ExposedPorts))
// 			for k := range ic.ExposedPorts {
// 				comp.Ports[k.Port()] = int32(k.Int())
// 				i++
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

// 	err = bldr.BuildComponent(ctx, "default", bcomp, RequestOptions{Output: os.Stdout})
// 	assert.Nil(t, err)

// }

// func Test_DockerOrchestrator_Deploy(t *testing.T) {
// 	st := &thrapb.Stack{ID: "test",
// 		Components: map[string]*thrapb.Component{
// 			"vault": &thrapb.Component{
// 				ID:   "vault",
// 				Name: "vault", Version: "0.10.3",
// 			},
// 			"consul": &thrapb.Component{
// 				ID:   "consul",
// 				Name: "consul", Version: "1.1.0",
// 			},
// 			"api": &thrapb.Component{
// 				ID:      "api",
// 				Name:    "api",
// 				Version: "latest",
// 				Head:    true,
// 				Build: &thrapb.Build{
// 					Dockerfile: "test.dockerfile",
// 					Context:    "../test-fixtures",
// 				},
// 			},
// 		},
// 	}

// 	bldr, _ := NewDockerRuntime()

// 	defer bldr.tearDown(context.Background(), st)

// 	_, _, err := bldr.Deploy(st, RequestOptions{Output: os.Stdout})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// <-time.After(2 * time.Second)

// }
