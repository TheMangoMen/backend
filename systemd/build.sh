#!/bin/bash

git -C /home/debian/backend reset HEAD --hard
git -C /home/debian/backend pull --rebase

go -C /home/debian/backend build cmd/app/main.go
