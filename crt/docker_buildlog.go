package crt

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

type dockerBuildStep struct {
	// partial or full hash id. docker build only returns
	// a partial id
	id string
	// cmd that was run in this step
	cmd string
	// true if the step was a cached run
	usedCache bool
	// all additional unparsed data
	data [][]byte
}

func (step *dockerBuildStep) ID() string {
	return step.id
}

func (step *dockerBuildStep) Cmd() string {
	return step.cmd
}

func (step *dockerBuildStep) UsedCache() bool {
	return step.usedCache
}

func (step *dockerBuildStep) Log() string {
	return string(bytes.Join(step.data, []byte("\n")))
}

// Step implements an interface for each step in the build log
type Step interface {
	// Hash id of the step
	ID() string
	// Cmd returns the docker command that was run
	Cmd() string
	// Returns true if the step used the cache
	UsedCache() bool
	// Returns the log data
	Log() string
}

// DockerBuildLog holds a docker build image log.  It is used to parse and
// get info about each step
type DockerBuildLog struct {
	// Multi-writer to write to supplied writer and buffer
	mw io.Writer
	// Copy of the log used to validate build
	*bytes.Buffer
}

// NewDockerBuildLog returns a new DockerBuildLog instance. It takes a writer
// for pass through to be written.
func NewDockerBuildLog(w io.Writer) *DockerBuildLog {
	buf := bytes.NewBuffer(nil)
	return &DockerBuildLog{
		Buffer: buf,
		mw:     io.MultiWriter(buf, w),
	}
}

func (log *DockerBuildLog) Write(b []byte) (int, error) {
	return log.mw.Write(b)
}

// Steps parses and returns all steps from the build log
func (log *DockerBuildLog) Steps() ([]Step, error) {
	var (
		b     = log.Buffer.Bytes() // bytes from log
		steps = make([]Step, 0)    // all steps
		step  *dockerBuildStep     // current step
	)

	lines := bytes.Split(b, []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		header := strings.TrimSpace(string(line[:5]))
		data := bytes.TrimSpace(line[5:])
		switch header {

		case "Step":
			if step != nil {
				steps = append(steps, step)
			}

			i := bytes.IndexRune(data, ':')
			if i < 0 {
				return nil, fmt.Errorf("parse instruction: '%s'", line)
			}
			step = &dockerBuildStep{
				cmd:  string(bytes.TrimSpace(data[i+1:])),
				data: make([][]byte, 0),
			}

		case "--->":
			id := string(data)
			_, err := hex.DecodeString(id)
			if err == nil {
				step.id = id
			} else if strings.Contains(id, "Using cache") {
				step.usedCache = true
			} else {
				// add the whole raw line as we don't know what line this is
				d := make([]byte, len(line))
				copy(d, line)
				step.data = append(step.data, d)
			}

		default:
			// add the whole raw line
			d := make([]byte, len(line))
			copy(d, line)
			step.data = append(step.data, d)

		}

	}
	// append last step
	steps = append(steps, step)
	return steps, nil
}
