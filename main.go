package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	gittools "github.com/sebastianRau/deployer/pkg/gitTools"
	ostools "github.com/sebastianRau/deployer/pkg/osTools"
	"github.com/sebastianRau/deployer/pkg/steps"
)

var (
	version = "1.0.0"
)

//go:embed keys/*
var keyFiles embed.FS

func main() {

	const (
		configTemplateFile string = "./temp/KW_BEM_missions/kw_bem_make.json.tpl"
		configDataFile     string = "./temp/KW_BEM_missions/kw_bem_make.data.json"
	)

	var (
		verbose      = flag.Bool("v", false, "verbose output")
		versionPrint = flag.Bool("version", false, "print version")
	)

	var (
		err   error
		check bool
	)

	flag.Parse()

	if *versionPrint {
		fmt.Printf("KW Mission Builder %s\n", version)
		os.Exit(0)

	}
	fmt.Printf("KW Mission Builder\n")

	check, err = CheckKnownHosts("github.com")
	if err != nil {
		panic(err)
	}
	if !check {
		fmt.Println("github.com is no present in you known hosts file!\nPlease call\n     ssh -T git@github.com \nto add it to knows hosts")
		os.Exit(2)
	}

	defer removeTemp()

	fmt.Println()
	id_BasicMissions, err := keyFiles.ReadFile("keys/id_BasicMissions")
	if err != nil {
		panic("No Keyfile found for id_BasicMissions")
	}
	id_BasicScripts, err := keyFiles.ReadFile("keys/id_BasicScripts")
	if err != nil {
		panic("No Keyfile found for id_BasicScripts")
	}

	_, _, err = gittools.UpdateKeyBytes(
		"git@github.com:SebastianRau/KW_BEM_BasicScripts.git",
		"./temp/KW_BEM_BasicScripts",
		id_BasicScripts,
		"",
		nil,
	)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("%-*s %s\n", 80, "Update KW_BEM_BasicScripts", "OK")
	}

	_, _, err = gittools.UpdateKeyBytes(

		"git@github.com:SebastianRau/KW_BEM_missions.git",
		"./temp/KW_BEM_missions",
		id_BasicMissions,
		"",
		nil,
	)

	if err != nil {
		panic(err.Error())
	} else {

		fmt.Printf("%-*s %s\n", 80, "Update KW_BEM_missions", "OK")
	}

	// read and unmarshal template
	st, err := steps.UnmarshalConfigTemplate(configTemplateFile, configDataFile)

	if err != nil {
		removeTemp()
		fmt.Println(err)
		os.Exit(1)
	}

	start := time.Now()
	st.Exceute(os.Stdout, *verbose)
	stop := time.Now()
	elapsed := stop.Sub(start)

	fmt.Printf("Time: %s\n", elapsed.Round(time.Millisecond))
	removeTemp()
	os.Exit(0)

}

func removeTemp() {
	err := ostools.Delete(
		"./temp/",
		nil,
	)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("%-*s %s\n", 80, "Removing tempfolder", "OK")
	}
}

func CheckKnownHosts(url string) (bool, error) {

	homeDir, err := os.UserHomeDir()
	knownHosts := homeDir + "/.ssh/known_hosts"
	bytes, err := os.ReadFile(knownHosts)
	if err != nil {
		return false, err
	}

	return strings.Contains(string(bytes), url), nil
}
