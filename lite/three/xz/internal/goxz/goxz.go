// Package main contains an example command line tool goxz built on top of the xz package.
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jamespfennell/xz"
	"io"
	"os"
)

const usage = `Compress or decompress in the xz format.

Usage:

   goxz [<option>...] <input_file> [<output_file>]

If <input_file> has extension .xz, its contents will be decompressed.
Otherwise, its contents will be compressed. These defaults can be
overridden with -c and -d.

If <input_file> is - input will be read from stdin; if <output_file>
is - output will be written to stdout.

If <output_file> is not provided, it defaults to the following:

 * If <input_file> is <file>.xz: decompress to <file>.
 * If <input_file> does not have extension .xz: compress to <input_file>.xz.
 * If <input_file> is stdin: compress to stdout.

Options:

  -c, -z, --compress  force compression
  -d, --decompress    force decompression
  -0 ... -9           specify the compression level (default is 6)
  -f, --force         overwrite the output file if it already exists
  -h, --help          print this help message
`

func main() {
	p, err := newParamsFromArgs(os.Args[1:])
	nilErrOrReturn(err)

	input, err := openInput(p.InputFile)
	nilErrOrReturn(err)
	defer func() {
		nilErrOrReturn(input.Close())
	}()

	output, err := openOutput(p.OutputFile, p.ForceOutputFileWrite)
	nilErrOrReturn(err)
	defer func() {
		nilErrOrReturn(output.Close())
	}()

	if p.IsCompression {
		w := xz.NewWriterLevel(output, p.Level)
		_, err = io.Copy(w, input)
		nilErrOrReturn(err)
		nilErrOrReturn(w.Close())
	} else {
		r := xz.NewReader(input)
		_, err = io.Copy(output, r)
		nilErrOrReturn(err)
		nilErrOrReturn(r.Close())
	}
}

type params struct {
	Level                int
	IsCompression        bool
	InputFile            string
	OutputFile           string
	ForceOutputFileWrite bool
}

func newParamsFromArgs(args []string) (p params, err error) {
	set := flag.NewFlagSet("", flag.ExitOnError)
	var compressionLevel [10]bool
	var forceCompression bool
	var forceDecompression bool
	for i := xz.BestSpeed; i <= xz.BestCompression; i++ {
		set.BoolVar(&compressionLevel[i], fmt.Sprintf("%d", i), false,
			fmt.Sprintf("Set compression level to %d", i))
	}
	for _, alias := range []string{"z", "c", "compress"} {
		set.BoolVar(&forceCompression, alias, false, "")
	}
	for _, alias := range []string{"d", "decompress"} {
		set.BoolVar(&forceDecompression, alias, false, "")
	}
	for _, alias := range []string{"f", "force"} {
		set.BoolVar(&p.ForceOutputFileWrite, alias, false, "")
	}
	set.Usage = func() {
		_, _ = fmt.Fprint(set.Output(), usage, "\n")
	}
	// In case of an error os.Exit will be called by the flag package, per the flag.ExitOnError setting above
	_ = set.Parse(args)
	p.Level, err = determineCompressionLevel(compressionLevel)
	if err != nil {
		return
	}

	if set.NArg() == 0 {
		_, _ = fmt.Fprintf(set.Output(), "Invalid args: <input_file> must be provided\n")
		set.Usage()
		os.Exit(1)
	}
	if set.NArg() > 2 {
		_, _ = fmt.Fprintf(set.Output(), "Invalid args: too many args passed (flags must come before <input_file>)\n")
		set.Usage()
		os.Exit(1)
	}
	p.InputFile = set.Arg(0)
	if p.IsCompression, err = isCompression(forceCompression, forceDecompression, p.InputFile); err != nil {
		return
	}
	if set.NArg() == 2 {
		p.OutputFile = set.Arg(1)
	} else {
		p.OutputFile, err = autoDetectOutputFile(p.InputFile, p.IsCompression)
		if err != nil {
			return
		}
	}
	return
}

func autoDetectOutputFile(inputFile string, isCompression bool) (string, error) {
	if inputFile == "-" {
		return "-", nil
	}
	if isCompression {
		return inputFile + ".xz", nil
	}
	if len(inputFile) >= 3 && inputFile[len(inputFile)-3:] == ".xz" {
		return inputFile[:len(inputFile)-3], nil
	}
	return "", fmt.Errorf(
		"cannot auto detect the output file to decompress to using the input file %s", inputFile)
}

func determineCompressionLevel(compressionLevel [10]bool) (int, error) {
	result := -1
	for level, set := range compressionLevel {
		if !set {
			continue
		}
		if result != -1 {
			return -1, fmt.Errorf("multiple compression levels cannot be specified (-%d -%d)", result, level)
		}
		result = level
	}
	if result == -1 {
		result = xz.DefaultCompression
	}
	return result, nil
}

func isCompression(forceCompression, forceDecompression bool, inputFile string) (bool, error) {
	if forceCompression && forceDecompression {
		return false, fmt.Errorf("--compress and --decompress cannot both be set")
	}
	if forceCompression {
		return true, nil
	}
	if forceDecompression {
		return false, nil
	}
	if inputFile == "-" {
		return true, nil
	}
	if len(inputFile) >= 3 && inputFile[len(inputFile)-3:] == ".xz" {
		return false, nil
	}
	return true, nil
}

func openInput(inputFile string) (io.ReadCloser, error) {
	if inputFile == "-" {
		return os.Stdin, nil
	}
	return os.Open(inputFile)
}

func openOutput(outputFile string, force bool) (io.WriteCloser, error) {
	if outputFile == "-" {
		return os.Stdout, nil
	}
	_, err := os.Stat(outputFile)
	if err == nil {
		// file already exists
		if force {
			return os.Create(outputFile)
		}
		return nil, fmt.Errorf("output file %s already exists", outputFile)
	}
	if errors.Is(err, os.ErrNotExist) {
		return os.Create(outputFile)
	}
	return nil, err
}

func nilErrOrReturn(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s.\n", err)
		os.Exit(1)
	}
}
