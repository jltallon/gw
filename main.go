// gw is a Golang workspace switcher
package main

import (
	"errors"
	"github.com/doblenet/go-doblenet/tracer"
// 	home "github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"path/filepath"
	"os"
 	"os/user"
// 	"syscall"
)

const (
	MODE_UNKNOWN = iota
	MODE_SWITCH
	MODE_CREATE
	MODE_PURGE
	MODE_SETUP	= 0xFF
)

var (
	mode	int = MODE_UNKNOWN
	force	bool
		
	branch,create,purge	bool
)


func main() {

	setupFlags()
	pflag.Usage = func() {
		tracer.Putln("gw - Go workspace switcher")
		tracer.Putln(`Usage:
	gw [-command] <workspace>
`)
		pflag.PrintDefaults()
	}
	
	pflag.Parse()
		
	if pflag.NArg() != 1 {
		tracer.Fatal("Required argument: <workspace>")
		os.Exit(1)
	}
	
	args := pflag.Args()
	modeConv(args)


	homeDir,err := getHome()
	if nil!=err {
 		tracer.FatalExit(2,"go-homedir: "+err.Error())
 	}
	
	w := New(homeDir,args[0])
	
	
	switch mode {
	case MODE_CREATE:
		if err:=w.Create(); nil!=err {
			tracer.FatalExit(3,err.Error())
		}
		
	case MODE_PURGE:
		if err:=w.Purge(); nil!=err {
			tracer.FatalExit(3,err.Error())
		}
	
	case MODE_SWITCH:
		fallthrough
		
	default:
		var gp string
		if gp,err=goPath(); nil!=err {
			tracer.Warn(err.Error())
		}
		if err:=w.EnterX(gp); nil!=err {
			tracer.FatalExit(3,err.Error())
		}
	}

	
	os.Exit(0)
}


func setupFlags() {
	
	pflag.BoolVarP(&force,"force","f",false,"Force creation/overwrite/delete.")
	
	pflag.BoolVarP(&create,"create","c",false,"Create a new workspace")
	pflag.BoolVarP(&branch,"branch","b",false,"Change into a new workspace (default)")
	pflag.BoolVarP(&purge, "purge","",false,"Remove a workspace")
	
	return
}

func modeConv(args []string) int {
	
	switch {
	case branch:
		mode = MODE_SWITCH
		
	case create:
		mode = MODE_CREATE
	
	case purge:
		mode = MODE_PURGE
		
	default:
	}
	return mode
}



func getHome() (string,error) {
	u,err := user.Current()
	if nil!=err {
		return "~/",err
	}
	return u.HomeDir,nil
}


func goPath() (string,error) {
	
	if ge := os.Getenv("GOPATH"); ""!=ge {
		return ge,nil
	}
	
	var home string
	var err error
	
	if home,err = getHome(); nil!=err {
		return "",err
	}
	
	home,err = filepath.EvalSymlinks(home+"/go")
	if nil==err {
		return home,nil
	}
	
	return "",errors.New("Could not resolve GOPATH")
}
