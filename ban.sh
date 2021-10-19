#!/usr/bin/env bash
ufw insert 1 deny from $1 comment "$2"
ss -K dst $1