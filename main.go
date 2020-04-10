// Package main provides all the functionality for GoThanks.
// GoThanks automatically stars Go's official repository and your go.mod github dependencies,
// providing a simple way to to say thanks to the maintainers of the modules you use and the contributors of Go itself.
//
// Usage:
//
// In order to run GoThanks you need to have a valid Github Access Token.
// You can pass the token as an argument to GoThanks or store it in an environmental variable named GITHUB_TOKEN.
//
// Inside the folder where your go.mod lives run:
//
// $ ./gothanks -github-token=xxxxxx
//
// or
//
// $ export GITHUB_TOKEN=xxxxx
// $ ./gothanks
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/v28/github"
	"github.com/sirkon/goproxy/gomod"
	"golang.org/x/oauth2"
)

func main() {
	githubToken := flag.String("github-token", os.Getenv("GITHUB_TOKEN"), "Github access token")
	assumeYes := flag.Bool("y", false, "Automatic yes to prompts and run non-interactively.")
	flag.Parse()

	if *githubToken == "" {
		fmt.Println("no Github token found")
		os.Exit(-1)
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	input, err := ioutil.ReadFile(dir + "/go.mod")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println("Welcome to GoThanks :)")

	// Asks user what he wants to do
	if *assumeYes == false {
		fmt.Println("\nYou are about to star you beloved dependencies.\n\nPress y to continue or n to abort")

		confirmed, err := confirm(os.Stdin)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		if !confirmed {
			fmt.Println("Aborting.")
			os.Exit(0)
		}
	}

	parseResult, err := gomod.Parse("", input)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	ctx := context.Background()
	client := githubClient(ctx, *githubToken)

	// Always send your love to Go!
	repos := []string{"github.com/golang/go"}
	for dep := range parseResult.Require {
		repos = append(repos, dep)
	}

	fmt.Print("\nSending your love..\n\n")
	for _, dep := range repos {

		if rep, ok := isGithubRepo(dep); ok {
			x, _, err := client.Activity.IsStarred(ctx, rep.owner, rep.repo)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if x {
				fmt.Printf("Repository %s is already starred!\n", rep.path)
				continue
			}

			fmt.Printf("Sending a star to %s\n", rep.path)

			_, err = client.Activity.Star(ctx, rep.owner, rep.repo)
			if err != nil {
				fmt.Printf("Could not star %s %s\n", rep.path, err)
			}
		}
	}

	fmt.Println("\nThank you!")
}

type githubRepo struct {
	path  string
	owner string
	repo  string
}

func githubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func confirm(in io.Reader) (bool, error) {
	reader := bufio.NewReader(in)

	char, _, err := reader.ReadRune()
	if err != nil {
		return false, err
	}

	if i := strings.ToLower(strings.TrimSpace(string(char))); i == "y" {
		return true, nil
	}

	return false, nil
}

func isGithubRepo(path string) (githubRepo, bool) {
	// Make sure we do not forget to star the Github mirrors of Go's subpackages
	path = strings.Replace(path, "golang.org/x/", "github.com/golang/", -1)

	re := regexp.MustCompile(`^github\.com\/[a-zA-Z\d-]+\/[a-zA-Z\d-]+`)
	repoPath := re.FindString(path)
	if repoPath != "" {
		parts := strings.Split(repoPath, "/")
		res := githubRepo{
			path:  repoPath,
			owner: parts[1],
			repo:  parts[2],
		}

		return res, true
	}

	return githubRepo{}, false
}
