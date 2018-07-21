package dockerfile

// Instruction implements a dockerfile instruction interface
type Instruction interface {
	Key() string
	String() string
}

// RawInstruction is a single dockerfile instruction
type RawInstruction struct {
	Op   string
	Data []byte
}

// Key returns the op key
func (ri *RawInstruction) Key() string {
	return ri.Op
}

func (ri *RawInstruction) String() string {
	return string(ri.Data)
}

// RawInstructions are a set of instructions in order
type RawInstructions []*RawInstruction

// HasOp returns true if the instruction set contains the operation
// in any one of its items
func (insts RawInstructions) HasOp(op string) bool {
	for _, i := range insts {
		if i.Op == op {
			return true
		}
	}

	return false
}
