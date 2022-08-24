#!/usr/bin/env just --justfile

run cmd="":
  go run . {{cmd}}

table:
  cd  .example/table/ && go run .

un:
  cd  .example/unicode/ && go run .

list:
  @just run list

cfg:
  cp .code-github-workspace.yaml ~/

update:
  go get -u
  go mod tidy -v