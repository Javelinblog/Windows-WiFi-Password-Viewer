package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
    fmt.Println("https://github.com/Javelinblog/Windows-WiFi-Password-Viewer\n")
	fmt.Println("\nRetrieving saved WiFi profiles...\n")


	// Execute the command to get the list of Wi-Fi profiles
	cmd := exec.Command("netsh", "wlan", "show", "profiles")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error retrieving WiFi profiles:", err)
		return
	}

	// Adjust the regular expression to handle exact spacing
	profileRegex := regexp.MustCompile(`\s{4}All User Profile\s{5}:\s(.*)`)
	matches := profileRegex.FindAllStringSubmatch(out.String(), -1)

	// Iterate through the matches to print SSIDs and their passwords
	for _, match := range matches {
		if len(match) > 1 {
			// Extract the SSID and remove only the carriage return characters
			ssid := match[1]
			ssid = strings.ReplaceAll(ssid, "\r", "") // Remove carriage returns
			fmt.Printf("SSID: \"%s\"\n", ssid)

			// Get the password for the SSID
			password, err := getPassword(ssid)
			if err != nil {
				fmt.Println("Error retrieving password for SSID", ssid, ":", err)
			} else {
				fmt.Println("Password:", password)
			}
			fmt.Println() // Print a blank line after each SSID
		}
	}

	// Wait for user input before exiting
	fmt.Println("\nOperation completed. Press Enter to exit...")
	fmt.Scanln()
}

// getPassword retrieves the WiFi password for the given SSID
func getPassword(ssid string) (string, error) {
	// Execute the command to get the WiFi profile details
	cmd := exec.Command("netsh", "wlan", "show", "profile", "name="+ssid, "key=clear")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Use a regular expression to extract the password
	passwordRegex := regexp.MustCompile(`Key Content\s*:\s*(.*)`)
	match := passwordRegex.FindStringSubmatch(out.String())
	if len(match) > 1 {
		// Extract and return the password, removing any carriage return characters
		password := match[1]
		password = strings.ReplaceAll(password, "\r", "")
		return password, nil
	}
	return "No password found", nil
}
