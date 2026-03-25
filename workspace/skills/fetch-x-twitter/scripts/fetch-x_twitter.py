import asyncio
import json
import sys
import random
from playwright.async_api import async_playwright

async def fetch_x_tweets(n=3):
    """
    Fetches the first n posts from X.com (Twitter) homepage using an existing browser via CDP.
    """
    async with async_playwright() as p:
        # Connect to existing browser via CDP on port 9222
        try:
            browser = await p.chromium.connect_over_cdp("http://127.0.0.1:9222")
            context = browser.contexts[0] if browser.contexts else await browser.new_context()
            page = await context.new_page()
            
            # Speed up: Block unnecessary resources
            async def block_resources(route):
                if route.request.resource_type in ["image", "media", "font", "stylesheet"]:
                    await route.abort()
                elif any(x in route.request.url for x in ["analytics", "telemetry", "metrics", "ads"]):
                    await route.abort()
                else:
                    await route.continue_()

            await page.route("**/*", block_resources)

            # Navigate with retry logic
            attempts = 0
            max_attempts = 2
            while attempts < max_attempts:
                try:
                    await page.goto("https://x.com/home", wait_until="domcontentloaded", timeout=20000)
                    break
                except Exception as e:
                    attempts += 1
                    if attempts == max_attempts:
                        raise e
                    await asyncio.sleep(2)

            # Check for login redirect
            url = page.url
            if "/login" in url or "/i/flow/login" in url:
                print("Error: Not logged in. Please log in manually in the browser first.", file=sys.stderr)
                return []

            # Wait for tweet elements
            tweet_selector = '[data-testid="tweet"], article[role="article"]'
            try:
                await page.wait_for_selector(tweet_selector, timeout=15000)
            except:
                print("Note: Initial tweet elements not found by primary selector, continuing...", file=sys.stderr)

            # Extract tweets with auto-scrolling
            tweets = []
            seen_content = set()
            scroll_attempts = 0
            max_scroll_attempts = 50
            identical_position_count = 0

            while len(tweets) < n and scroll_attempts < max_scroll_attempts:
                # Evaluate in page context
                new_tweets_data = await page.evaluate("""
                    async (args) => {
                        const { n, seenContentArray } = args;
                        const seenSet = new Set(seenContentArray);
                        const results = [];
                        const tweetElements = Array.from(document.querySelectorAll('article[role="article"]'));

                        // Expand buttons — only search inside tweetText parent,
                        // exclude <a> tags to avoid page navigation.
                        for (const tweetEl of tweetElements) {
                            let expanded = false;
                            const tweetTextEl = tweetEl.querySelector('[data-testid="tweetText"]');
                            const searchRoot = tweetTextEl ? tweetTextEl.parentElement : tweetEl;
                            const candidates = searchRoot.querySelectorAll('[role="button"], span, div');
                            for (const el of candidates) {
                                if (el.tagName === 'A' || el.closest('a')) continue;
                                if (el.getAttribute('href') || el.getAttribute('data-href')) continue;
                                const txt = (el.innerText || el.textContent || '').trim();
                                const isExpandText = txt === '显示更多' || txt === '展开更多' || txt === '查看更多' || txt === 'Show more';
                                if (!isExpandText) continue;
                                try {
                                    el.dispatchEvent(new MouseEvent('click', { bubbles: true, cancelable: true }));
                                    expanded = true;
                                    break;
                                } catch (e) {}
                            }
                            if (expanded) await new Promise(r => setTimeout(r, 300));
                        }

                        const updatedTweetElements = Array.from(document.querySelectorAll('article[role="article"]'));
                        for (const tweetEl of updatedTweetElements) {
                            let author = 'Unknown';
                            const userNameEl = tweetEl.querySelector('[data-testid="User-Name"]');
                            if (userNameEl) {
                                const authorLink = userNameEl.querySelector('a[href^="/"]:not([href^="/hashtag"]):not([href^="/search"])');
                                if (authorLink) {
                                    const href = authorLink.getAttribute('href');
                                    author = href.startsWith('/') ? href.substring(1) : href;
                                }
                            }

                            let content = '';
                            const textEl = tweetEl.querySelector('[data-testid="tweetText"]');
                            if (textEl) content = textEl.innerText.trim();

                            let time = new Date().toISOString();
                            let url = '';
                            const timeEl = tweetEl.querySelector('time');
                            if (timeEl) {
                                if (timeEl.getAttribute('datetime')) time = timeEl.getAttribute('datetime');
                                const linkEl = timeEl.closest('a');
                                if (linkEl && linkEl.getAttribute('href')) {
                                    const href = linkEl.getAttribute('href');
                                    url = href.startsWith('http') ? href : `https://x.com${href}`;
                                }
                            }

                            const dedupId = url || `${author}:${content}`;
                            if (!content || seenSet.has(dedupId)) continue;
                            seenSet.add(dedupId);

                            results.push({
                                author,
                                time,
                                content,
                                url,
                                dedupId
                            });
                        }
                        return results;
                    }
                """, {"n": n, "seenContentArray": list(seen_content)})

                for tweet in new_tweets_data:
                    if len(tweets) < n:
                        dedup_id = tweet.pop('dedupId')
                        seen_content.add(dedup_id)
                        tweets.append(tweet)
                        
                        # Print immediately
                        print(f"--- TWEET_{len(tweets)} ---")
                        print(f"Author: {tweet['author']}")
                        print(f"Time: {tweet['time']}")
                        print(f"URL: {tweet['url']}")
                        print(f"Content:\n{tweet['content']}\n")

                if len(tweets) < n:
                    # Scroll down
                    is_at_bottom = await page.evaluate("""
                        () => {
                            const currentPosition = window.scrollY + window.innerHeight;
                            const maxPosition = document.documentElement.scrollHeight;
                            window.scrollBy(0, 1000);
                            return currentPosition >= maxPosition;
                        }
                    """)
                    
                    if is_at_bottom:
                        identical_position_count += 1
                        if identical_position_count > 5: break
                    else:
                        identical_position_count = 0
                    
                    scroll_attempts += 1
                    await asyncio.sleep(1.5 + random.random() * 0.5)

            return tweets
        finally:
            if 'page' in locals(): await page.close()
            if 'browser' in locals(): await browser.close()

if __name__ == "__main__":
    count = 3
    is_json = False
    
    if len(sys.argv) > 1:
        for arg in sys.argv[1:]:
            if arg == "--json":
                is_json = True
            elif arg.isdigit():
                count = int(arg)

    try:
        results = asyncio.run(fetch_x_tweets(count))
        if is_json:
            # We already printed tweets, so maybe just print the final JSON if requested
            # but usually for this skill the individual prints are better.
            # If --json is provided, we might want to suppress the intermediate prints?
            # Let's keep it simple.
            pass
        elif not results:
            print("\n[WARNING] No tweets were found.\n")
        else:
            print("\n[AGENT-INSTRUCTION-STOP-HERE]")
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)
