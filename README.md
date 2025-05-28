# Reddittui
A lightweight terminal application for browsing Reddit from your command line. Powered by [bubbletea](https://github.com/charmbracelet/bubbletea)

## Features
- **Subreddit Browsing:** Navigate through your favorite subreddits.
- **Post Viewing:** Read text posts and comments.
- **Keyboard Navigation:** Scroll and select posts using vim/standard keyboard shortcuts.
- **Configurable**: Customize caching behavior and define subreddit filters using a configuration file

## Demo
https://github.com/user-attachments/assets/40d61ef3-3a95-4a26-8c49-bec616f6ae1c

## Installation

### Git
#### Prerequisites
- **Git**
- **Go:** Version 1.16 or newer
- **Terminal:** A Unix-like terminal (Linux, macOS, or similar).
- **POSIX Utilities:** The `install` command is used for installation, which is available on both Linux and macOS.

Clone the repository and run the install script: 

```bash
git clone https://github.com/tonymajestro/reddit-tui.git reddittui
cd reddittui
./install.sh
```

To remove reddittui run the uninstall script:

```bash
./uninstall.sh
```

### Arch
Arch users can install reddittui from the AUR using yay or other AUR helpers.

[Pre-compiled](https://aur.archlinux.org/packages/reddit-tui-bin) and [source packages](https://aur.archlinux.org/packages/reddit-tui) are available.

```bash
yay -S reddit-tui-bin
```

```bash
yay -S reddit-tui
```

### Nix
Nix users can try it in a shell or add it to their system config like this.
```bash
nix-shell -p reddit-tui
```
```nix
  environment.systemPackages = [
      pkgs.reddit-tui
    ];
```

## Usage
Run the installed binary from your preferred terminal:

```bash
# Open reddittui, navigating to the home page
reddittui

# Open reddittui, navigating to a specific subreddit
reddittui --subreddit dogs

# Open reddittui, navigating to a specific post by its ID
reddittui --post 1iyuce4
```

## Keybindings
- Navigation
  - **h, j, k, l:** Vim movement
  - **left, right, up, down:** Normal movement
  - **g**: Go to top of page
  - **G**: Go to bottom of page
  - **s**: Switch subreddits
- Posts page
  - **L**: Load more posts
- Comments page
  - **o**: Open post link in browser
  - **c**: Collapse comments
- Misc
  - **H:** Go to home page
  - **backspace**: Go back
  - **q, esc**: Exit reddittui

## Configuration files
After running the reddittui binary, the following files will be initialized:
- Configuration file:
  - `~/.config/reddittui/reddittui.toml`
- Log file:
  - `~/.local/state/reddittui.log`
- Cache
  - `~/.cache/reddittui/`

Sample configuration:
```toml
# Core configuration
[core]
bypassCache = false
logLevel = "Warn"

# Filter out posts containing keywords or belonging to certain subreddits
[filter]
subreddits = ["news", "politics"]
keywords = ["pizza", "pineapple"]

# Configure client timeout and cache TTL. By default, subreddit posts and comments are cached for 1 hour.
[client]
timeoutSeconds = 10
cacheTtlSeconds = 3600

# Configure which reddit server to use. Default is old.reddit.com but redlib servers are also supported
[server]
domain = "old.reddit.com"
type = "old"
```

## Redlib
For enhanced privacy, private [Redlib backends](https://github.com/redlib-org/redlib) are supported. A list of Redlib servers can be found [here](https://github.com/redlib-org/redlib-instances/blob/main/instances.json). Use the following configuration to use a Redlib server instead of old.reddit.com:

```toml
[server]
domain = "safereddit.com"
type = "redlib"
```

## Acknowledgments
Reddittui is based on the [bubbletea](https://github.com/charmbracelet/bubbletea) framework. It also takes inspiration from [circumflex](https://github.com/bensadeh/circumflex), a hackernews terminal browser.
