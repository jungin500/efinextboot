package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type BcdeditEntry struct {
	name        string
	identifier  string
	description string
}

func checkCodepage() int {
	chcpRawOut, _ := exec.Command("cmd.exe", "/c", "chcp").Output()
	chcpStrings := strings.Split(string(chcpRawOut), " ")
	chcp, _ := strconv.Atoi(strings.TrimSpace(chcpStrings[len(chcpStrings)-1]))
	return chcp
}

func parseBytes(byteItem []byte) string {
	// Check current codepage of console
	cpid := checkCodepage()

	if cpid != 65001 {
		// Convert codepage of current window
		exec.Command("cmd.exe", "/c", "chcp 65001").Run()
	}

	return string(byteItem)
}

func parseBcdeditEntries(stdout []byte) []BcdeditEntry {
	decoded := parseBytes(stdout)
	lines := strings.Split(decoded, "\r\n")

	var entries []BcdeditEntry
	var currentEntry *BcdeditEntry = nil

	var isCurrentEntryTitle bool = false
	// Explicitly skip first line (is empty!)
	for idx := 1; idx < len(lines); idx++ {
		if currentEntry == nil {
			currentEntry = new(BcdeditEntry)
		}

		// Assume first and Last line always has empty line!
		if len(strings.TrimSpace(lines[idx])) == 0 {
			// Skip {fwbootmgr} as it doesn't have any effects
			if currentEntry.identifier != "{fwbootmgr}" {
				entries = append(entries, *currentEntry)
			}
			currentEntry = nil
			isCurrentEntryTitle = true
			continue
		}
		if isCurrentEntryTitle {
			currentEntry.name = strings.TrimSpace(lines[idx])
			isCurrentEntryTitle = false
			continue
		}
		if strings.HasPrefix(lines[idx], "identifier") {
			currentEntry.identifier = strings.TrimSpace(strings.TrimPrefix(lines[idx], "identifier"))
			continue
		}
		if strings.HasPrefix(lines[idx], "description") {
			currentEntry.description = strings.TrimSpace(strings.TrimPrefix(lines[idx], "description"))
			continue
		}
		if strings.HasPrefix(lines[idx], "description") {
			currentEntry.description = strings.TrimSpace(strings.TrimPrefix(lines[idx], "description"))
			continue
		}
	}

	// Append last element if left
	if currentEntry != nil {
		entries = append(entries, *currentEntry)
		currentEntry = nil
	}

	return entries
}

func chooseEntry(entries []BcdeditEntry, verbose bool) *BcdeditEntry {
	fmt.Println("[EFI Boot Entries]")
	for idx, entry := range entries {
		if verbose {
			fmt.Printf("%-02d: %s\n  - Name: %s\n  - Identifier: %s\n", idx, entry.description, entry.name, entry.identifier)
		} else {
			fmt.Printf("%2d: %s\n", idx, entry.description)
		}
	}
	fmt.Print("\nChoose nextboot entry by index: ")

	var nextIdxStr string
	var nextIdx int
	for true {
		_, err := fmt.Scanln(&nextIdxStr)
		if err != nil {
			fmt.Printf("Invalid selection! try again: ")
			continue
		}

		nextIdx, err = strconv.Atoi(nextIdxStr)
		if err != nil || !(0 <= nextIdx && nextIdx < len(entries)) {
			fmt.Printf("Invalid selection! try again: ")
		} else {
			break
		}
	}

	return &entries[nextIdx]
}

func updateBcdedit(entry *BcdeditEntry) {
	cmd := exec.Command("bcdedit.exe", "/set", "{fwbootmgr}", "bootsequence", entry.identifier)
	output, err := cmd.Output()

	if err != nil {
		fmt.Printf("Failed to set entry: %s\n", err.Error())
	} else {
		fmt.Printf("\nbcdedit: %s\n", parseBytes(output))
	}
}

func main() {
	// Check privilege error
	if exec.Command("bcdedit.exe").Run() != nil {
		fmt.Println("Privilege Error")
		os.Exit(1)
	}

	// Argument parsing
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "Verbose entry list output")
	flag.Parse()

	// Check boot sequences
	cmd := exec.Command("bcdedit.exe", "/enum", "firmware")
	output, _ := cmd.Output()

	entries := parseBcdeditEntries(output)
	selecedEntry := chooseEntry(entries, verbose)
	fmt.Printf("Updating nextboot to \"%s\"\n", selecedEntry.description)

	updateBcdedit(selecedEntry)
	fmt.Printf("Press any key to reboot now (Cancel with Ctrl+C) ...")
	fmt.Scanln()

	exec.Command("shutdown.exe", "/r", "/t", "0").Run()
	fmt.Printf("Rebooting ... (Do manually if doesn't)")
	fmt.Scanln() // Make terminal window persist
}
