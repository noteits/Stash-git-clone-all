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
func ExecuteCMD(usrName, usrPas, stashHost, prjKey, branch string, p stash.Repository, output chan string, wg *sync.WaitGroup, err error) {
	defer wg.Done()
	gitCmd := fmt.Sprintf("%s%s@%s/scm/%s/%s.git", stashHost[:8], usrName, stashHost[8:], strings.ToLower(prjKey), p.Slug)
	var gc = *exec.Command("git", "clone", "-b", branch, gitCmd)
	err = gc.Run()
	if err != nil {
		output <- err.Error()
	}
	output <- gitCmd
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

	response := make(chan string)

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
			delete(prjSet, i)
		}
		//gitCmd = fmt.Sprintf("git clone https://%s@%s/scm/%s/%s.git", usrName, stashHost, prjKey, p.Slug)
		//println(gitCmd)

		//println(prjSet[i].Name, prjSet[i].Slug)
	}

	(&wg).Add(len(prjSet))

	for i := range prjSet {
		go ExecuteCMD(usrName, usrPas, stashHost, prjKey, branch, prjSet[i], response, &wg, err)
	}

	go func() {
		for i := range response {
			fmt.Println(i)
		}
	}()

	wg.Wait()
}
