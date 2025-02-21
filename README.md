# Reddittui
A lightweight terminal application for browsing Reddit from your command line. Powered by [bubbletea](https://github.com/charmbracelet/bubbletea)

## Features
- **Subreddit Browsing:** Navigate through your favorite subreddits.
- **Post Viewing:** Read text posts and comments.
- **Keyboard Navigation:** Scroll and select posts using vim/standard keyboard shortcuts.
- **Configurable**: Customize caching behavior and define subreddit filters using a configuration file

## Demo
https://github.com/user-attachments/assets/40d61ef3-3a95-4a26-8c49-bec616f6ae1c

## Prerequisites

- **Go:** Version 1.16 or newer
- **Terminal:** A Unix-like terminal (Linux, macOS, or similar).
- **POSIX Utilities:** The `install` command is used for installation, which is available on both Linux and macOS.

## Installation
Clone the repository and run the install script: 

```bash
git clone https://github.com/tonymajestro/reddit-tui.git reddittui
cd reddittui
./install.sh
```

## Usage
Run the installed binary from your preferred terminal:

```bash
reddittui
```

## Keybindings
- Navigation
  - **h, j, k, l:** vim movement
  - **left, right, up, down:** normal movement
- Posts page
  - **s**: Switch subreddits
- Comments page
  - **o**: Open post link in browser
  - **c**: Collapse comments
- Misc
  - **H:** Go to home page
  - **backspace**: Go back
  - **q, esc**: Exit reddittui

## Configuration
After running the reddittui binary, a configuration file will be created at **~/.config/reddittui/reddittui.toml**. Example configuration file:

```toml
# Core configuration
[core]
bypassCache = false
logLevel = "Debug"
clientTimeout = 10

# Filter out posts containing keywords or belonging to certain subreddits
[filter]
subreddits = ["news", "politics", "videogames"]
keywords = ["trump", "nazi", "pizza"]
```

## Acknowledgments
Reddittui is based on the [bubbletea](https://github.com/charmbracelet/bubbletea) framework. It also takes inspiration from [circumflex](https://github.com/bensadeh/circumflex), a hackernews terminal browser.
