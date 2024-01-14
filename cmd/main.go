package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Zyko0/please"
	"github.com/otiai10/copy"
)

var (
	loaderCodeFmt = []byte(`package main

import (
	"github.com/Zyko0/please"
)

func init() {
	please.SetMode(please.Mode("%s")) // TODO: tmp because no private repo resolution atm
	please.GlitchMe()
}
`)
)

var (
	debug bool
	mode  string
)

func debugLog(format string, args ...any) {
	if !debug {
		return
	}
	fmt.Printf(format+"\n", args...)
}

func errorLog(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func resolvePackagePath(pkg, version string) (string, error) {
	gopath, _ := os.LookupEnv("GOPATH")
	gopath = path.Clean(gopath)
	// Try /src
	srcPath := path.Join(gopath, "src", pkg)
	if fi, err := os.Stat(srcPath); err == nil && fi.IsDir() {
		return srcPath, nil
	}
	// Normalize pkg mod with '!'
	var normalized string
	for _, r := range pkg {
		if r >= 'A' && r <= 'Z' {
			normalized += "!" + string(r+32)
		} else {
			normalized += string(r)
		}
	}
	pkg = normalized
	// Try /pkg/mod without version
	pkgPath := path.Join(gopath, "pkg", "mod", pkg)
	if fi, err := os.Stat(pkgPath); err == nil && fi.IsDir() {
		return pkgPath, nil
	}
	// Try /pkg/mod without version
	pkgPathVer := path.Join(gopath, "pkg", "mod", pkg+version)
	if fi, err := os.Stat(pkgPathVer); err == nil && fi.IsDir() {
		return pkgPathVer, nil
	}
	// Look for latest /pkg/mod version suffix
	pkgPathParent, _ := path.Split(pkgPath)
	_, name := path.Split(pkg)
	var candidates []fs.FileInfo
	err := filepath.WalkDir(pkgPathParent, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if d.Name() == "" {
			return nil
		}
		dirName := strings.Split(d.Name(), "@")
		if dirName[0] == name {
			info, err := d.Info()
			if err != nil {
				return err
			}
			candidates = append(candidates, info)
		}

		return nil
	})
	if err != nil {
		return "", err
	}
	// Pick the most recent version "naively"
	if len(candidates) > 0 {
		sortFunc := func(i, j int) bool {
			return candidates[i].Name() > candidates[j].Name()
		}
		// If "@latest" has been specified, try to naively find by strcmp
		if version != "@latest" {
			sortFunc = func(i, j int) bool {
				return candidates[i].ModTime().After(candidates[j].ModTime())
			}
		}
		sort.SliceStable(candidates, sortFunc)
		dirName := candidates[0].Name()
		return path.Clean(path.Join(pkgPathParent, dirName)), nil
	}

	return "", errors.New("couldn't resolve go package location")
}

func validateMode() error {
	mode = strings.ToUpper(mode)
	switch please.Mode(mode) {
	case please.None, please.Default, please.Medium, please.Unsafe:
		return nil
	default:
		return errors.New("err: invalid glitch mode '" + mode + "'")
	}
}

func main() {
	gopath, ok := os.LookupEnv("GOPATH")
	if !ok {
		log.Fatal("no GOPATH environment variable found")
	}

	flag.BoolVar(&debug, "debug", false, "Logs extra information (default: false)")
	flag.StringVar(&mode, "mode", "default", "Select a glitch mode (none,default,medium,unsafe)")
	flag.Usage()
	flag.Parse()
	goPkg := flag.Arg(0)
	if goPkg == "" {
		errorLog("err: missing ebitengine go package argument")
		flag.Usage()
		return
	}
	if err := validateMode(); err != nil {
		errorLog(err.Error())
		flag.Usage()
		return
	}

	debugLog("GOPATH: %s", gopath)
	goPkgLocalExec := flag.Arg(1)
	if goPkgLocalExec != "" {
		debugLog("Go application path: %s", path.Clean(path.Join(goPkg, goPkgLocalExec)))
	}
	goPkgRoot := path.Clean(goPkg)
	goPkgRootNoVersion := goPkgRoot
	var version string
	if idx := strings.LastIndex(goPkgRoot, "@"); idx >= 0 {
		version = goPkgRoot[idx:]
		debugLog("Version: %s", version)
		goPkgRootNoVersion = goPkgRoot[:idx]
	}

	// Create a temporary directory
	debugLog("Creating a temporary directory...")
	pathTmpDir, err := os.MkdirTemp("", "pleaseglitch_")
	if err != nil {
		errorLog("err: couldn't create temporary directory: %v", err)
		return
	}
	defer func() {
		debugLog("Deleting temporary directory...")
		err := os.RemoveAll(pathTmpDir)
		if err != nil {
			errorLog("err: couldn't remove temporary directory: %v", err)
			return
		}
		debugLog("Temporary directory deleted from: %s", pathTmpDir)
	}()
	debugLog("Created a temporary directory at: %s", pathTmpDir)

	// go mod init __glitch => create a temporary go mod to allow download versioned packages
	os.Setenv("GO111MODULE", "on")
	debugLog("Executing 'go mod init %s'", "__glitch")
	cmd := exec.Command("go", "mod", "init", "__glitch")
	cmd.Dir = pathTmpDir
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		debugLog(string(out))
	}
	if err != nil {
		errorLog("err: couldn't execute 'go mod init %s': %v", "__glitch", err)
		return
	}
	// go get
	cmd = exec.Command("go", "get", goPkgRoot)
	cmd.Dir = pathTmpDir
	cmd.Env = os.Environ()
	debugLog("Executing 'go get %s'", goPkgRoot)
	out, err = cmd.CombinedOutput()
	if len(out) > 0 {
		debugLog(string(out))
	}
	if err != nil {
		errorLog("err: couldn't execute 'go get %s': %v", goPkgRoot, err)
		return
	}

	// Resolve go package location
	pathTmpExec := path.Join(pathTmpDir, goPkgLocalExec)
	loaderPath := path.Join(pathTmpExec, "pleaseglitchloader.go")
	debugLog("Resolving go package location...")
	pathGoRepoPkgRoot, err := resolvePackagePath(goPkgRootNoVersion, version)
	if err != nil {
		errorLog("err: couldn't locate cloned go package: %v\n", err)
		return
	}
	debugLog("Go package found at: %s", pathGoRepoPkgRoot)
	// Copy cloned directory to tmp one
	debugLog("Copying go package to temporary directory...")
	err = copy.Copy(pathGoRepoPkgRoot, pathTmpDir, copy.Options{
		Specials:      true,
		AddPermission: fs.ModePerm,
	})
	if err != nil {
		errorLog("err: couldn't copy to temporary directory %v\n", err)
		return
	}
	// Write auto init file into user package
	debugLog("Creating loader file at: %s...", loaderPath)
	err = os.WriteFile(
		loaderPath,
		[]byte(fmt.Sprintf(string(loaderCodeFmt), mode)),
		os.FileMode(0644),
	)
	if err != nil {
		errorLog("err: couldn't write loader file to temporary directory: %v", err)
		return
	}
	debugLog("Loader file created")

	// go get glitch package
	debugLog("Executing 'go get github.com/Zyko0/please'")
	cmd = exec.Command("go", "get", "github.com/Zyko0/please")
	cmd.Dir = pathTmpDir
	cmd.Env = os.Environ()
	out, err = cmd.CombinedOutput()
	if len(out) > 0 {
		debugLog(string(out))
	}
	if err != nil {
		errorLog("err: couldn't execute 'go get github.com/Zyko0/please': %v", err)
		return
	}
	// TODO: this is tmp local tests
	//go mod edit -replace github.com/pselle/bar=/Users/pselle/Projects/bar
	cmd = exec.Command("go", "mod", "edit", "-replace", "github.com/Zyko0/please=C:/Users/Zyko/go/src/github.com/Zyko0/please")
	cmd.Dir = pathTmpDir
	cmd.Env = os.Environ()
	out, err = cmd.CombinedOutput()
	if len(out) > 0 {
		debugLog(string(out))
	}
	if err != nil {
		errorLog("err: couldn't execute 'go mod edit LOCALTEST': %v", err)
		return
	}
	// go mod tidy
	debugLog("Executing 'go mod tidy'")
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = pathTmpDir
	cmd.Env = os.Environ()
	out, err = cmd.CombinedOutput()
	if len(out) > 0 {
		debugLog(string(out))
	}
	if err != nil {
		errorLog("err: couldn't execute 'go mod tidy': %v", err)
		return
	}
	// go run
	debugLog("Executing 'go run .'")
	cmd = exec.Command("go", "run", ".")
	cmd.Dir = pathTmpExec
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Start command asynchronously
	err = cmd.Start()
	if len(out) > 0 {
		debugLog(string(out))
	}
	if err != nil {
		errorLog("err: couldn't execute 'go run .': %v", err)
		return
	}
	// Wait for the command to complete
	debugLog("Waiting for application to complete")
	err = cmd.Wait()
	if err != nil {
		errorLog("err: couldn't wait for the application to complete: %v", err)
		return
	}
}
