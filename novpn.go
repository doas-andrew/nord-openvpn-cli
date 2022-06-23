package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
)

type NordRecommendation struct {
	Hostname string
	Load int
	Technologies []NordProtocol
}

type NordProtocol struct {
	Identifier string
}

func main() {
	// Define and parse command line flags
	protoPtr := flag.String("proto", "tcp", "Connection protocol. Defaults to TCP.\n  Accepts: TCP, UDP")
	flag.Parse()

	// Validate user inputs
	proto := strings.ToLower(*protoPtr)
	if proto != "udp" && proto != "tcp" {
		fmt.Printf("\"%v\" is not a valid connection protocol. Use TCP or UDP.", proto)
		return
	}

	// Get recommended servers from NordVPN
	resp, err := http.Get("https://api.nordvpn.com/v1/servers/recommendations")
	if err != nil { log.Fatalln("Could not reach https://api.nordvpn.com/v1/servers/recommendations \n", err) }

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil { log.Fatalln(err) }

	results := []NordRecommendation{}
	err = json.Unmarshal(body, &results)
	if err != nil { log.Fatalln(err) }

	if len(results) == 0 {
		fmt.Println("NordVPN did not return any recommended servers. https://nordvpn.com/servers/tools/")
		return
	}

	// Sort by .Load ASC
	sort.Slice(results, func(i, j int) bool {
		return results[i].Load < results[j].Load
	})

	var oFileAbs = ""
	var oDir = os.Getenv("MY_NOVPN_DIR")
	if oDir == "" {
		oDir = "$HOME/.config/novpn"
	}
	oDir = os.ExpandEnv(oDir)

	for _, rec := range results {
		// Check if this recommendation supports the specified protocol
		supportsProto := false
		for _, tech := range rec.Technologies {
			if tech.Identifier == "openvpn_"+proto {
				supportsProto = true
				break
			}
		}
		// Skip it if the specified protocol is not supported
		if !supportsProto {
			continue
		}

		oFileAbs = path.Join(oDir, fmt.Sprint(rec.Hostname, ".", proto, ".ovpn"))
		
		if _, err := os.Stat(oFileAbs); os.IsNotExist(err) {
			// File does not exist, request it from NordVPN
			resp, err = http.Get(fmt.Sprint("https://downloads.nordcdn.com/configs/files/ovpn_", proto, "/servers/", rec.Hostname, ".", proto, ".ovpn"))
			if err != nil { log.Fatalln(err) }
			defer resp.Body.Close()
	
			// Create empty file
			newFile, err := os.Create(oFileAbs)
			if err != nil { log.Fatalln(err) }
			defer newFile.Close()
	
			// Copy temp file from NordVPN into local file
			_, err = io.Copy(newFile, resp.Body)
			if err != nil { log.Fatalln(err) }
			
			// File saved successfully, break loop
			fmt.Println("Downloaded ", oFileAbs)
			break
		} else if err != nil {
			// Some other error has occurred
			log.Fatalln(err)
		} else {
			// We already have the file, break loop
			fmt.Println("Found ", oFileAbs)
			break
		}
	}

	if oFileAbs == "" {
		fmt.Println("Could not find server that supports openvpn_"+proto)
		return
	}

	// Connect with config file
	cmd := exec.Command("sudo", "openvpn", "--config", oFileAbs, "--auth-user-pass", path.Join(oDir, "up.txt"), "--auth-nocache")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil { log.Fatalln(err) }
}
