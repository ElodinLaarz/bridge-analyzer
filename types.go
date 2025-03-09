package game

import (
	"fmt"
	"sort"
)

type PlayerName int

const (
	North PlayerName = iota
	East
	South
	West
)

func (p PlayerName) String() string {
	switch p {
	case North:
		return "N"
	case East:
		return "E"
	case South:
		return "S"
	case West:
		return "W"
	default:
		return "?"
	}
}

type Suit int

const (
	Club Suit = iota
	Diamond
	Heart
	Spade
)

func (s Suit) String() string {
	switch s {
	case Club:
		return "C"
	case Diamond:
		return "D"
	case Heart:
		return "H"
	case Spade:
		return "S"
	default:
		return "?"
	}
}

type Value int

const (
	Two Value = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

func (v Value) String() string {
	switch v {
	case Two:
		return "2"
	case Three:
		return "3"
	case Four:
		return "4"
	case Five:
		return "5"
	case Six:
		return "6"
	case Seven:
		return "7"
	case Eight:
		return "8"
	case Nine:
		return "9"
	case Ten:
		return "10"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	default:
		return "?"
	}
}

type Card struct {
	suit  Suit
	value Value
}

type Hand []Card

func (h Hand) String() string {
	s := fmt.Sprintf(" HCP %d\n", h.Points())
	s += fmt.Sprintf(" AltPts %.1f\n", h.AlternatePoints())
	// Sort into suits
	sorted := make([][]Card, 4)
	for _, card := range h {
		sorted[card.suit] = append(sorted[card.suit], card)
	}
	// Sort each suit by value
	for suit := range sorted {
		sort.Slice(sorted[suit], func(i, j int) bool {
			return sorted[suit][i].value < sorted[suit][j].value
		})
		s += fmt.Sprintf(" %s: ", Suit(suit))
		for _, card := range sorted[suit] {
			s += fmt.Sprintf(" %s", card.value)
		}
		s += "\n"
	}
	return s
}

type Bridge struct {
	hands  map[PlayerName]Hand
	dealer PlayerName

	nsTricksTaken int
	ewTricksTaken int

	contractSuit *Suit
}

func (b *Bridge) String() string {
	s := ""
	for player := North; player <= West; player++ {
		s += fmt.Sprintf("%s:\n%v\n", PlayerName(player), b.hands[player])
	}
	return s
}

type Deck []Card
