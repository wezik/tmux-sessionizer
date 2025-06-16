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
Run below script to install the latest release:

```bash
curl -L https://github.com/wezik/thop/releases/latest/download/thop -o ./thop
chmod +x ./thop
sudo mv ./thop /usr/local/bin/
```

## Usage
```bash
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
To change it you best option is to override it in your `.bashrc` or other shell config file:

```bashrc
export EDITOR='vim'
```

## Current state
This project is in a somewhat early experimental stage, it's destination is set but things can still change.

## Ideas / TODO's
- Panes support
- Integration tests
- General config file, launchable from command
- List active tmux sessions, and add them to configured templates (for easy switching)
- Add showcase and example template to README
- Review the Makefile
