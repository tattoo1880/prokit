#!/bin/zsh

git add .
git commit -m "update"
git push origin main


git tag -a v0.2.0 -m "release v0.2.0"
git push origin v0.2.0