#!/usr/bin/env bash

# This is only meant to be run in dev, don't run in prod!
sudo tbot-server service stop && echo "service stopped"
sudo tbot-server service disable && echo "service disabled"

svcPath="/etc/systemd/system/tbot.service"
if [[ -f "$svcPath" ]]; then
  sudo rm "$svcPath" && echo "service file removed"
else
  echo "no service file to remove"
fi

if [[ -f "$HOME/bootstrap.key" ]]; then
  rm "$HOME/bootstrap.key" && echo "bootstrap key file removed"
else
  echo "no bootstrap key file to remove"
fi

goose down
sleep 2
goose up