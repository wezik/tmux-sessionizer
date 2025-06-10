# ⚡ Thop - Tmux Hopper
Fast and lightweight interactive CLI for defining and jumping between projects / tmux sessions.

## About
Light and quick to use way of managing tmux sessions

### Why not tmuxinator?
Tmuxinator is a great tool, but I found it to be too troublesome and too complex for my own need.
Thop is designed to be lightweight, simple to install, and extremely quick to use.

### Features:
- Fast navigation to desired project / session from anywhere (including from inside of a Tmux session)
- Easy to edit yaml templates
- Run commands in all/desired windows/panes

## Dependencies
- [fzf](https://github.com/junegunn/fzf)
- [tmux](https://github.com/tmux/tmux) 1.8+ (except for 2.5)

## Installation
WIP - no released binary yet, for now install from source (requires GO 1.24.2+):

```bash
git clone https://github.com/wezik/thop.git
cd thop
go build -o ./thop
sudo mv ./thop /usr/local/bin/
```

## Usage
```
thop <command> [args...]
```

### Commands:
```
select [name]           Open a project in a new session, or switch to an existing one.
                        - If no [name] is given, interactive selection is launched.
                        (aliases: s, <no args>)

create [name] [cwd]     Create a new project template.
                        - If only [name] is given, uses current working directory.
                        - If neither are given, both default to current directory name.
                        (aliases: c, a, add, append, new)

edit [name]             Edit a project template.
                        - If no [name] is given, interactive selection is launched.
                        - Uses editor defined in config, or $EDITOR env variable if not set.
                        (aliases: e)

delete [name]           Delete a project template.
                        - If no [name] is given, interactive selection is launched.
                        (aliases: d)

help                    Show this help message.
                        (aliases: any non-command string)

config                  Edit thop configuration
                        - Uses editor defined in config, or $EDITOR env variable if not set

Notes:
- Interactive selection is powered by fzf for commands without explicit [name].
```

## Editor

Thop uses your shell's default editor for opening files (`$EDITOR`)
To change it you have 2 options:

A) Add this to your `.bashrc`  

```bash
export EDITOR='vim'
``` 

B) Run below command, add `editor: vim` in a new line and save it

```bash
thop config
```

## Current state
This project is in a very early experimental stage, I am figuring out how I want this thing to work and what it should be capable of.

## Ideas / TODO's
- Refactor tmux unit tests to testify
- Panes support
- Integration tests
- General config file, launchable from command
- List active tmux sessions, and allow to attach to them
- Consider replacing cli entrypoint with cobra
- Add showcase and example template to README
