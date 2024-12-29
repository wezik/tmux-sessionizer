#!/bin/bash

# passing current directory to app
APP_DIR=$(dirname "$(realpath "$0")")
ORIGINAL_DIR="$(pwd)"
cd $APP_DIR || exit

# "magic" capture, using script to grab stderr without interfering with the app (fzf is acting weird with a simple redirect)
stderr=$(script -q -c "./zig-out/bin/tmux-sessionizer \"$ORIGINAL_DIR\" \"$@\"" /dev/null 2>&1 | tee /dev/stderr)
# stderr=$(script -q -c "zig build run -- \"$ORIGINAL_DIR\" \"$@\"" /dev/null 2>&1 | tee /dev/stderr)

# take the last line and operate based on cmd received, also clears it from the screen
last_line=$(echo "$stderr" | tail -n 1 | tr -d '\r' | tr -d '\n')
if [[ $last_line == *"message:"* ]]; then
        message=$(echo "$last_line" | awk -F':' '{print $2}')
        echo -e -n "\033[1A\033[K"
        echo "$message"
elif [[ $last_line == *"error:"* ]]; then
        error=$(echo "$last_line" | awk -F':' '{print $2}')
        echo -e -n "\033[1A\033[K"
        echo "$error"
        exit 1
elif [[ $last_line == *"attach_hook:"* ]]; then
        session_name=$(echo "$last_line" | awk -F':' '{print $2}')
        echo -e -n "\033[1A\033[K"
        tmux attach -t "$session_name"
fi
