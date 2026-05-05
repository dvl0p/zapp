# ⚡ Zapp

A lightweight CLI tool written in Go for managing network ports and the processes using them.
Say goodbye to hunting down PIDs just to free up `localhost`.

## Features
- Quickly see all `LISTENING` ports
- Free up blocked port with a single command

## Installation
Make sure you have [Go](https://github.com/golang/go) installed, then:
```Bash
go install github.com/dvl0p/zapp@latest
```

## Usage

**List all active listening ports**
```Bash
zapp
```

**Free a specific port (kills underlying process)**
```Bash
# For running on port PORT:
# zapp free PORT
# For example on port 8080:
zapp free 8080
```

