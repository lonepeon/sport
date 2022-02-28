#! /usr/bin/env bash

success() {
  echo "$(tput setaf 2)[OK] $(tput setaf 0) $*"
}

fail() {
  echo "$(tput setaf 1)[ERR]$(tput setaf 0) $*"
}
