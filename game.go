// Package game provides an API for playing bridge to a user or AI agent.
package game

import (
	"math/rand"
	"sort"
)

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

func New() *Bridge {
	b := &Bridge{
		hands: make(map[PlayerName]Hand),
	}
	b.Deal()
	return b
}

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
	b.ewTricksTaken = 0
}
