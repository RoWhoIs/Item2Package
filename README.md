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