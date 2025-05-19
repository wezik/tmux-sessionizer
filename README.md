# Tmux/Project sessionizer

## What is this?

### CLI app that lets you:
- define terminal multiplexer templates per project
- automatically attach custom environments / execute shell scripts specific to a project
- define per project util scripts with aliases to quickly run them
- jump between projects quickly

## Why?

### My own needs
I often had to jump between projects multiple times a day, setting custom gcloud configurations, running sql proxies, etc.
I also like writing small util scripts that help me with my daily tasks / testing.
This app lets me define a template that auto sets-up my environment, with a nice fzf searchable list of projects.

## Dependencies
- [fzf](https://github.com/junegunn/fzf)
- [tmux](https://github.com/tmux/tmux)

## Installation
```bash
git clone https://github.com/wezik/tmux-sessionizer.git
cd tmux-sessionizer
./install.sh
```

## Usage
```
Usage: phop [command]

Available commands:
  a, add, c, create     Create a new project in the current working directory
  d, delete, r, remove  Delete a project
  e, edit               Edit a project with the given editor, defaults to nano
  s, script             Manage scripts (see phop script help for more info)
  l, list               List projects
  h, help               Show this help
