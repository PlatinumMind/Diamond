package cpp

import (
	"os"
	"strings"
)

type Output struct {
	BinName string
	Outpath string
}

type CppExec struct {
	Src string
	Cxx string
	Includes []string
	CxxFlags []string
	Output Output
}


func MakeCppExec(cfg CppExec) string {
	cmd := generateCommand(cfg.Cxx, cfg.Src, cfg.Includes , cfg.CxxFlags, cfg.Output.BinName, cfg.Output.Outpath)
	if _, err := os.Stat(cfg.Output.Outpath); os.IsNotExist(err) {
		os.Mkdir(cfg.Output.Outpath, 0775)
	}
	return cmd
}

func generateCommand(conpiler string, source string, includes []string, flags []string, binName string, outPath string) string {
	return conpiler + " " + source + " " + strings.Join(includes, " ") + " " + strings.Join(flags, " ") + " -o " + outPath + "/" + binName
}