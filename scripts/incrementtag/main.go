package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	tagsOutput, err := exec.Command("git", "tag").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	tags := strings.Split(strings.TrimSpace(string(tagsOutput)), "\n")
	if len(tags) == 0 {
		log.Fatal("no tags found")
	}

	tag := tags[len(tags)-1]
	parts := strings.Split(tag, ".")
	if len(parts) != 3 {
		log.Fatalf("malformed tag: %v", tag)
	}

	// If the last part of the tag had something like "-alpha", chop it off.
	lastPartStr, _, _ := strings.Cut(parts[2], "-")

	num, err := strconv.Atoi(lastPartStr)
	if err != nil {
		log.Fatal(err)
	}

	newNum := num + 1

	// NOTE: parts[0] still has the "v" prefix, which we need for Go.
	fmt.Printf("%s.%s.%d\n", parts[0], parts[1], newNum)
}
