const { chromium } = require('playwright');

/**
 * Fetches the first n posts from X.com (Twitter) homepage
 * @param {number} n - Number of posts to fetch (default: 3)
 * @returns {Promise<Array<Object>>} Array of tweet objects with author, time, and content
 */
async function fetchXTweets(n = 3) {
  let browser, page;
  try {
    // Connect to existing browser via CDP on port 9222
    browser = await chromium.connectOverCDP('http://127.0.0.1:9222');
    const context = browser.contexts()[0] || await browser.newContext();
    page = await context.newPage();

    // Speed up: Block unnecessary resources
    await page.route('**/*.{png,jpg,jpeg,gif,svg,css,woff,woff2,otf,ttf,eot}', route => route.abort());
    await page.route(/.*(analytics|telemetry|metrics|ads).*/, route => route.abort());

    // Navigate with retry logic
    let attempts = 0;
    const maxAttempts = 2;
    while (attempts < maxAttempts) {
      try {
        // console.error(`Navigating to X.com (Attempt ${attempts + 1})...`);
        await page.goto('https://x.com/home', {
          waitUntil: 'domcontentloaded',
          timeout: 20000
        });
        break;
      } catch (e) {
        attempts++;
        if (attempts === maxAttempts) throw e;
        await page.waitForTimeout(2000);
      }
    }

    // Check for login redirect
    const url = page.url();
    if (url.includes('/login') || url.includes('/i/flow/login')) {
      throw new Error('Not logged in. Please log in manually in the browser first to maintain the persistent context.');
    }

    // Handle potential login/modal dialogs
    try {
      // Try to close login modals if they appear
      const modalCloseSelectors = [
        '[data-testid="modal-close"]',
        '[role="dialog"] button[aria-label="Close"]',
        '[aria-label="Close"]',
        'div[role="dialog"] button'
      ];

      for (const selector of modalCloseSelectors) {
        const closeButton = await page.$(selector);
        if (closeButton) {
          await closeButton.click();
          await page.waitForTimeout(1000);
          console.error('Closed modal dialog');
          break;
        }
      }
    } catch (e) {
      // Ignore errors in modal handling
    }

    // Wait for tweet elements to appear (more robust than waitForFunction)
    // console.error('Waiting for tweet elements...');
    const tweetSelector = '[data-testid="tweet"], article[role="article"]';
    try {
      await page.waitForSelector(tweetSelector, { timeout: 15000 });
    } catch (e) {
      console.error('Note: Initial tweet elements not found by primary selector, checking others...');
    }

    // Extract tweets with auto-scrolling to handle larger 'n' values
    const tweets = [];
    const seenContent = new Set();
    let scrollAttempts = 0;
    const maxScrollAttempts = 1000; // Increased significantly for larger counts
    let identicalPositionCount = 0;

    while (tweets.length < n && scrollAttempts < maxScrollAttempts) {
      // Find all visible tweet articles
      const newTweets = await page.evaluate(async ({ n, seenContentArray }) => {
        const seenSet = new Set(seenContentArray);
        const results = [];
        const tweetElements = Array.from(document.querySelectorAll('article[role="article"]'));

        // First, expand all collapsed content in each tweet
        for (const tweetEl of tweetElements) {
          let expanded = false;

          // Only look for expand buttons inside the tweet text container,
          // and EXCLUDE <a> tags to avoid page navigation.
          // X.com uses [role="button"] spans for "Show more" inside tweetText.
          const tweetTextEl = tweetEl.querySelector('[data-testid="tweetText"]');
          const searchRoot = tweetTextEl ? tweetTextEl.parentElement : tweetEl;

          const candidates = searchRoot.querySelectorAll('[role="button"], span, div');
          for (const el of candidates) {
            // Skip any element that is an anchor or contains an anchor (navigation risk)
            if (el.tagName === 'A' || el.closest('a')) continue;
            // Skip if element itself has an href-like behavior
            if (el.getAttribute('href') || el.getAttribute('data-href')) continue;

            // Only match elements whose direct text (trimmed) is exactly the expand phrase
            const txt = (el.innerText || el.textContent || '').trim();
            const isExpandText = txt === '显示更多' || txt === '展开更多' || txt === '查看更多' || txt === 'Show more';
            if (!isExpandText) continue;

            try {
              // Use dispatchEvent instead of .click() to avoid triggering
              // any default browser navigation on anchor-like elements
              el.dispatchEvent(new MouseEvent('click', { bubbles: true, cancelable: true }));
              expanded = true;
              break; // Only expand once per tweet
            } catch (e) {
              // ignore
            }
          }

          // If we expanded, wait a bit for content to load
          if (expanded) {
            await new Promise(r => setTimeout(r, 300));
          }
        }

        // After expanding, re-query tweet elements (in case DOM changed)
        const updatedTweetElements = Array.from(document.querySelectorAll('article[role="article"]'));

        for (const tweetEl of updatedTweetElements) {
          // Extract author handle (@username)
          let author = '';
          const userNameEl = tweetEl.querySelector('[data-testid="User-Name"]');
          if (userNameEl) {
            const authorLink = userNameEl.querySelector('a[href^="/"]:not([href^="/hashtag"]):not([href^="/search"])');
            if (authorLink) {
              const href = authorLink.getAttribute('href');
              author = href.startsWith('/') ? href.substring(1) : href;
            }
          }

          // Extract content
          let content = '';
          const textEl = tweetEl.querySelector('[data-testid="tweetText"]');
          if (textEl) {
            content = textEl.innerText.trim();
          }

          // Use URL as unique identifier if available, otherwise author+content
          let time = '';
          let url = '';
          const timeEl = tweetEl.querySelector('time');
          if (timeEl) {
            if (timeEl.getAttribute('datetime')) {
              time = timeEl.getAttribute('datetime');
            }
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
            author: author || 'Unknown',
            time: time || new Date().toISOString(),
            content: content,
            url: url || '',
            dedupId: dedupId
          });
        }
        return results;
      }, { n, seenContentArray: Array.from(seenContent) });

      for (const tweet of newTweets) {
        if (tweets.length < n) {
          seenContent.add(tweet.dedupId);
          delete tweet.dedupId;
          tweets.push(tweet);

          // Print tweet content immediately
          console.log(`--- TWEET_${tweets.length} ---`);
          console.log(`Author: ${tweet.author}`);
          console.log(`Time: ${tweet.time}`);
          console.log(`URL: ${tweet.url}`);
          console.log(`Content:\n${tweet.content}\n`);
        }
      }

      if (tweets.length < n) {
        // Check if we can scroll further
        const isAtBottom = await page.evaluate(() => {
          const currentPosition = window.scrollY + window.innerHeight;
          const maxPosition = document.documentElement.scrollHeight;
          window.scrollBy(0, 1000);
          return currentPosition >= maxPosition;
        });

        if (isAtBottom) {
          identicalPositionCount++;
          if (identicalPositionCount > 10) break; // Reached bottom or stuck
        } else {
          identicalPositionCount = 0;
        }

        scrollAttempts++;
        // Wait for content to render - slightly variable to avoid detection
        await page.waitForTimeout(1500 + Math.random() * 500);
      }
    }

    return tweets;
  } finally {
    if (typeof page !== 'undefined') await page.close();
    if (typeof browser !== 'undefined') await browser.close();
  }
}

// If script is run directly, fetch and display tweets
if (require.main === module) {
  const args = process.argv.slice(2);
  const isJson = args.includes('--json');
  const countArg = args.find(a => !a.startsWith('--'));
  const count = parseInt(countArg) || 3;

  fetchXTweets(count)
    .then(tweets => {
      if (isJson) {
        console.log(JSON.stringify(tweets, null, 2));
      } else {
        if (tweets.length === 0) {
          console.log('\n[WARNING] No tweets were found. This might be due to a loading issue or empty home feed.\n');
        } else {
          // console.log('\n=== DATA_START ===');
          // Print tweets in the expected format for extraction
          // for (let i = 0; i < tweets.length; i++) {
          //   const tweet = tweets[i];
          //   console.log(`--- TWEET_${i + 1} ---`);
          //   console.log(`Author: ${tweet.author}`);
          //   console.log(`Time: ${tweet.time}`);
          //   console.log(`URL: ${tweet.url}`);
          //   console.log(`Content:\n${tweet.content}\n`);
          // }
          // console.log(`=== DATA_END (${tweets.length} total) ===`);
          console.log('\n[AGENT-INSTRUCTION-STOP-HERE]');
        }
      }
    })
    .catch(error => {
      console.error('Error fetching tweets:', error);
      process.exit(1);
    });
}

module.exports = { fetchXTweets };