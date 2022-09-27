package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/microctar/licorice/app/constant"
)

var (
	server     bool
	confdir    string
	help       bool
	inputfile  string
	outputfile string
	port       int
	rule       string
	target     string
	version    bool
)

const (
	usage = "USAGE\n\t%s [flags]\n"
)

func init() {
	// ...(args, cmdargs, default_val, describtion)
	flag.BoolVar(&server, "api", true, "enable api mode")
	flag.StringVar(&confdir, "config", "", "look for configuration file at the path")
	flag.BoolVar(&help, "help", false, "Display this information")
	flag.StringVar(&inputfile, "input", "", "specify input file")
	flag.IntVar(&port, "port", 6060, "specify server port")
	flag.StringVar(&outputfile, "output", "config", "specify output file")
	flag.StringVar(&rule, "rule", "ACL4SSR.ini", "specify rule file")
	flag.StringVar(&target, "target", "clash", "specify target")
	flag.BoolVar(&version, "version", false, "print licorice version")
	flag.Usage = Usage
	flag.Parse()

}

func Usage() {
	fmt.Fprintf(flag.CommandLine.Output(), usage, "licorice")
	flag.PrintDefaults()
}

func main() {

	// print program version
	if version {
		fmt.Printf("licorice %s %s %s with %s %s\n", constant.Version, runtime.GOOS, runtime.GOARCH, runtime.Version(), constant.BuildTime)
		return
	}

	if help {
		flag.Usage()
		return
	}

	// run as commmand line tool
	if inputfile != "" {
		RunCMD()
		return
	}

	// run as server
	if server {
		RunServer()
		return
	}

}
