#!/bin/bash

git -C /home/debian/backend reset HEAD --hard
git -C /home/debian/backend pull --rebase

/usr/local/go/bin/go -C /home/debian/backend mod tidy
/usr/local/go/bin/go -C /home/debian/backend build cmd/app/main.go
