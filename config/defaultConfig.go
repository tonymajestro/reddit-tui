package config

const defaultConfiguration = `
#
# Default configuration for reddittui.
# Uncomment to configure
#

#[core]
#bypassCache = false
#logLevel = "Warn"

#[filter]
#keywords = ["drama"]
#subreddits = ["news", "politics"]

#[client]
#timeoutSeconds = 10
#cacheTtlSeconds = 3600

#[server]
#domain = "old.reddit.com"
#type = "old"
`
