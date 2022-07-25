package main

import (
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
	node "github.com/tidwall/go-node"
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
	run(cfg)
}

func run(cfg Recipient) {
	vm := node.New(nil)
	if cfg.Pre.Cmd != "" {
		println("preparing build")
		
		v := vm.Run("child_process.execSync('"+cfg.Pre.Cmd+"')")
		if err := v.Error(); err != nil {
			log.Fatal(err)
		}
		println(v.String())
	}
	for _, step := range cfg.Step {
		print("[" + step.Name + "] started running\n")
		v := vm.Run("child_process.execSync('"+step.Cmd+"')")
		if err := v.Error(); err != nil {
			log.Fatal(err)
		}
		println(v.String())
		println("[" + step.Name + "] finished running\n")
	}

	if cfg.Clean.Cmd != "" {
		println("cleaning build:")
		v := vm.Run("child_process.execSync('"+cfg.Pre.Cmd+"')")
		if err := v.Error(); err != nil {
			log.Fatal(err)
		}
		println(v.String())
	}
}
