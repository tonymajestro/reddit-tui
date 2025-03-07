package main

import (
	"bytes"
	"reddittui/components"
	"reddittui/config"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

const (
	testTimeout    = 20 * time.Second
	testDomain     = "safereddit.com"
	testServerType = "redlib"
)

var testConf = config.Config{
	Core: config.CoreConfig{
		BypassCache: true,
		LogLevel:    "Warn",
	},
	Client: config.ClientConfig{
		TimeoutSeconds: int(testTimeout.Seconds()),
	},
	Server: config.ServerConfig{
		Domain: testDomain,
		Type:   testServerType,
	},
}

func TestStartup(t *testing.T) {
	t.Logf("Testing startup...")

	tui := components.NewRedditTui(testConf, "", "")
	tm := teatest.NewTestModel(t, tui, teatest.WithInitialTermSize(300, 100))

	t.Logf("\tVerify the loading screen shows on startup...")
	WaitFor(t, tm, "loading reddit.com...")

	t.Logf("\tVerify the home page loads...")
	WaitFor(t, tm, "The front page of the internet")
}

func TestSwitchSubreddit(t *testing.T) {
	t.Logf("Testing switching subreddit...")
	configuration := config.NewConfig()
	configuration.Core.BypassCache = true

	tui := components.NewRedditTui(configuration, "", "")
	tm := teatest.NewTestModel(t, tui, teatest.WithInitialTermSize(300, 100))

	t.Logf("\tVerify home page loads...")
	WaitFor(t, tm, "The front page of the internet")

	t.Logf("\tVerify subreddit selection modal shows...")
	WaitForWithInputs(t, tm, "s", "Choose a subreddit:")

	t.Logf("\tVerify dogs subreddit loads...")
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("dogs"),
	})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	WaitFor(t, tm, "r/dogs")
}

func TestReturnToHomePage(t *testing.T) {
	t.Logf("Testing returning to the home page after switching subreddits...")
	configuration := config.NewConfig()
	configuration.Core.BypassCache = true

	tui := components.NewRedditTui(configuration, "", "")
	tm := teatest.NewTestModel(t, tui, teatest.WithInitialTermSize(300, 100))

	t.Logf("\tVerify home page loads...")
	WaitFor(t, tm, "The front page of the internet")

	t.Logf("\tVerify subreddit selection modal shows...")
	WaitForWithInputs(t, tm, "s", "Choose a subreddit:")

	t.Logf("\tVerify dogs subreddit loads...")
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("dogs"),
	})
	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})
	WaitFor(t, tm, "r/dogs ")

	t.Logf("\tVerify home page loads")
	WaitForWithInputs(t, tm, "H", "The front page of the internet")
}

func TestShowComments(t *testing.T) {
	t.Logf("Testing show post comments...")
	configuration := config.NewConfig()
	configuration.Core.BypassCache = true

	tui := components.NewRedditTui(configuration, "", "")
	tm := teatest.NewTestModel(t, tui, teatest.WithInitialTermSize(300, 100))

	t.Logf("\tVerify home page loads...")
	WaitFor(t, tm, "The front page of the internet")

	t.Logf("\tVerify subreddit selection modal shows...")
	WaitForWithInputs(t, tm, "s", "Choose a subreddit:")

	t.Logf("\tVerify dogs subreddit loads...")
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("dogs"),
	})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	WaitFor(t, tm, "r/dogs", "/r/dogs")

	t.Logf("\tVerify comments header loads...")
	time.Sleep(time.Second)
	WaitForWithInputs(t, tm, "l", "submitted", "ago by", "point", "comment")
}

func WaitFor(t *testing.T, tm *teatest.TestModel, messages ...string) {
	WaitForWithInputs(t, tm, "", messages...)
}

func WaitForWithInputs(t *testing.T, tm *teatest.TestModel, inputs string, messages ...string) {
	if len(inputs) > 0 {
		tm.Send(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune(inputs),
		})
	}

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		for _, message := range messages {
			if !bytes.Contains(bts, []byte(message)) {
				return false
			}
		}

		return true
	}, teatest.WithCheckInterval(time.Millisecond*50), teatest.WithDuration(testTimeout))
}
