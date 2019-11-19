// Ch10ex04 takes any number of packages as parameters (default: package of
// current directory) and prints all packages in the workspace that depend on it
// transitively.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
)

// Package represents one package as returned by go list.
type Package struct {
	ImportPath string   // import path
	Deps       []string // all (recursively) imported dependencies
}

func main() {
	args := []string{"."}
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	// Get canonical import paths for command line arguments
	packages, err := runList(args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ch10ex04: %v\n", err)
		os.Exit(1)
	}
	paths := getUniqueImportPaths(packages)

	// Get all packages in workspace
	packages, err = runList("...")

	// Get packages that depend on argument packages
	dependers := checkDeps(paths, packages)

	for _, d := range dependers {
		fmt.Println(d)
	}
}

// runList runs the command
//
//     go list -e -json PKG [PKG...]
//
// for the packages provided as arguments and returns a slice of Package.
func runList(args ...string) ([]Package, error) {
	params := append([]string{"list", "-e", "-json"}, args...)
	cmd := exec.Command("go", params...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var packages []Package
	dec := json.NewDecoder(stdout)
	for {
		var pkg Package
		err := dec.Decode(&pkg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return packages, nil
}

// getUniqueImportPaths looks at all packages and returns a sorted slice of
// unique import paths found.
func getUniqueImportPaths(packages []Package) []string {
	var paths []string
	pmap := make(map[string]bool)
	for _, p := range packages {
		if pmap[p.ImportPath] {
			continue
		}
		pmap[p.ImportPath] = true
		paths = append(paths, p.ImportPath)
	}

	sort.Strings(paths)
	return paths
}

// checkDeps returns the import paths of all packages for which any of their
// dependencies exists in paths.
func checkDeps(paths []string, packages []Package) []string {
	var dependers []string
	for _, p := range packages {
		if containsAny(p.Deps, paths) {
			dependers = append(dependers, p.ImportPath)
		}
	}
	return dependers
}

// containsAny returns true if any element of s1 exists in s2.
func containsAny(s1, s2 []string) bool {
	for _, el := range s1 {
		if contains(s2, el) {
			return true
		}
	}
	return false
}

// contains returns true if slice s conatains element elem.
func contains(s []string, elem string) bool {
	for _, e := range s {
		if e == elem {
			return true
		}
	}
	return false
}
