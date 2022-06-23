## About
This command line utility automatically downloads OpenVPN configuration files from NordVPN and uses them to connect with the OpenVPN protocol.

## Requirements
Ensure that you have the following installed and they're on your $PATH:
- OpenVPN: `which openvpn`
- Golang: `which go`

## Installation
- Clone the project: `git clone 'https://github.com/doas-andrew/nord-openvpn-cli.git'`
- Navigate into the project directory: `cd nord-openvpn-cli`
- Build the binary: `go build novpn.go`
- Move the binary to a directory that is on your $PATH: `mv novpn ~/.local/bin/novpn`
- Create the config directory: `mkdir -p ~/.config/novpn`
- Create the file used to hold your credentials: `touch $HOME/.config/novpn/up.txt`
- Find your "service credentials" on NordVPN
- Put the username on the first line of up.txt
- Put the password on the second line of up.txt
- Re-start or re-source your shell

## Usage
- `novpn -proto TCP`
- `novpn -proto UDP`
- `novpn` (defaults to TCP)
