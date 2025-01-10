package reddit

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var subredditSuggestions = []string{
	"15minutefood",
	"adviceanimals",
	"all",
	"animalsbeingbros",
	"animalsbeingderps",
	"animalsbeingjerks",
	"anime",
	"anime_irl",
	"apple",
	"art",
	"askreddit",
	"askscience",
	"aww",
	"awwducational",
	"backpacking",
	"baking",
	"battlestations",
	"beamazed",
	"bestof",
	"bikinibottomtwitter",
	"biology",
	"bitcoin",
	"boardgames",
	"bodyweightfitness",
	"books",
	"buildapc",
	"camping",
	"canada",
	"careerguidance",
	"cars",
	"cats",
	"cfb",
	"changemyview",
	"chatgpt",
	"chemistry",
	"comicbooks",
	"compsci",
	"contagiouslaughter",
	"cooking",
	"coolguides",
	"cozyplaces",
	"crappydesign",
	"creepy",
	"cryptocurrency",
	"dadjokes",
	"damnthatsinteresting",
	"dataisbeautiful",
	"dating",
	"dating_advice",
	"daytrading",
	"design",
	"destinythegame",
	"digitalpainting",
	"diwhy",
	"diy",
	"dnd",
	"documentaries",
	"drawing",
	"dundermifflin",
	"eatcheapandhealthy",
	"economics",
	"eldenring",
	"entertainment",
	"entrepreneur",
	"ethereum",
	"europe",
	"expectationvsreality",
	"explainlikeimfive",
	"eyebleach",
	"facepalm",
	"fantasy",
	"fantasyfootball",
	"fauxmoi",
	"femalefashionadvice",
	"fitness",
	"food",
	"foodhacks",
	"formula1",
	"fortnitebr",
	"frugal",
	"funny",
	"funnyanimals",
	"futurology",
	"gadgets",
	"gameofthrones",
	"games",
	"gaming",
	"gardening",
	"genshin_impact",
	"getmotivated",
	"gifrecipes",
	"gifs",
	"google",
	"hair",
	"hardware",
	"health",
	"healthyfood",
	"highqualitygifs",
	"history",
	"historymemes",
	"holup",
	"homeautomation",
	"homeimprovement",
	"homestead",
	"howto",
	"humansbeingbros",
	"iama",
	"idiotsincars",
	"indieheads",
	"interestingasfuck",
	"internetisbeautiful",
	"iphone",
	"itookapicture",
	"japantravel",
	"jokes",
	"keto",
	"kpop",
	"leagueoflegends",
	"learnprogramming",
	"lifehacks",
	"lifeprotips",
	"listentothis",
	"loseit",
	"mademesmile",
	"makeupaddiction",
	"malefashionadvice",
	"maliciouscompliance",
	"marvelmemes",
	"marvelstudios",
	"math",
	"maybemaybemaybe",
	"mealprepsunday",
	"meditation",
	"memes",
	"mildlyinfuriating",
	"mildlyinteresting",
	"minecraft",
	"minecraftmemes",
	"mma",
	"modernwarfareii",
	"motorcycles",
	"moviedetails",
	"movies",
	"music",
	"mypeopleneedme",
	"nails",
	"nasa",
	"natureisfuckinglit",
	"nba",
	"netflixbestof",
	"nevertellmetheodds",
	"news",
	"nfl",
	"nintendoswitch",
	"nosleep",
	"nostupidquestions",
	"nottheonion",
	"nutrition",
	"oddlysatisfying",
	"oddlyspecific",
	"offmychest",
	"oldschoolcool",
	"onepiece",
	"outdoors",
	"outoftheloop",
	"overwatch",
	"painting",
	"parenting",
	"pcgaming",
	"pcmasterrace",
	"personalfinance",
	"pettyrevenge",
	"philosophy",
	"photography",
	"photoshopbattles",
	"pics",
	"podcasts",
	"pokemon",
	"pokemongo",
	"politics",
	"popculturechat",
	"premierleague",
	"prequelmemes",
	"productivity",
	"programmerhumor",
	"programming",
	"ps4",
	"ps5",
	"psychology",
	"rarepuppers",
	"reactiongifs",
	"recipes",
	"relationship_advice",
	"relationshipmemes",
	"roadtrip",
	"running",
	"science",
	"sciencememes",
	"scifi",
	"shoestring",
	"showerthoughts",
	"singularity",
	"skincareaddiction",
	"slowcooking",
	"sneakers",
	"soccer",
	"socialskills",
	"solotravel",
	"space",
	"spacex",
	"sports",
	"standupshots",
	"starterpacks",
	"starwars",
	"steam",
	"stockmarket",
	"streetwear",
	"strength_training",
	"survival",
	"tattoos",
	"taylorswift",
	"technicallythetruth",
	"technology",
	"television",
	"teslamotors",
	"thriftstorehauls",
	"tifu",
	"tinder",
	"todayilearned",
	"travel",
	"travelhacks",
	"trippinthroughtime",
	"unexpected",
	"unitedkingdom",
	"unpopularopinion",
	"upliftingnews",
	"videos",
	"wallstreetbets",
	"watchpeopledieinside",
	"wearethemusicmakers",
	"wholesomememes",
	"woahdude",
	"woodworking",
	"worldnews",
	"writingprompts",
	"youshouldknow",
}

var (
	inputStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	inputContainerStyle = lipgloss.NewStyle().Margin(2, 4, 1, 4)
)

type (
	selectSubredditMsg struct{}
	hideSubredditMsg   struct{}
)

func ShowSubredditInput() tea.Cmd {
	return func() tea.Msg {
		return selectSubredditMsg{}
	}
}

func HideSubredditInput() tea.Cmd {
	return func() tea.Msg {
		return hideSubredditMsg{}
	}
}

type SubredditInput struct {
	model textinput.Model
	focus bool
}

func NewSubredditInput() SubredditInput {
	model := textinput.New()
	model.ShowSuggestions = true
	model.SetSuggestions(subredditSuggestions)
	model.CharLimit = 30

	return SubredditInput{model: model}
}

func (s *SubredditInput) Focus() tea.Cmd {
	s.focus = true
	s.model.Reset()
	return s.model.Focus()
}

func (s *SubredditInput) Blur() {
	s.focus = false
	s.model.Blur()
}

func (s SubredditInput) Init() tea.Cmd {
	return textinput.Blink
}

func (s SubredditInput) Update(msg tea.Msg) (SubredditInput, tea.Cmd) {
	if !s.focus {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			s.Blur()
			return s, focusListPage
		case "ctrl+c":
			return s, tea.Quit
		case "enter":
			s.Blur()
			return s, fetchSubredditPosts(s.model.Value())
		}
	}

	var cmd tea.Cmd
	s.model, cmd = s.model.Update(msg)
	return s, cmd
}

func (s SubredditInput) View() string {
	selectionView := fmt.Sprintf("Choose a subreddit:\n%s", inputStyle.Render(s.model.View()))
	return inputContainerStyle.Render(selectionView)
}
