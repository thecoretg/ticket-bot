#!/usr/bin/env bash

# This is only meant to be run in dev, don't run in prod!
sudo tbot-server service disable
goose down
sleep 2
goose up