---
name: fetch-x-twitter
description: "Fetch or Get posts from X.com (Twitter)"
instructions: |
  Execute the script with Python. Because you need to forward the raw response stream directly to the user (bypassing the model), examine your system context to find the current Channel and Chat ID, and pass them as flags:
  ```powershell
  python "d:\git\picoclaw\workspace\skills\fetch-x-twitter\scripts\fetch-x-twitter.py" [n] --channel "<channel>" --chat_id "<chat_id>"
  ```
---

# fetch-x-twitter Skill

Fetch the first `n` posts from X.com (Twitter) and forward the raw response stream to the user as soon as possible.

## Steps for Agent

1. Determine the current `channel` and `chat_id` from your system context block (look for "Current Session").
2. Run the script passing those arguments. The script will automatically forward them in real-time to the user bypassing the LLM.
3. Finish the turn silently (or with a brief summary).

## Example
"获取 X 首页最近 5 条动态" → Execute → Forward the raw response  stream to the user → Done

