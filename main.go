package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"strings"
	"sync"

	"github.com/xoom/stash"
)

// ExecuteCMD executes cmd command for inputed params.
func ExecuteCMD(usrName, usrPas, stashHost, prjKey, branch string, p stash.Repository, wg *sync.WaitGroup) {
	defer wg.Done()
	gitCmd := fmt.Sprintf("%s%s@%s/scm/%s/%s.git", stashHost[:8], usrName, stashHost[8:], strings.ToLower(prjKey), p.Slug)
	gc := exec.Command("git", "clone", "-b", branch, gitCmd)
	err := gc.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(gitCmd)
}

func main() {

	var usrName, usrPas, stashHost, prjKey, branch string
	var p stash.Repository
	var prjSet map[int]stash.Repository
	var wg sync.WaitGroup

	flag.StringVar(&stashHost, "bp", "https://some-stash-host.net", "Stash host")
	flag.StringVar(&usrName, "un", "name", "Stash user name")
	flag.StringVar(&usrPas, "up", "pass", "Stash user pass")
	flag.StringVar(&prjKey, "pk", "project", "Project Key")
	flag.StringVar(&branch, "br", "master", "repo branch")

	flag.Parse()

	bu, err := url.Parse(stashHost)
	if err != nil {
		log.Fatal(err)
	}

	cl := stash.NewClient(usrName, usrPas, bu)

	prjSet, err = cl.GetRepositories()

	if err != nil {
		log.Fatal(err)
	}

	for i := range prjSet {
		p = prjSet[i]
		if p.Project.Key != prjKey {
			continue
		}
		wg.Add(1)
		go ExecuteCMD(usrName, usrPas, stashHost, prjKey, branch, prjSet[i], &wg)
		//gitCmd = fmt.Sprintf("git clone https://%s@%s/scm/%s/%s.git", usrName, stashHost, prjKey, p.Slug)
		//println(gitCmd)

		//println(prjSet[i].Name, prjSet[i].Slug)
	}
	wg.Wait()
}
