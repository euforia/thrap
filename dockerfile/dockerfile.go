package dockerfile

import (
	"bytes"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Stage is a group of instructions containing a single FROM statement
type Stage []Instruction

// ID returns the as field from a FROM instruction
func (stage Stage) ID() string {
	for i := range stage {
		if f, ok := stage[i].(*From); ok {
			return f.As
		}
	}
	return ""
}

// HasOp returns true if the instruction set contains the operation
// in any one of its items
func (stage Stage) HasOp(op string) bool {
	for _, i := range stage {
		if i.Key() == op {
			return true
		}
	}
	return false
}

// GetOp returns the first occurence of an instruction associated to the op
// and its index in the stage
func (stage Stage) GetOp(op string) (Instruction, int) {
	for i := range stage {
		if stage[i].Key() == op {
			return stage[i], i
		}
	}
	return nil, -1
}

// Dockerfile holds a completely parsed dockerfile
type Dockerfile struct {
	Stages []Stage
}

// StepCount returns the total number of steps in the dockerfile
func (df *Dockerfile) StepCount() int {
	var c int
	for _, st := range df.Stages {
		c += len(st)
	}
	return c
}

// AddInstruction adds an instuction to the given stage it handles the addition
// of instruction based on the type instruction applying the associated logic
func (df *Dockerfile) AddInstruction(stage int, inst Instruction) error {
	st := df.Stages[stage]

	key := inst.Key()
	var nst Stage

	switch key {
	case KeyFrom:
		if st.HasOp(KeyFrom) {
			return errors.New("already has FROM")
		}
		nst = append([]Instruction{inst}, st...)

	case KeyEntrypoint:
		nst = append(st, inst)

	case KeyArg:
		// Only add if does not exist
		ai := inst.(*Arg)
		for _, s := range st {
			if s.Key() != KeyArg {
				continue
			}
			if s.(*Arg).Name == ai.Name {
				return nil
			}
		}
		nst = df.insertInstruction(st, inst)

	default:
		nst = df.insertInstruction(st, inst)

	}

	df.Stages[stage] = nst

	return nil
}

// insertInstruction inserts the instruction after the from statement.  It copies
// the exsting instructions into a new slice, inserts the new instruction and
// returns the new slice
func (df *Dockerfile) insertInstruction(st Stage, inst Instruction) Stage {
	_, j := st.GetOp(KeyFrom)

	nst := make(Stage, len(st)+1)
	copy(nst, st[:j+1])
	copy(nst[j+2:], st[j+1:])
	nst[j+1] = inst

	return nst
}

func (df *Dockerfile) String() string {
	var str string
	for _, stage := range df.Stages {
		for _, s := range stage {
			// Nil instruction treated as new line
			if s == nil {
				str += "\n"
				continue
			}
			str += s.Key() + " " + s.String() + "\n"

		}
		str += "\n"
	}

	return str
}

// RawDockerfile contains parsed elements of a dockerfile
type RawDockerfile struct {
	Stages []RawInstructions
}

// ParseRaw parses a raw dockerfile to concrete dockerfile
func ParseRaw(raw *RawDockerfile) *Dockerfile {
	df := &Dockerfile{Stages: make([]Stage, len(raw.Stages))}

	for j, stage := range raw.Stages {
		df.Stages[j] = make(Stage, len(stage))
		for i, s := range stage {
			ki, err := ParseInstruction(s)
			if err == nil {
				df.Stages[j][i] = ki
			} else {
				df.Stages[j][i] = stage[i]
			}
		}
	}

	return df
}

// ParseFile parses a dockerfile at the given path into a raw format
func ParseFile(fpath string) (d *RawDockerfile, err error) {
	var b []byte
	b, err = ioutil.ReadFile(fpath)
	if err == nil {
		d, err = ParseBytes(b)
	}

	return
}

// ParseBytes parses dockerfile bytes to a set of instructions
func ParseBytes(b []byte) (*RawDockerfile, error) {
	var (
		df    = &RawDockerfile{Stages: make([]RawInstructions, 0)}
		lines = bytes.Split(b, []byte("\n"))
		out   = make(RawInstructions, 0, len(lines))
		buf   = make([]byte, 0)
	)

	for _, line := range lines {

		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			out = append(out, &RawInstruction{Op: KeyComment, Data: line[1:]})
			continue
		}

		i := bytes.IndexByte(line, ' ')
		if i == 0 {
			buf = append(buf, append([]byte("\n"), line...)...)
			continue
		}

		if op, ok := isCmd(line[:i]); ok {

			if len(buf) > 0 {
				out[len(out)-1].Data = append(out[len(out)-1].Data, buf...)
				buf = make([]byte, 0)
			}

			if op == KeyFrom && len(out) != 0 {
				stage := out
				df.Stages = append(df.Stages, stage)
				out = make(RawInstructions, 0, len(lines))
			}

			inst := &RawInstruction{string(line[:i]), line[i+1:]}
			out = append(out, inst)

		} else {
			buf = append(buf, append([]byte("\n"), line...)...)
		}

	}

	df.Stages = append(df.Stages, out)
	return df, nil
}
