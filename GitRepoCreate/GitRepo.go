package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {

	// Define cmd line arguments
	var name string
	flag.StringVar(&name, "name", "", "Name of the git repo / project")
	flag.StringVar(&name, "n", "", "Name of the git repo / project")

	var username string
	flag.StringVar(&username, "username", "", "Specify git username")
	flag.StringVar(&username, "u", "", "Specify git username")

	var email string
	flag.StringVar(&email, "email", "", "Specify git email address")
	flag.StringVar(&email, "e", "", "Specify git email address")

	var sshkey string
	flag.StringVar(&sshkey, "key", "~/.ssh/id_rsa", "Specify full path to SSH key")
	flag.StringVar(&sshkey, "k", "~/.ssh/id_rsa", "Specify full path to SSH key")

	var checkVersion bool
	flag.BoolVar(&checkVersion,"version", false, "Display current version")
	flag.BoolVar(&checkVersion,"v", false, "Display current version")


	flag.Parse()

	if checkVersion {
		fmt.Println("GoGit v0.0.1")
		os.Exit(0)
	}

	if name != "" {
		gitInit := exec.Command("git", "init")
		gitInit.Run()

		gitSetBranch := exec.Command("git", "branch", "-M", "main")
		gitSetBranch.Run()

		f, err := os.Create("README.md")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		readmeString := "# " + name
		_, err2 := f.WriteString(readmeString)
		if err2 != nil {
			log.Fatal(err2)
		}

		// Create Github Repo - Private by default
		postBody, _ := json.Marshal(map[string]string{
			"name": name,
			"private": "true",
		})
		responseBody := bytes.NewBuffer(postBody)

		url := "https://api.github.com/user/repos"
		GOGIT := os.Getenv("GOGIT")
		var bearer = "Token " + GOGIT
		req, err := http.NewRequest("POST", url, responseBody)
		req.Header.Add("Authorization", bearer)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}
		defer resp.Body.Close()
	} else {
		fmt.Println("The name flag must be specified. Please use -h / -help for usage. Exiting")
		os.Exit(0)
	}

	if username != "" {
		gitUn := exec.Command("git","config", "user.name", username)
		gitUn.Run()
	}

	if email != "" {
		gitEmail := exec.Command("git", "config", "user.email", email)
		gitEmail.Run()
	}

	if sshkey != "" {
		sshCommand := "ssh -i " + sshkey + " -F /dev/null"
		gitSshKey := exec.Command("git", "config", "core.sshCommand", sshCommand)
		gitSshKey.Run()
	}

}
