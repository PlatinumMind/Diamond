package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type Step struct {
	Name string
	Cmd  string
}

type Pre struct {
	Cmd string
}

type Clean struct {
	Cmd string
}

type Recipient struct {
	Vars  map[string]string
	Pre   Pre
	Step  []Step
	Clean Clean
}

func main() {
	var cfg Recipient
	content, err := os.ReadFile("diamond.build")
	if err != nil {
		panic(err)
	}
	toml.Unmarshal(content, &cfg)
	fmt.Printf("%v\n", cfg.Vars["test1"])
	run(cfg)
}

func run(cfg Recipient) {

	if cfg.Pre.Cmd != "" {
		println("preparing build: ")
		cmd := Vars(cfg.Pre.Cmd, cfg.Vars)
		cmd = strings.ReplaceAll(cmd, "'", "\\'")
		cmd = strings.ReplaceAll(cmd, "\"", "\\\"")
		if runtime.GOOS == "windows" {
			o, err := exec.Command("cmd", "/c", cmd).Output()
			if err != nil {
				log.Fatal(err)
			}
			println(string(o))
		} else {
			o, err := exec.Command("sh", "-c", cmd).Output()
			if err != nil {
				log.Fatal(err)
			}
			println(string(o))
		}
	}
	for _, step := range cfg.Step {
		print("[" + step.Name + "] started running\n")
		cmd := Vars(step.Cmd, cfg.Vars)
		cmd = strings.ReplaceAll(cmd, "'", "\\'")
		cmd = strings.ReplaceAll(cmd, "\"", "\\\"")
		if runtime.GOOS == "windows" {
			o, err := exec.Command("cmd", "/c", cmd).Output()
			if err != nil {
				log.Fatal(err)
			}
			println(string(o))
		} else {
			o, err := exec.Command("sh", "-c", cmd).Output()
			if err != nil {
				log.Fatal(err)
			}
			println(string(o))
		}
		println("[" + step.Name + "] finished running\n")
	}

	if cfg.Clean.Cmd != "" {
		println("cleaning build:")
		cmd := Vars(cfg.Clean.Cmd, cfg.Vars)
		cmd = strings.ReplaceAll(cmd, "'", "\\'")
		cmd = strings.ReplaceAll(cmd, "\"", "\\\"")
		if runtime.GOOS == "windows" {
			o, err := exec.Command("cmd", "/c", cmd).Output()
			if err != nil {
				log.Fatal(err)
			}
			println(string(o))
		} else {
			o, err := exec.Command("sh", "-c", cmd).Output()
			if err != nil {
				log.Fatal(err)
			}
			println(string(o))
		}
	}
}

func Vars(input string, data map[string]string) string {
	m1 := regexp.MustCompile(`{(.*?)}`)
	m2 := regexp.MustCompile(`{|}`)
	return m1.ReplaceAllStringFunc(input, func(s string) string {
		return data[m2.Split(s, 10)[1]]
	})
}
