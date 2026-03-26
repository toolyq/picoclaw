---
name: fetch-x-twitter
description: "Fetch or Get posts from X.com (Twitter)"
instructions: |
  Execute the script with Python:
  ```powershell
  python "D:\s\picoclaw\.picoclaw\workspace\skills\fetch-x-twitter\scripts\fetch-x-twitter.py" [n]
  ```
---

# fetch-x-twitter Skill

Fetch the first `n` posts from X.com (Twitter) homepage with author, UTC timestamp, and direct URL.

## Steps for Agent

1. **Send to User**: Call the `message` tool with this extracted content.
   - Use the exact data from between the markers
   - Do NOT modify or summarize
   - Pass it directly to the message tool

2. **End Task**: 
   - After calling message tool, return empty string
   - Do NOT continue iteration
   - Do NOT call LLM again
   - Task is COMPLETE

## Example
"获取 X 首页最近 5 条动态" → Execute → Extract data → Call message tool → Return empty → Done

