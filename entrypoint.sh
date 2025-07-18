#!/bin/sh

set -e

export TOGETHERAI_API_KEY = "${INPUT_TOGETHERAI_API_KEY}"
export  GITHUB_TOKEN = "${INPUT_GITHUB_TOKEN}"

/app/ai-reviewer