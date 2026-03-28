name: browser-use
description: "Browser automation using browser-use library with CDP connection. Supports complex web tasks like navigation, form filling, data extraction, and multi-step workflows."
instructions: |
  Execute the script with Python. Pass the task description and optionally channel/chat_id for real-time updates:
  ```powershell
  python "D:\s\picoclaw\.picoclaw\workspace\skills\browser-use\scripts\browser-use.py" --task "<task description>" [--channel "<channel>"] [--chat_id "<chat_id>"]
  ```
  
  The browser must have remote debugging enabled on port 9222.
  
  Examples:
  - Navigate to a website and extract information
  - Fill forms and submit them
  - Perform multi-step workflows on web pages
  - Scrape data from dynamic websites

---

# browser-use Skill

Automate browser tasks using the browser-use library with an existing browser connected via CDP.

## Prerequisites

1. Chrome/Edge browser must be running with remote debugging enabled:
   ```powershell
   # Chrome
   chrome.exe --remote-debugging-port=9222
   
   # Edge
   msedge.exe --remote-debugging-port=9222
   ```

2. Required Python packages:
   ```powershell
   pip install browser-use playwright langchain-openai
   ```

3. OpenAI API key set in environment:
   ```powershell
   $env:OPENAI_API_KEY = "your-api-key"
   ```

## Usage

### Basic Usage
```powershell
python "D:\s\picoclaw\.picoclaw\workspace\skills\browser-use\scripts\browser-use.py" --task "Go to example.com and find the main heading"
```

### With Real-time Updates
```powershell
python "D:\s\picoclaw\.picoclaw\workspace\skills\browser-use\scripts\browser-use.py" --task "Search for Python tutorials on Google" --channel "telegram" --chat_id "123456"
```

### Complex Tasks
```powershell
python "D:\s\picoclaw\.picoclaw\workspace\skills\browser-use\scripts\browser-use.py" --task "Go to GitHub, search for 'browser-use', open the first result, and summarize the README"
```

## Features

- **CDP Connection**: Connects to existing browser on port 9222
- **AI-Powered**: Uses LLM to understand and execute complex tasks
- **Real-time Updates**: Optionally sends progress updates to user
- **Error Handling**: Robust error handling and retry logic
- **Screenshots**: Can capture screenshots during execution

## Steps for Agent

1. Ensure the user's browser has remote debugging enabled on port 9222
2. Verify OPENAI_API_KEY is set in environment
3. Run the script with the task description
4. The script will execute the task and return results
5. Optionally forward real-time updates to the user

## Example Tasks

- "Navigate to news.ycombinator.com and list the top 5 headlines"
- "Go to Amazon and search for 'wireless headphones', then extract product names and prices"
- "Fill out the contact form on example.com with test data"
- "Login to a website and navigate to the dashboard"