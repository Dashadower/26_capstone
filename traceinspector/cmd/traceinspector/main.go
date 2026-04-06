package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"traceinspector"
	"traceinspector/imp"
)

func main() {
	// https://pkg.go.dev/flag#String
	input_path := flag.String("gofile", "", "")
	_ = flag.Bool("print-cfg", false, "whether to just print cfg and exit")
	_ = flag.Bool("print-imp", false, "whether to just print the translated Imp code and exit")
	_ = flag.Bool("interpret-imp", false, "whether to just interpret the translated Imp code and exit")
	flag.Parse()
	if *input_path == "" {
		panic("need to pass input go file path with --gofile")
	}

	just_print_cfg := false
	just_print_imp := false
	just_interpret := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "print-cfg" {
			just_print_cfg = true
		} else if f.Name == "print-imp" {
			just_print_imp = true
		} else if f.Name == "interpret-imp" {
			just_interpret = true
		}
	})

	fset := token.NewFileSet()
	// https://pkg.go.dev/go/parser#Mode
	file, err := parser.ParseFile(fset, *input_path, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error while parsing", input_path, "-", err)
		return
	}

	if just_print_cfg {
		traceinspector.Print_cfg(file, fset)
		return
	}

	imp_functions := imp.Translate_ast_file_to_imp(file, fset)

	if just_print_imp {
		for _, fun := range imp_functions {
			fmt.Println(fun)
		}
	}

	if just_interpret {
		interpreter := imp.ImpInterpreter{Functions: imp_functions}
		interpreter.Interpret_main()
	}
}
