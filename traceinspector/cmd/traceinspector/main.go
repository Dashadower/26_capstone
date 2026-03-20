package traceinspector

import (
	"flag"
	"fmt"
)

func main() {
	// https://pkg.go.dev/flag#String
	input_path := flag.String("gofile", "", "")
	// print_cfg := flag.String("--print-cfg", "", "")
	flag.Parse()
	if *input_path == "" {
		panic("need to pass input go file path with --gofile")
	}

	just_print_cfg := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "--print-cfg" {
			just_print_cfg = true
		}
	})
	if just_print_cfg {
		// print_cfg()
	}

	fmt.Printf("Parsing %s\n", *input_path)

}
