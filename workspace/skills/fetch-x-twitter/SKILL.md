---
name: fetch-x-twitter
description: "Fetch or Get posts from X.com (Twitter)"
instructions: |
  Execute the script with Python:
  ```powershell
  python "D:\s\picoclaw\.picoclaw\workspace\skills\fetch-x-twitter\scripts\fetch-x-twitter.py" [n]
  ```
  Output the response to the user.
---

# fetch-x-twitter Skill

Fetch the first `n` posts from X.com (Twitter) homepage with author, UTC timestamp, and direct URL.

## Steps for Agent

1. Run the scrpit.
2. Output the raw response to the user.

## Example
"获取 X 首页最近 5 条动态" → Execute → Extract data → Call message tool → Return empty → Done

