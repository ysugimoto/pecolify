package pecolify

import (
	"bufio"
	"github.com/peco/peco"
	"io/ioutil"
	"os"
	"strings"
)

type Runner struct {
	tmpFile    *os.File
	pipeReader *os.File
	pipeWriter *os.File

	oldStdout     *os.File
	oldArgs       []string
	captureBuffer string
}

func New() *Runner {
	return &Runner{
		oldStdout: os.Stdout,
		oldArgs:   os.Args,
	}
}

func (r *Runner) Transform(data []string) (string, error) {
	var e error

	// create tmpfile
	r.tmpFile, _ = ioutil.TempFile(os.TempDir(), "peconize")
	defer os.Remove(r.tmpFile.Name())

	// swap stdout
	r.pipeReader, r.pipeWriter, e = os.Pipe()
	os.Stdout = r.pipeWriter

	if e = r.runPeco(data); e != nil {
		r.pipeWriter.Close()
		if e != peco.ErrUserCanceled {
			return "", e
		} else {
			return "", nil
		}
	}

	// create blocking channel
	blocker := make(chan int, 1)

	go r.captureOutput(blocker)

	r.pipeWriter.Close()
	<-blocker

	return r.captureBuffer, nil
}

func (r *Runner) captureOutput(blocker chan int) {
	var lines string

	defer r.pipeReader.Close()

	scanner := bufio.NewScanner(r.pipeReader)
	for scanner.Scan() {
		lines += scanner.Text()
	}

	r.captureBuffer = lines
	blocker <- 1
}

func (r *Runner) runPeco(data []string) error {
	buffer := strings.Join(data, "\n")
	ioutil.WriteFile(r.tmpFile.Name(), []byte(buffer), 0644)
	os.Args = []string{"peconize", r.tmpFile.Name()}

	// restore
	defer func() {
		os.Stdout = r.oldStdout
		os.Args = r.oldArgs
	}()

	cli := peco.CLI{}
	if err := cli.Run(); err != nil {
		return err
	}

	return nil
}
