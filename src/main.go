package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"diamond.build/src/langs/cpp"
	"github.com/pelletier/go-toml/v2"
)

type Target struct {
	Cpp_exe []cpp.CppExec
}

func main() {
	cfg := RunDiamond()
	// run cpp execs
	for _, v := range cfg.Cpp_exe {
		cmd := cpp.MakeCppExec(v)
		println(cmd)
		if runtime.GOOS == "windows" {
			o := exec.Command("cmd", "/c", cmd)

			var error bytes.Buffer
			o.Stderr = &error
			var out bytes.Buffer
			o.Stdout = &out

			err := o.Run()
			if err != nil {
				fmt.Println("Build failed:", error.String())
				return
			}

			println("Build Compelete with 0 error", out.String())
		} else {
			o := exec.Command("sh", "-c", cmd)

			var error bytes.Buffer
			o.Stderr = &error
			var out bytes.Buffer
			o.Stdout = &out

			err := o.Run()
			if err != nil {
				fmt.Println("Build failed:", error.String())
				return
			}

			println("Build Compelete with 0 error", out.String())
		}
	}
}

func RunDiamond() Target {
	var targets Target
	contents, err := os.ReadFile("./diamond.build")
	if err != nil {
		fmt.Println(err)
	}
	toml.Unmarshal(contents, &targets)
	return targets
}
