package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/vulcan-frame/vulcan-game/gamedata"
	"github.com/vulcan-frame/vulcan-game/vulcan/pkg/filewriter"
	"github.com/vulcan-frame/vulcan-util/xsync"
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
			_, _ = fmt.Fprint(os.Stderr, fmt.Sprintf("%s", msg))
			os.Exit(1)
		}
	}()

	flag.Parse()

	jsonBaseDir = path.Join(filewriter.BasePath(), jsonBaseDir)
	slog.Info("json directory", "dir", jsonBaseDir)

	gamedata.Load(jsonBaseDir)
}
