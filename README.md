# mhgu-switch-save-data-editor-go
MHGU Switch Save Data Editor (Go)
A Work-in-Progress command-line tool for editing Monster Hunter Generations Ultimate (MHGU) save data for both Switch and 3DS versions.

https://img.shields.io/badge/Go-1.21%252B-blue https://img.shields.io/badge/Status-WIP-yellow https://img.shields.io/badge/License-MIT-green

📋 Overview
This is a command-line application written in Go for editing save files from Monster Hunter Generations Ultimate (released as Monster Hunter XX in Japan). The tool aims to provide functionality for both the Nintendo Switch and Nintendo 3DS versions of the game.

⚠️ Important Notice: This project is currently under active development (WIP). Features are incomplete, and the structure is subject to change.

✨ Features (Planned)
Cross-Platform Support: Edit save files from both the Nintendo Switch and Nintendo 3DS versions.

Item/Equipment Editing: Modify items, weapons, armor, and character resources (Zenny, Hunter Points, etc.).

Character & Palico Stats: Adjust hunter rank, playtime, Palico details, and other character data.

Quest & Guild Card Editing: Modify quest completion status and guild card information.

Command-Line Interface (CLI): Easy-to-use terminal commands for quick save file operations.

Save Data Backup: Automatic backup of original save files before modification.

🚀 Getting Started
Prerequisites
Go 1.21 or higher: Download and install Go

A Monster Hunter Generations Ultimate (MHGU) save file.

Switch: You will need to extract your save file from the Nintendo Switch.

3DS: You will need to extract your save file from the 3DS or an emulator.

Installation
Clone the repository:

bash
git clone https://github.com/nurulhuda-git/mhgu-switch-save-data-editor-go.git
cd mhgu-switch-save-data-editor-go
Build the application:

bash
go build -o mhgu-editor main.go
This will create an executable named mhgu-editor (or mhgu-editor.exe on Windows) in the project directory.