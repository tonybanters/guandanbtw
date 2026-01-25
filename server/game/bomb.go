package game

func detect_bomb(cards []Card, level Rank) Combination {
	n := len(cards)
	if n == 4 && is_four_joker_bomb(cards) {
		return Combination{
			Type:       Comb_Bomb,
			Cards:      cards,
			Bomb_Power: 1000,
		}
	}
	if is_straight_flush(cards, level) {
		return Combination{
			Type:       Comb_Bomb,
			Cards:      cards,
			Bomb_Power: 900 + straight_flush_value(cards, level),
		}
	}
	if n >= 4 && n <= 10 {
		if power, ok := is_n_of_kind_bomb(cards, level); ok {
			return Combination{
				Type:       Comb_Bomb,
				Cards:      cards,
				Bomb_Power: power,
			}
		}
	}

	return Combination{Type: Comb_Invalid}
}

func is_four_joker_bomb(cards []Card) bool {
	if len(cards) != 4 {
		return false
	}

	red_count := 0
	black_count := 0

	for _, c := range cards {
		switch c.Rank {
		case Rank_Red_Joker:
			red_count++
		case Rank_Black_Joker:
			black_count++
		default:
			return false
		}
	}

	return red_count == 2 && black_count == 2
}

func is_straight_flush(cards []Card, level Rank) bool {
	if len(cards) < 5 {
		return false
	}

	non_wild, wild := separate_wilds(cards, level)

	if len(non_wild) == 0 {
		return false
	}

	var suit Suit = -1
	for _, c := range non_wild {
		if c.Rank == Rank_Black_Joker || c.Rank == Rank_Red_Joker {
			return false
		}
		if suit == -1 {
			suit = c.Suit
		} else if c.Suit != suit {
			return false
		}
	}

	rank_present := make(map[Rank]bool)
	for _, c := range non_wild {
		rank_present[c.Rank] = true
	}

	natural_order := []Rank{
		Rank_Ace, Rank_Two, Rank_Three, Rank_Four, Rank_Five,
		Rank_Six, Rank_Seven, Rank_Eight, Rank_Nine, Rank_Ten,
		Rank_Jack, Rank_Queen, Rank_King, Rank_Ace,
	}

	needed_len := len(cards)
	wilds_available := len(wild)

	for start := 0; start <= len(natural_order)-needed_len; start++ {
		gaps := 0
		for i := 0; i < needed_len; i++ {
			rank := natural_order[start+i]
			if !rank_present[rank] {
				gaps++
			}
		}
		if gaps <= wilds_available {
			return true
		}
	}

	return false
}

func straight_flush_value(cards []Card, level Rank) int {
	non_wild, _ := separate_wilds(cards, level)

	max_val := 0
	for _, c := range non_wild {
		v := rank_value(c.Rank, level)
		if v > max_val {
			max_val = v
		}
	}

	return len(cards)*10 + max_val
}

func is_n_of_kind_bomb(cards []Card, level Rank) (int, bool) {
	n := len(cards)
	if n < 4 || n > 10 {
		return 0, false
	}

	non_wild, wild := separate_wilds(cards, level)

	rank_counts := count_ranks(non_wild)

	for rank, count := range rank_counts {
		if rank == Rank_Black_Joker || rank == Rank_Red_Joker {
			continue
		}
		total := count + len(wild)
		if total >= n {
			power := n*100 + rank_value(rank, level)
			return power, true
		}
	}

	if len(wild) >= n {
		power := n*100 + rank_value(level, level)
		return power, true
	}

	return 0, false
}
