/**
 * pecolify: peco works your program
 *
 * @package pecolify
 * @author Yoshiaki Sugimoto
 * @license MIT
 */
package pecolify

import (
	"bufio"
	"context"
	"github.com/peco/peco"
	"io/ioutil"
	"os"
	"strings"
)

//
// Runner struct
//
// stack some temporary data
//
type Runner struct {
	tmpFile    *os.File
	pipeReader *os.File
	pipeWriter *os.File

	oldStdout     *os.File
	oldArgs       []string
	captureBuffer string
}

//
// Create Runner
//
// @public
// @returns *Runner
//
func New() *Runner {
	return &Runner{
		oldStdout: os.Stdout,
		oldArgs:   os.Args,
	}
}

//
// Data pass to peco and get selected line
//
// @public
// @param data []string
// returns string, error
//
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
	blocker := make(chan struct{}, 1)

	go r.captureOutput(blocker)

	r.pipeWriter.Close()
	<-blocker

	return r.captureBuffer, nil
}

//
// Capture peco's output
//
// @private
// @param blocker chan int
// returns void
//
func (r *Runner) captureOutput(blocker chan int) {
	var lines string

	defer r.pipeReader.Close()

	scanner := bufio.NewScanner(r.pipeReader)
	for scanner.Scan() {
		lines += scanner.Text()
	}

	r.captureBuffer = lines
	blocker <- struct{}{}
}

//
// Run the peco
//
// @private
// @param data []string
// returns error
//
func (r *Runner) runPeco(data []string) error {
	buffer := strings.Join(data, "\n")
	ioutil.WriteFile(r.tmpFile.Name(), []byte(buffer), 0644)

	// Make dummy arguments.
	// Peco will process from tempfile...
	os.Args = []string{"peconize", r.tmpFile.Name()}

	// restore
	defer func() {
		os.Stdout = r.oldStdout
		os.Args = r.oldArgs
	}()

	p := peco.New()
	p.Run(context.Background())
	p.PrintResults()

	return nil
}
