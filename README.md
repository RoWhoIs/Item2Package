# Item2Package

## Description

Item2Package is a simple utility that compiles Roblox assets into .zip files without any user account needed. This tool is designed to streamline the process of packaging Roblox assets, making compiling and using assets in programs like Novetus and Blender easier.

## Installation

To install Item2Package, follow these steps:

1. Download the latest release from the GitHub repository.
2. Extract the .zip file to your desired location.
3. In a command prompt, run item2package (`./item2package [args]` or `item2package [args]`)

## Usage

| Flag | Type | Description                                             | Required | Default Value | Example      |
|------|------|---------------------------------------------------------|----------|---------------|--------------|
| `-i` | int  | The item you're wishing to package                      | Yes      | None          | `-i 1082932` |
| `-l` | bool | Legacy mode, converts the rbxm for use in older clients | No       | false         | `-l true`    |
| `-v` | bool | Verbosely log operations                                | No       | false         | `-v true`    |
| `-o` | dir  | Set the output directory for the zip file               | No       | .             | `-o ~`       |

> [!IMPORTANT]  
> The `-o` flag does not currently work and will continue to output files to the working directory

## FaQ

> This was flagged by my antivirus, is this software trustworthy?

This is because none of the packages are signed (that costs money for Windows and MacOS)
If you'd like, you can compile the source code yourself. We used `go build -ldflags "-s -w" -o item2package main.go` as our compilation command.

> This keeps getting a 429 error!

Wait for about one minute. This is caused due to none of the requests made by Item2Package being authenticated, so Roblox enforces stricter network policies.

> How do I contribute?

Fork this repository, create a new branch, and open a pull request. Thanks in advanced!
