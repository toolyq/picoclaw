---
name: fetch-x-twitter
description: "Fetch or Get posts from X.com (Twitter) with nodejs"
instructions: Execute the script wiht nodejs ```node "D:\s\picoclaw\.picoclaw\workspace\skills\fetch-x-twitter\scripts\fetch-x-twitter.js" [n]```
---

# fetch-x-twitter Skill

Fetch the first `n` posts from X.com (Twitter) homepage with author, UTC timestamp, and direct URL.

## Steps for Agent

1. **Execute Script**: Run the Node.js script:
   ```powershell
   node "D:\s\picoclaw\.picoclaw\workspace\skills\fetch-x-twitter\scripts\fetch-x-twitter.js" [n]
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
- **Login**: Requires valid login session in `D:\s\picoclaw\.picoclaw\workspace\ai_chrome_data`

## Example
"获取 X 首页最近 5 条动态" → Execute → Extract data → Call message tool → Return empty → Done

