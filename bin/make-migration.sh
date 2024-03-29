#!/usr/bin/env bash
time=$(date +%s)
touch $(pwd)/migrations/"$time"-"$1".sql