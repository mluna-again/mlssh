#! /usr/bin/env bash

pane_count=$(tmux list-panes | wc -l)
if [ -z "$TMUX" ] || (( pane_count > 1 )); then
  echo "please run this script inside an empty tmux window :)"
  exit 1
fi

if [ ! -d stress ]; then
  mkdir stress || exit
fi

echo "making keys..."
count=0
while (( count < 10 )); do
  keypath="stress/id_rsa${count}"
  if [ -f "$keypath" ]; then
    count=$(( count + 1 ))
    continue
  fi

  if ! ssh-keygen -t rsa -b 4096 -f "$keypath" -N '' &>/dev/null; then
    echo failed at creating ssh key
    exit 1
  fi

  count=$(( count + 1 ))
done
echo "done."

tmux split-window || exit
tmux split-window || exit
tmux split-window || exit
tmux select-layout tiled || exit
tmux split-window || exit
tmux split-window || exit
tmux split-window || exit
tmux select-layout tiled || exit
tmux split-window || exit
tmux split-window || exit
tmux select-layout tiled || exit

count=0
while read -r id; do
  tmux send-keys -t "$id" ssh Space -i Space stress/id_rsa${count} Space mikoshi.net.co Enter
  sleep 1
  count=$(( count + 1 ))
done < <(tmux list-panes -F '#{pane_id}')
