package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/peterh/liner"
)

var diceRegexp = regexp.MustCompile(`^(\d+)d(\d+) *((-|\+) *(\d+))?$`)

func main() {
	rand.Seed(time.Now().UnixNano())

	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	for {
		str, err := line.Prompt("> ")
		if err == liner.ErrPromptAborted {
			break
		} else if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		} else if str == "" {
			continue
		}

		submatch := diceRegexp.FindStringSubmatch(str)
		if submatch == nil {
			fmt.Fprintln(os.Stderr, "invalid syntax.")
			continue
		}
		dice, err := strconv.ParseInt(submatch[1], 10, 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		sides, err := strconv.ParseInt(submatch[2], 10, 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		var total int64
		if submatch[3] != "" {
			bonus, err := strconv.ParseInt(submatch[5], 10, 0)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			if submatch[4] == "+" {
				total += bonus
			} else {
				total -= bonus
			}
		}
		for i := 0; i < int(dice); i++ {
			total += int64(rand.Intn(int(sides)) + 1)
		}

		fmt.Println(total)
		line.AppendHistory(str)
	}
}
