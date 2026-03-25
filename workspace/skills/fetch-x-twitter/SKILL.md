---
name: fetch-x-twitter
description: "Fetch or Get posts from X.com (Twitter)"
instructions: |
  Execute the script with Python:
  ```powershell
  python "D:\s\picoclaw\.picoclaw\workspace\skills\fetch-x-twitter\scripts\fetch-x-twitter.py" [n]
  ```
  Or with Node.js:
  ```powershell
  node "D:\s\picoclaw\.picoclaw\workspace\skills\fetch-x-twitter\scripts\fetch-x-twitter.js" [n]
  ```
---

# fetch-x-twitter Skill

Fetch the first `n` posts from X.com (Twitter) homepage with author, UTC timestamp, and direct URL.

## Steps for Agent

1. **Execute Script**: Run the Node.js or Python script:
   ```powershell
   python "D:\s\picoclaw\.picoclaw\workspace\skills\fetch-x-twitter\scripts\fetch-x-twitter.js" [n]
   # Or
   node "D:\s\picoclaw\.picoclaw\workspace\skills\fetch-x-twitter\scripts\fetch-x-twitter.py" [n]
   ```
   *(Default `n` is 3)*

2. **Extract Data**: Find the content between `=== DATA_START ===` and `=== DATA_END ===` markers.

3. **Send to User**: Call the `message` tool with this extracted content.
   - Use the exact data from between the markers
   - Do NOT modify or summarize
   - Pass it directly to the message tool

4. **End Task**: 
   - After calling message tool, return empty string
   - Do NOT continue iteration
   - Do NOT call LLM again
   - Task is COMPLETE

## Important Notes
- **Browser**: Requires a running browser with remote debugging enabled on port 9222 (e.g., `chrome.exe --remote-debugging-port=9222`).
- **Login**: Ensure you are already logged in to X (Twitter) in the browser instance before running.

## Example
"获取 X 首页最近 5 条动态" → Execute → Extract data → Call message tool → Return empty → Done

