package main

import (
 	"errors"
	"os"
	"os/exec"
)


type Workspace struct {
	base	string
	name	string
}


func New(basedir, name string) Workspace {
	return Workspace{basedir,name}
}


func (w *Workspace) Create() error {

	dir := w.String() + "/"
	if err := os.Mkdir(w.String(),0755); nil!=err {
		return errors.New("mkdir "+dir+": "+err.Error())
	}
	
	sub := []string{"src","pkg","bin"}
	
	for _,v := range sub {
		
		sd := dir + v
		if err := os.Mkdir(sd,0755); nil!=err {
			return errors.New("mkdir "+sd+": "+err.Error())
		}
	}
	
	return nil
}

func (w *Workspace) Purge() error {

	// v1: NO-OP
	return nil
}


func (w *Workspace) Enter() error {
	return wsenter(w,"")
}

func (w *Workspace) EnterX(goPath string) error {
	return wsenter(w,goPath)
}

func wsenter(w *Workspace, goPath string) error {
		
	env := os.Environ()
	goEnv := make([]string,0,len(env)+2)	
	
	gp := w.String()
	if ""!=goPath {
		gp = gp+":"+goPath
	}
	goEnv=append(env,"GOPATH="+gp)
	goEnv=append(goEnv,("VIRTUALGO="+w.name))
	
	cmd := exec.Command("/bin/bash","-i")
	cmd.Env = goEnv
	cmd.Dir = w.String()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}



func (w *Workspace) Name() string {
	return w.name
}
func (w *Workspace) Base() string {
	return w.base
}

func (w *Workspace) String() string {
	return w.base+"/"+w.name
}

