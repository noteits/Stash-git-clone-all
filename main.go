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

func stashQueries(repos map[int]stash.Repository) <-chan stash.Repository {
	out := make(chan stash.Repository)
	go func() {
		defer close(out)
		for i := range repos {
			out <- repos[i]
		}
	}()
	return out
}

// ExecuteCMD executes cmd command for inputed params.
func ExecuteCMD(repoCh <-chan stash.Repository, usrName, usrPas, stashHost, branch, prjKey string, wg *sync.WaitGroup) {
	defer wg.Done()
	for repo := range repoCh {
		gitCmd := fmt.Sprintf("%s%s@%s/scm/%s/%s.git", stashHost[:8], usrName, stashHost[8:], strings.ToLower(prjKey), repo.Slug)
		gc := exec.Command("git", "clone", "-b", branch, gitCmd)
		err := gc.Run()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(gitCmd)
	}
}

func main() {

	var usrName, usrPas, stashHost, prjKey, branch string
	var rangeRepos int
	var prjSet map[int]stash.Repository
	var wg sync.WaitGroup

	flag.StringVar(&stashHost, "bp", "https://some-stash-host.net", "Stash host")
	flag.StringVar(&usrName, "un", "name", "Stash user name")
	flag.StringVar(&usrPas, "up", "pass", "Stash user pass")
	flag.StringVar(&prjKey, "pk", "project", "Project Key")
	flag.StringVar(&branch, "br", "master", "repo branch")
	flag.IntVar(&rangeRepos, "rr", 5, "Range of repos")

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

	reposStash := stashQueries(prjSet)

	wg.Add(rangeRepos)

	for i := 0; i < rangeRepos; i++ {
		go ExecuteCMD(reposStash, usrName, usrPas, stashHost, branch, prjKey, &wg)
	}
	wg.Wait()
}
