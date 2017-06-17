package main

import (
	"flag"
	"log"
	"net/url"

	"fmt"
	"github.com/xoom/stash"
	"os/exec"
	"strings"
)

func main() {
	var usrName, usrPas, stashHost, prjKey string

	flag.StringVar(&stashHost, "bp", "https://some-stash-host.net", "Stash host")
	flag.StringVar(&usrName, "un", "name", "Stash user name")
	flag.StringVar(&usrPas, "up", "pass", "Stash user pass")
	flag.StringVar(&prjKey, "pk", "project", "Project Key")

	flag.Parse()

	bu, err := url.Parse(stashHost)
	if err != nil {
		log.Fatal(err)
	}
	cl := stash.NewClient(usrName, usrPas, bu)

	var prjSet map[int]stash.Repository
	prjSet, err = cl.GetRepositories()

	if err != nil {
		log.Fatal(err)
	}
	var gitCmd string
	var p stash.Repository
	var gc *exec.Cmd

	for i := range prjSet {
		p = prjSet[i]
		if p.Project.Key != prjKey {
			continue
		}
		//gitCmd = fmt.Sprintf("git clone https://%s@%s/scm/%s/%s.git", usrName, stashHost, prjKey, p.Slug)
		//println(gitCmd)

		gitCmd = fmt.Sprintf("%s%s@%s/scm/%s/%s.git", stashHost[:8], usrName, stashHost[8:], strings.ToLower(prjKey), p.Slug)
		println(gitCmd)
		gc = exec.Command("git", "clone", gitCmd)
		err = gc.Run()
		if err != nil {
			println(err.Error())
		}

		//println(prjSet[i].Name, prjSet[i].Slug)
	}
}
