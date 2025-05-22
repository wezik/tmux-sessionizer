# Tmux/Project sessionizer/hopper (phop)
_name is clearly work in progress_

## What is this?
**phop** is a handy CLI tool for managing terminal multiplexer sessions with ease. It lets you define reusable project environments using a simple configuration format and automates the creation and switching of sessions. **phop** sets up your terminal with the right windows, scripts, and tools instantly.

### Features:
- Fast navigation to any project from anywhere
- Define project templates for terminal multiplexer (tmux) environments
- Attach scripts and hooks to specific windows or whole sessions
- Friendly and flexible configuration

## Why?
I frequently switch between multiple projects through the day, often needing to change gcloud configurations, spin up SQL proxies, or run specific workflows. Repeating these setup steps manually was tedious and time consuming.

**phop** solves this by letting me define reusable templates that automatically set up my environment exactly how I need it. It also integrates seamlessly with my collection of small utility scripts for daily tasks and testing.

With a fuzzy searchable (fzf) list of projects, I can jump into any project in seconds, with everything already configured and ready to go.

## Dependencies
- [fzf](https://github.com/junegunn/fzf)
- [tmux](https://github.com/tmux/tmux)

## System wide installation from source (requires Go installed)
```bash
git clone https://github.com/wezik/tmux-sessionizer.git
cd tmux-sessionizer
go build -o ./build/phop
sudo mv ./build/phop /usr/local/bin/
```

## Usage
```
Usage: phop [command]

Available commands:
  a, add, c, create     Create a new project in the current working directory
  d, delete, r, remove  Delete a project
  e, edit               Edit a project with the given editor, defaults to system default
  s, script             Manage scripts (see phop script help for more info)
  l, list               List projects
  h, help               Show this help
```

## Current state
This project is in a very early experimental stage, I am figuring out how I want this thing to work and what it should be capable of.
