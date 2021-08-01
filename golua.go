package main

import (
	"fmt"
	"os"

	"github.com/1lann/golua/lua"
	"github.com/1lann/golua/ssap"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

const hello = `
package main

import "fmt"

var message string = "Hello, World!"

func printStr(str string) {
	fmt.Println(str)
}

func main() {
	printStr(message)
}
`

func main() {
	os.Chdir("./sample")
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.LoadAllSyntax,
	}, ".")
	if err != nil {
		panic(err)
	}

	if len(pkgs[0].Errors) > 0 {
		for _, err := range pkgs[0].Errors {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}

	hello, ssaPkgs := ssautil.AllPackages(pkgs, ssa.SanityCheckFunctions)
	hello.Build()

	_, _ = hello, ssaPkgs

	// ssaPkgs[0].Func("Phi").WriteTo(os.Stdout)

	l := lua.New()
	ssap.Process(ssaPkgs[0], l)
	fmt.Println(string(l.Assemble()))

	// for _, pkg := range ssaPkgs {
	// 	fmt.Println(pkg.Pkg.Name())

	// 	for _, imp := range hello.AllPackages() {
	// 		if imp.Pkg.Name() == "fmt" {
	// 			imp.Func("Fprintln").WriteTo(os.Stdout)
	// 		}
	// 	}
	// }

	// // Print out the package.
	// hello.WriteTo(os.Stdout)

	// // Print out the package-level functions.
	// hello.Func("init").WriteTo(os.Stdout)
	// hello.Func("main").WriteTo(os.Stdout)
	// hello.Func("printStr").WriteTo(os.Stdout)
	// for _, imp := range hello.Prog.AllPackages() {
	// 	if imp.Pkg.Name() == "fmt" {
	// 		imp.Func("Println").WriteTo(os.Stdout)
	// 	}
	// }
	// // for k := range hello.Members {
	// // 	fmt.Println(k)
	// // }
}

func demo() {
	a := make(chan int, 5)

	go func() {
		for v := range a {
			fmt.Println(v)
		}
	}()

	for i := 0; i < 5; i++ {
		a <- i
	}
	close(a)
}
