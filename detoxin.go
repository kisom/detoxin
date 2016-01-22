package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

var ansiCodes = regexp.MustCompile(`(\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[mGK])`)
var backspace = regexp.MustCompile(`(.\x08+)`)

func detox(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		line = ansiCodes.ReplaceAllString(line, "")
		line = backspace.ReplaceAllString(line, "")
		fmt.Fprintln(w, line)
	}
}

func inplaceDetox(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	tmp, err := ioutil.TempFile("", "detox")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	detox(file, tmp)
	tmp.Close()
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	tmp, err = os.Open(tmp.Name())
	if err != nil {
		return err
	}

	_, err = io.Copy(file, tmp)
	return err
}

func detoxFile(path string, inplace bool) {
	if inplace {
		err := inplaceDetox(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	} else {
		file, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		detox(file, os.Stdout)
		file.Close()
	}
}

func main() {
	var inplace, failed bool

	flag.BoolVar(&inplace, "i", false, "detox files in-place")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		if inplace {
			fmt.Fprintf(os.Stderr,
				"inplace doesn't make sense when reading from stdin.")
			os.Exit(1)
		}
		detox(os.Stdin, os.Stdout)
	} else {
		for _, path := range args {
			detoxFile(path, inplace)
		}
	}

	if failed {
		os.Exit(1)
	}
}
