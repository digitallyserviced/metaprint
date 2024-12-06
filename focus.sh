#!/bin/zsh -l

emulate -L zsh


export DISPLAY=:0
# bspc subscribe node | stdbuf -o0 -i0 grep node_focus | while read line; do typeset -a argss=(${=line});bspc query -T -n ${argss[4]} | fx .client.className ; done
# bspc subscribe node | stdbuf -o0 -i0 grep node_focus | while read line; do typeset -a argss=(${=line});bspc query -T -n ${argss[4]} | jq .client.className ; done&
bspc subscribe node | stdbuf -o0 -i0 grep node_focus | while read line; do typeset -a argss=(${=line}); xdotool getactivewindow getwindowname ; done &

