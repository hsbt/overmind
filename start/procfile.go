package start

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/DarthSim/overmind/utils"
)

type procfileEntry struct {
	Name    string
	Command string
}

type procfile []procfileEntry

func parseProcfile(procfile string, portBase, portStep int) (pf procfile) {
	re, _ := regexp.Compile("^(\\w+):\\s+(.+)$")

	f, err := os.Open(procfile)
	utils.FatalOnErr(err)

	port := portBase
	names := make(map[string]bool)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if len(scanner.Text()) > 0 {
			params := re.FindStringSubmatch(scanner.Text())
			if len(params) != 3 {
				continue
			}

			name, cmd := params[1], params[2]

			if names[name] {
				utils.Fatal("Process names must be uniq")
			}
			names[name] = true

			pf = append(pf, procfileEntry{
				name,
				strings.Replace(cmd, "$PORT", strconv.Itoa(port), -1),
			})

			port += portStep
		}
	}

	utils.FatalOnErr(scanner.Err())

	if len(pf) == 0 {
		utils.Fatal("No entries was found in Procfile")
	}

	return
}

func (p procfile) MaxNameLength() (nl int) {
	for _, e := range p {
		if l := len(e.Name); nl < l {
			nl = l
		}
	}
	return
}
