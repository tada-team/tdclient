package examples

import (
	"flag"
	"fmt"
	"os"

	"github.com/tada-team/tdproto"
)

type Settings struct {
	Server        string
	Verbose       bool
	TeamUid       string
	Chat          string
	Token         string
	DryRun        bool
	Deep          int
	requireTeam   bool
	requireChat   bool
	requireToken  bool
	requireDryRun bool
	requireDeep   int
}

func NewSettings() Settings { return Settings{} }

func (s *Settings) RequireTeam() {
	flag.StringVar(&s.TeamUid, "team", "", "team uid")
	s.requireTeam = true
}

func (s *Settings) RequireChat() {
	flag.StringVar(&s.Chat, "chat", "", "chat jid")
	s.requireChat = true
}

func (s *Settings) RequireToken() {
	flag.StringVar(&s.Token, "token", "", "bot or user token")
	s.requireToken = true
}

func (s *Settings) RequireDryRun() {
	flag.BoolVar(&s.DryRun, "dryrun", false, "read or del pull")
	s.requireDryRun = true
}

func (s *Settings) RequireDeep() {
	flag.IntVar(&s.Deep, "deep", 5, "only")
	s.requireDryRun = true
}

func (s *Settings) Parse() {
	flag.StringVar(&s.Server, "server", "https://web.tada.team", "server address")
	flag.BoolVar(&s.Verbose, "verbose", false, "verbose logging")
	flag.Parse()

	ok := true

	if s.requireTeam {
		if s.TeamUid == "" {
			fmt.Println("-team required")
			ok = false
		} else if !tdproto.ValidUid(s.TeamUid) {
			fmt.Println("invalid team uid")
			ok = false
		}
	}

	if s.requireChat {
		if s.Chat == "" {
			fmt.Println("-chat required")
			ok = false
		} else if !tdproto.NewJID(s.Chat).Valid() {
			fmt.Println("invalid chat")
			ok = false
		}
	}

	if s.requireToken && s.Token == "" {
		fmt.Println("-token required")
		ok = false
	}

	if !ok {
		os.Exit(0)
	}
}
