## About
This command line utility automatically downloads OpenVPN configuration files from NordVPN and connects using the file.

## Installation
- Ensure you have `openvpn` installed and it's on your $PATH (check with `which openvpn`)
- `git clone` this repo
- `cd` into the directory
- `go build novpn.go`
- Place the binary in a directory that is on your $PATH (~/.local/bin)
- `mkdir $HOME/ovpn`
- `touch $HOME/ovpn/up.txt`
- Find your "service credentials" on NordVPN
- Put the username on the first line of up.txt
- Put the password on the second line of up.txt
- Re-start or re-source your shell

## Usage
- `novpn -proto TCP`
- `novpn -proto UDP`
- `novpn` (defaults to TCP)
