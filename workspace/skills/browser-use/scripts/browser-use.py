"""
Browser automation using browser-use library with CDP connection.
Connects to an existing browser on localhost:9222 and executes AI-driven tasks.
"""

import asyncio
import json
import sys
import os
import argparse
import urllib.request
from datetime import datetime

# Ensure UTF-8 output
if hasattr(sys.stdout, 'reconfigure'):
    sys.stdout.reconfigure(encoding='utf-8')

try:
    from browser_use import Browser, ChatOpenAI
except ImportError as e:
    print(f"Error: Required packages not installed. Run: pip install browser-use langchain-openai", file=sys.stderr)
    print(f"Missing: {e}", file=sys.stderr)
    sys.exit(1)


async def run_browser_task(
    task: str,
    channel_id: str = "",
    chat_id: str = "",
    model_name: str = "nvidia/qwen/qwen3-next-80b-a3b-instruct",
    headless: bool = False,
    max_steps: int = 50,
    save_screenshots: bool = True
):
    """
    Execute a browser automation task using browser-use with CDP connection.
    
    Args:
        task: Natural language description of the task to perform
        channel_id: Optional channel for real-time updates
        chat_id: Optional chat ID for real-time updates
        model_name: LLM model to use (default: gpt-4o)
        headless: Whether to run in headless mode (ignored when using CDP)
        max_steps: Maximum number of steps to execute
        save_screenshots: Whether to save screenshots during execution
    """
    
    # Initialize LLM using local model
    api_key = "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"
    if not api_key:
        print("Error: OPENAI_API_KEY environment variable is not set.", file=sys.stderr)
        sys.exit(1)
    llm = ChatOpenAI(
        model="mistralai/mistral-small-4-119b-2603", 
        temperature=0.7,
        # base_url="http://localhost:1234/v1",
        base_url="https://integrate.api.nvidia.com/v1",
        api_key="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",
        timeout=600
    )
    
    browser = None
    results = []
    
    try:
        # Send initial status
        await send_update(
            channel_id, chat_id,
            f"🚀 Starting browser task: {task[:100]}{'...' if len(task) > 100 else ''}"
        )
        
        # Create browser instance with CDP connection
        browser = Browser(headless=headless, cdp_url="http://127.0.0.1:9222")
        
        # Create agent
        from browser_use import Agent
        
        agent = Agent(
            task=task,
            llm=llm,
            browser=browser,
            max_actions_per_step=5,
            llm_timeout=600,
            step_timeout=600,
            use_vision=save_screenshots,
        )
        
        # Track step progress
        step_count = 0
        
        async def on_step_end(callback_agent):
            nonlocal step_count
            step_count += 1
            last_item = callback_agent.history.history[-1] if callback_agent.history.history else None
            if last_item:
                action_preview = ""
                if last_item.model_output and last_item.model_output.action:
                    action_obj = last_item.model_output.action[0]
                    action_preview = str(action_obj.model_dump() if hasattr(action_obj, 'model_dump') else action_obj)[:50]
                
                success = True
                result_text = ""
                if last_item.result and len(last_item.result) > 0:
                    success = not last_item.result[-1].error
                    result_text = last_item.result[-1].error or last_item.result[-1].extracted_content or ""
                
                step_info = {
                    "step": step_count,
                    "action": action_preview,
                    "result": result_text,
                    "success": success,
                }
                results.append(step_info)
                
                status = "✅" if success else "❌"
                await send_update(
                    channel_id, chat_id,
                    f"{status} Step {step_count}: {action_preview}"
                )
                
                print(f"\n--- Step {step_count} ---")
                print(f"Action: {action_preview}")
                print(f"Result: {result_text}")

        # Run the agent with progress tracking
        history = await agent.run(max_steps=max_steps, on_step_end=on_step_end)
        
        # Get final result
        final_result = history.final_result() if history else "Task completed"
        
        await send_update(
            channel_id, chat_id,
            f"✨ Task completed!\n\nResult: {final_result}"
        )
        
        print(f"\n{'='*50}")
        print(f"Task completed in {step_count} steps")
        print(f"Final result: {final_result}")
        print(f"{'='*50}\n")
        
        return {
            "success": True,
            "steps": step_count,
            "final_result": final_result,
            "details": results
        }
        
    except Exception as e:
        error_msg = f"Error executing task: {str(e)}"
        print(f"\n❌ {error_msg}", file=sys.stderr)
        
        await send_update(channel_id, chat_id, f"❌ Error: {str(e)[:200]}")
        
        return {
            "success": False,
            "error": str(e),
            "steps": len(results),
            "details": results
        }
        
    finally:
        # Clean up
        if browser:
            try:
                await browser.close()
            except:
                pass


async def send_update(channel_id: str, chat_id: str, content: str):
    """Send update to user via webhook if configured."""
    webhook_url = os.environ.get("PICOCLAW_WEBHOOK_URL")
    
    if webhook_url and channel_id and chat_id:
        try:
            payload = json.dumps({
                "channel": channel_id,
                "chat_id": chat_id,
                "content": content
            }).encode()
            req = urllib.request.Request(
                webhook_url,
                data=payload,
                headers={'Content-Type': 'application/json'}
            )
            await asyncio.to_thread(urllib.request.urlopen, req, timeout=5)
        except Exception as e:
            print(f"Webhook push failed: {e}", file=sys.stderr)
    else:
        # Print to stdout if no webhook configured
        print(content)


def main():
    parser = argparse.ArgumentParser(
        description="Browser automation using browser-use with CDP connection"
    )
    parser.add_argument(
        "--task", "-t",
        required=True,
        help="Natural language task description"
    )
    parser.add_argument(
        "--channel", "-c",
        default="",
        help="Channel ID for real-time updates"
    )
    parser.add_argument(
        "--chat_id", "-ch",
        default="",
        help="Chat ID for real-time updates"
    )
    parser.add_argument(
        "--model", "-m",
        default="nvidia/mistralai/mistral-small-4-119b-2603",
        help="LLM model to use (default: nvidia/mistralai/mistral-small-4-119b-2603)"
    )
    parser.add_argument(
        "--max-steps", "-s",
        type=int,
        default=50,
        help="Maximum number of steps (default: 50)"
    )
    parser.add_argument(
        "--json", "-j",
        action="store_true",
        help="Output result as JSON"
    )
    parser.add_argument(
        "--no-vision",
        action="store_true",
        help="Disable vision/screenshots to avoid CDP timeouts"
    )
    
    args = parser.parse_args()
    
    # Run the task
    result = asyncio.run(run_browser_task(
        task=args.task,
        channel_id=args.channel,
        chat_id=args.chat_id,
        model_name=args.model,
        max_steps=args.max_steps,
        save_screenshots=not args.no_vision
    ))
    
    if args.json:
        print(json.dumps(result, ensure_ascii=False, indent=2))
    
    sys.exit(0 if result.get("success", False) else 1)


if __name__ == "__main__":
    main()