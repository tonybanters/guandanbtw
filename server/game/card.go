package game

type Suit int

const (
	Suit_Hearts Suit = iota
	Suit_Diamonds
	Suit_Clubs
	Suit_Spades
	Suit_Joker
)

type Rank int

const (
	Rank_Two Rank = iota
	Rank_Three
	Rank_Four
	Rank_Five
	Rank_Six
	Rank_Seven
	Rank_Eight
	Rank_Nine
	Rank_Ten
	Rank_Jack
	Rank_Queen
	Rank_King
	Rank_Ace
	Rank_Black_Joker
	Rank_Red_Joker
)

type Card struct {
	Suit Suit
	Rank Rank
	Id   int
}

type Deck struct {
	Cards []Card
}

func New_Deck() *Deck {
	deck := &Deck{Cards: make([]Card, 0, 108)}
	id := 0

	for copy := 0; copy < 2; copy++ {
		for suit := Suit_Hearts; suit <= Suit_Spades; suit++ {
			for rank := Rank_Two; rank <= Rank_Ace; rank++ {
				deck.Cards = append(deck.Cards, Card{Suit: suit, Rank: rank, Id: id})
				id++
			}
		}
		deck.Cards = append(deck.Cards, Card{Suit: Suit_Joker, Rank: Rank_Black_Joker, Id: id})
		id++
		deck.Cards = append(deck.Cards, Card{Suit: Suit_Joker, Rank: Rank_Red_Joker, Id: id})
		id++
	}

	return deck
}

func (d *Deck) Shuffle() {
	for i := len(d.Cards) - 1; i > 0; i-- {
		j := rand_int(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
}

func (d *Deck) Deal() [4][]Card {
	var hands [4][]Card
	for i := 0; i < 4; i++ {
		hands[i] = make([]Card, 0, 27)
	}

	for i, card := range d.Cards {
		hands[i%4] = append(hands[i%4], card)
	}

	return hands
}

func rank_value(rank Rank, level Rank) int {
	if rank == Rank_Red_Joker {
		return 100
	}
	if rank == Rank_Black_Joker {
		return 99
	}
	if rank == level {
		return 98
	}

	base_order := []Rank{
		Rank_Two, Rank_Three, Rank_Four, Rank_Five, Rank_Six,
		Rank_Seven, Rank_Eight, Rank_Nine, Rank_Ten,
		Rank_Jack, Rank_Queen, Rank_King, Rank_Ace,
	}

	for i, r := range base_order {
		if r == rank {
			return i
		}
	}
	return -1
}

func card_value(card Card, level Rank) int {
	return rank_value(card.Rank, level)
}

func Card_Value(card Card, level Rank) int {
	return rank_value(card.Rank, level)
}

func Is_Wild(card Card, level Rank) bool {
	return card.Suit == Suit_Hearts && card.Rank == level
}

func compare_cards(a, b Card, level Rank) int {
	va := card_value(a, level)
	vb := card_value(b, level)
	if va > vb {
		return 1
	}
	if va < vb {
		return -1
	}
	return 0
}
