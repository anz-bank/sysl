package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func (e exit) Error() string {
	return e.message
}

// main3 is the real main function. It takes its output streams and command-line
// arguments as parameters to support testability.
func main3(stdout, stderr io.Writer, args []string) error {
	flags := flag.NewFlagSet(args[0], flag.PanicOnError)

	root := flags.String("root", ".", "sysl root directory for input files (default: .)")
	output := flags.String("o", "", "output file name")
	mode := flags.String("mode", "textpb", "output mode")

	flags.Parse(args[1:])

	switch *mode {
	case "", "textpb", "json":
	default:
		return fmt.Errorf("Invalid -mode %#v", *mode)
	}

	filename := flags.Arg(0)

	fmt.Fprintf(stderr, "Args: %v\n", flags.Args())
	fmt.Fprintf(stderr, "Root: %s\n", *root)
	fmt.Fprintf(stderr, "Module: %s\n", filename)
	fmt.Fprintf(stderr, "Mode: %s\n", *mode)
	format := strings.ToLower(*output)
	toJSON := *mode == "json" || *mode == "" && strings.HasSuffix(format, ".json")
	fmt.Fprintf(stderr, "%s\n", filename)
	mod, err := Parse(filename, *root)
	if err != nil {
		return err
	}
	if mod != nil {
		if toJSON {
			if *output == "-" {
				return FJSONPB(stdout, mod)
			}
			return JSONPB(mod, *output)
		}
		if *output == "-" {
			return FTextPB(stdout, mod)
		}
		return TextPB(mod, *output)
	}
	return nil
}

// main2 calls main3 and handles any errors it returns. It takes its output
// streams and command-line arguments and even main3 as parameters to support
// testability.
func main2(
	stdout, stderr io.Writer, args []string,
	main3 func(stdout, stderr io.Writer, args []string) error,
) int {
	if err := main3(stdout, stderr, args); err != nil {
		fmt.Fprintln(stderr, err.Error())
		if err, ok := err.(exit); ok {
			return err.code
		}
		return 1
	}
	return 0
}

// main is as small as possible to minimise its no-coverage footprint.
func main() {
	if rc := main2(os.Stdout, os.Stderr, os.Args, main3); rc != 0 {
		os.Exit(rc)
	}
}
