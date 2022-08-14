package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type AppData struct {
	Name      string
	Version   string
	BuildDate string
}

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

type Instruction struct {
	Vars  map[string]string
	Pre   Pre
	Step  []Step
	Clean Clean
}

func main() {
	year, month, day := time.Now().Date()

	appData := AppData{
		Name:      "Diamond Build",
		Version:   "0.0.3",
		BuildDate: fmt.Sprintf("%s-%s-%s", strconv.Itoa(year), month, strconv.Itoa(day)),
	}
	if os.Args[1] == "" {
		var cfg Instruction
		content, err := os.ReadFile("diamond.build")
		if err != nil {
			fmt.Fprintln(io.Writer(os.Stderr), "No diamond.build file found")
		}
		toml.Unmarshal(content, &cfg)
		run(cfg)
	}
	if os.Args[1] == "--version" {
		versionData := fmt.Sprintf(`%s by Platinum Mind
Version: %s
Build Date: %s`, appData.Name, appData.Version, appData.BuildDate)
		println(versionData)

	}
}

func run(cfg Instruction) {

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
