#!/bin/bash

# Set script to exit immediately on error
set -e

cdklocal bootstrap aws://000000000000/us-east-1 && cdklocal synth && cdklocal deploy
