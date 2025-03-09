// Package game provides an API for playing bridge to a user or AI agent.
package game

import (
	"fmt"
	"math/rand"
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

func (h Hand) Points() int {
	points := 0
	for _, card := range h {
		switch card.value {
		case Ace:
			points += 4
		case King:
			points += 3
		case Queen:
			points += 2
		case Jack:
			points += 1
		}
	}
	return points
}

func (h Hand) AlternatePoints() float32 {
	alternatePoints := float32(h.Points())
	// AK = 2
	// A = +1
	// K = +0.5
	suitPoints := make([]float32, 4)
	countBySuit := make([]int, 4)
	for _, card := range h {
		countBySuit[card.suit] += 1
		switch card.value {
		case Ace:
			if suitPoints[card.suit] != 0 {
				suitPoints[card.suit] = 2
				continue
			}
			suitPoints[card.suit] += 1
		case King:
			if suitPoints[card.suit] != 0 {
				suitPoints[card.suit] = 2
				continue
			}
			suitPoints[card.suit] += 0.5
		}
	}
	for _, suitPoint := range suitPoints {
		alternatePoints += suitPoint
	}

	// + longest two suits
	sort.Slice(countBySuit, func(i, j int) bool {
		return countBySuit[i] > countBySuit[j]
	})
	alternatePoints = alternatePoints + float32(countBySuit[0]) + float32(countBySuit[1])
	return alternatePoints
}

type Bridge struct {
	hands  map[PlayerName]Hand
	dealer PlayerName

	nsTricksTaken int
	esTricksTaken int
}

func (b *Bridge) String() string {
	s := ""
	for player := North; player <= West; player++ {
		s += fmt.Sprintf("%s:\n%v\n", PlayerName(player), b.hands[player])
	}
	return s
}

func New() *Bridge {
	b := &Bridge{
		hands: make(map[PlayerName]Hand),
	}
	b.Deal()
	return b
}

type Deck []Card

func (d *Deck) Shuffle() {
	rand.Shuffle(len(*d), func(i, j int) {
		(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
	})
}

func getFullDeck() Deck {
	deck := make([]Card, 0, 52)
	for suit := Club; suit <= Spade; suit++ {
		for value := Two; value <= Ace; value++ {
			deck = append(deck, Card{suit: suit, value: value})
		}
	}
	return deck
}

func newShuffledDeck() Deck {
	deck := getFullDeck()
	deck.Shuffle()
	return deck
}

func (b *Bridge) Deal() {
	b.Reset()
	deck := newShuffledDeck()

	// Technically, we could just assign in groups of 13, but... this feels
	// like dealing.
	for i := 0; i < len(deck); i++ {
		player := PlayerName(i % 4)
		b.hands[player] = append(b.hands[player], deck[i])
	}
	b.dealer = PlayerName((b.dealer + 1) % 4)
}

func (b *Bridge) Reset() {
	b.hands = make(map[PlayerName]Hand)
	b.nsTricksTaken = 0
	b.esTricksTaken = 0
}
