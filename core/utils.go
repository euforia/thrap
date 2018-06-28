package core

// type foo struct {
// 	st *thrapb.Stack

// 	d *orchestrator.DockerOrchestrator
// 	r *registry.DockerRuntime
// }

// func (c *foo) foo() {
// 	for _, comp := range c.st.Components {
// 		cont := thrapb.NewContainer(c.st.ID, comp.ID)

// 		if comp.IsBuildable() {
// 			// Set fully qualified name including stack id for components
// 			// being built
// 			cont.Container.Image = filepath.Join(c.st.ID, comp.Name) + ":" + comp.Version
// 		} else {
// 			// Set user provided name otherwise
// 			cont.Container.Image = comp.Name + ":" + comp.Version
// 			//c.fleshOut(cont, comp.Name, comp.Version)
// 		}

// 	}
// }

// func (c *foo) fleshOut(cont *thrapb.Container, name, tag string) error {
// 	err := c.r.ImagePull(context.Background(), cont.Container.Image)
// 	if err != nil {
// 		return err
// 	}

// 	cfg, err := c.r.ImageConfig(name, tag)
// 	if err == nil {
// 		return nil
// 	}

// }
