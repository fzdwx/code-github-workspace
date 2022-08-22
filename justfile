#!/usr/bin/env just --justfile

run cmd="":
  go run . {{cmd}}

list:
  @just run list

cfg:
  cp .code-github-workspace.yaml ~/

update:
  go get -u
  go mod tidy -v