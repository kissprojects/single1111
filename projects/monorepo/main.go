package main

import (
	"fmt"
	"github.com/kissprojects/single/comps/go/file"
	"github.com/kissprojects/single/comps/go/os"
	"log"
	"os/exec"
	"strings"
)

func main() {
	migrateFolderToMonorepo("git@github.com:linksoft-dev/sigeflex.git",
		"comps/golang/mailer",
		"comps/go")
}

func runCommand(command string) *exec.Cmd {
	os.ExecuteCommand(command)
	cmd := exec.Command(command)
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
	return cmd
}

// migrateFolderToMonorepo this function migrate some folder in some repository to new monorepo repository
func migrateFolderToMonorepo(sourceRepoUrl, sourceFolder, destFolder string) {
	// source repository
	repoName := getRepoName(sourceRepoUrl)
	sourceRepoFolder := fmt.Sprintf("/tmp/%s", repoName)
	if err := file.CreateDirIfNotExists(sourceRepoFolder); err != nil {

	}
	runCommand(fmt.Sprintf("git clone %s %s", sourceRepoUrl, sourceRepoFolder))
	runCommand(fmt.Sprintf("cd %s", sourceRepoFolder))
	runCommand("git branch -b monorepo")
	runCommand(fmt.Sprintf("git filter-branch --subdirectory-filter %s -- --all", sourceFolder))
	runCommand(fmt.Sprintf(`git commit -m "mono repo — moving %s"`, sourceFolder))
	runCommand(fmt.Sprintf(`git push`))
	//runCommand(fmt.Sprintf("git mv %s %s", sourceFolder, destFolder))

	// monorepo repository
	runCommand(fmt.Sprintf("cd %s", destFolder))
	runCommand("git branch -b monorepo")
	runCommand(fmt.Sprintf("git remote add %s %s", repoName, sourceRepoUrl))
	runCommand("git fetch")
	runCommand(fmt.Sprintf("git merge %s/monorepo --allow-unrelated-histories", repoName))

}

func getRepoName(url string) string {
	// Split the URL using "/" and ".git"
	parts := strings.Split(url, "/")
	lastPart := parts[len(parts)-1]
	repoName := strings.TrimSuffix(lastPart, ".git")
	return repoName
}
