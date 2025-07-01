package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/go-pantheon/roma/gamedata"
	"github.com/go-pantheon/roma/vulcan/pkg/filewriter"
)

var (
	jsonBaseDir string
	fullStderr  bool
)

func init() {
	flag.StringVar(&jsonBaseDir, "json_dir", "gen/parser/json", "json dir path, eg: -json_dir gen/parser/json")
	flag.BoolVar(&fullStderr, "full_std_err", false, "")
}

func main() {
	defer func() {
		if p := recover(); p != nil {
			msg := xsync.CatchErr(p).Error()
			if !fullStderr {
				lines := strings.Split(msg, "\n")
				if len(lines) > 1 {
					msg = lines[0]
				}
			}

			_, _ = fmt.Fprintf(os.Stderr, "%s\n", msg)
			os.Exit(1)
		}
	}()

	flag.Parse()

	jsonBaseDir = path.Join(filewriter.BasePath(), jsonBaseDir)
	slog.Info("json directory", "dir", filewriter.SprintGenPath(jsonBaseDir))

	gamedata.Load(jsonBaseDir)
}
