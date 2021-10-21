package examples

import (
	"flag"
	"fmt"
	"os"

	"github.com/tada-team/tdproto"
)

type settings struct {
	Server        string
	Verbose       bool
	TeamUid       string
	Chat          string
	Token         string
	DryRun        bool
	requireTeam   bool
	requireChat   bool
	requireToken  bool
	requireDryRun bool
}

func NewSettings() settings { return settings{} }

func (s *settings) RequireTeam() {
	flag.StringVar(&s.TeamUid, "team", "", "team uid")
	s.requireTeam = true
}

func (s *settings) RequireChat() {
	flag.StringVar(&s.Chat, "chat", "", "chat jid")
	s.requireChat = true
}

func (s *settings) RequireToken() {
	flag.StringVar(&s.Token, "token", "", "bot or user token")
	s.requireToken = true
}

func (s *settings) RequireDryRun() {
	flag.BoolVar(&s.DryRun, "dryrun", false, "read or del pull")
	s.requireDryRun = true
}

func (s *settings) Parse() {
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
		} else if !tdproto.JID(s.Chat).Valid() {
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
