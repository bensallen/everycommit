package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type pullRequest struct {
	owner   string
	project string
	id      int
}

func parseURL(rawurl string) (*pullRequest, error) {
	URL, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if URL.Hostname() != "github.com" {
		return nil, fmt.Errorf("Only Github.com pull requests are supported")
	}
	elements := strings.Split(URL.Path, "/")

	if !(len(elements) == 5 && elements[3] == "pull") {
		return nil, fmt.Errorf("URL doesn't match expected format, eg. https://github.com/<owner>/<project>/pull/<ID>")
	}

	id, err := strconv.Atoi(elements[4])
	if err != nil {
		return nil, fmt.Errorf("URL doesn't have a numeric pull request ID")
	}

	return &pullRequest{owner: elements[1], project: elements[2], id: id}, nil
}

func (pr *pullRequest) commits(ctx context.Context, client *github.Client) ([]*github.RepositoryCommit, error) {
	s := client.PullRequests
	commits, _, err := s.ListCommits(ctx, pr.owner, pr.project, pr.id, &github.ListOptions{})
	return commits, err
}

func openRepo(repopath string) (*git.Worktree, error) {
	repo, err := git.PlainOpen(repopath)
	if err != nil {
		return nil, err
	}
	return repo.Worktree()
}

func run(cmd *exec.Cmd, wt *git.Worktree, commit *github.RepositoryCommit, out io.Writer) error {

	fmt.Fprintf(out, "Checking out commit: %s\n", commit.GetSHA())

	err := wt.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commit.GetSHA()),
	})
	if err != nil {
		return fmt.Errorf("git checkout, %s\n", err)
	}

	fmt.Fprintf(out, "Running: %s %s\n", cmd.Path, cmd.Args)

	return cmd.Run()
}

func main() {
	var (
		help     bool
		url      string
		repopath string
	)

	flag.BoolVar(&help, "h", false, "Display this help screen and quit")
	flag.StringVar(&url, "u", "", "Pull request URL, eg: https://github.com/<owner>/<project>/pull/<ID>")
	flag.StringVar(&repopath, "r", "", "Directory of cloned repository")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [options] cmd [args...]:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()

	if help {
		flag.Usage()
		os.Exit(2)
	}

	if url == "" || repopath == "" {
		flag.Usage()
		os.Exit(2)
	}

	wt, err := openRepo(repopath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	pr, err := parseURL(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := github.NewClient(nil)
	commits, err := pr.commits(ctx, client)

	for _, commit := range commits {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = repopath

		if err := run(cmd, wt, commit, os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "error, %s\n", err)
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					os.Exit(status.ExitStatus())
				}
			}
			os.Exit(1)
		}
	}
}
