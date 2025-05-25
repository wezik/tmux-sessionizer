package cli

const helpMessage = `Phop - Project Hopper CLI

Usage:
  phop <command> [args...]

Commands:
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

Notes:
- Interactive selection is powered by fzf for commands without explicit [name].
- Aliases allow fast access for power users (e.g., 'phop c', 'phop s').`
