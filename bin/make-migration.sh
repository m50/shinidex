#!/usr/bin/env bash
time=$(date +%s)
touch $(pwd)/pkg/database/migrations/"$time"-"$1".sql
