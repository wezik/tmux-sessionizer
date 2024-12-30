#!/bin/bash

# passing current directory to app
APP_DIR=$(dirname "$(realpath "$0")")
ORIGINAL_DIR="$(pwd)"
cd $APP_DIR || exit

# "magic" capture, using script to grab stderr without interfering with the app (fzf is acting weird with a simple redirect)
stderr=$(script -q -c "./zig-out/bin/tmux-sessionizer \"$ORIGINAL_DIR\" \"$1\"" /dev/null 2>&1 | tee /dev/stderr)
# stderr=$(script -q -c "zig build run -- \"$ORIGINAL_DIR\" \"$1\"" /dev/null 2>&1 | tee /dev/stderr)

# take the last line and operate based on cmd received, also clears it from the screen
last_line=$(echo "$stderr" | tail -n 1 | tr -d '\r' | tr -d '\n')

IFS=":" read -r var1 var2 var3 <<< "$last_line"
if [[ $var1 != *"signal"* ]]; then
        exit 0
fi

echo -e -n "\033[1A\033[K"

if [[ $var2 == "edit" ]]; then
        editor=$2
        if [[ -z $editor ]]; then
                editor="nano"
        fi
        $editor $var3
elif [[ $var2 == "tmux_attach" ]]; then
        tmux attach -t "$var3"
else
        echo "unhandled $var1:$var2:$var3 (dev is at fault)"
        exit 1
fi

