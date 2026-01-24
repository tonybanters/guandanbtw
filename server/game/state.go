package game

type Game_Phase int

const (
	Phase_Waiting Game_Phase = iota
	Phase_Deal
	Phase_Play
	Phase_Tribute
	Phase_End
)

type Tribute_Info struct {
	From_Seat int
	To_Seat   int
	Done      bool
}

type Game_State struct {
	Phase          Game_Phase
	Level          Rank
	Team_Levels    [2]int
	Hands          [4][]Card
	Current_Turn   int
	Current_Lead   Combination
	Lead_Player    int
	Pass_Count     int
	Finish_Order   []int
	Tributes       []Tribute_Info
	Tribute_Leader int
}

func New_Game_State() *Game_State {
	return &Game_State{
		Phase:        Phase_Waiting,
		Level:        Rank_Two,
		Team_Levels:  [2]int{0, 0},
		Finish_Order: make([]int, 0, 4),
	}
}

func (g *Game_State) Get_Cards_By_Id(seat int, ids []int) []Card {
	id_set := make(map[int]bool)
	for _, id := range ids {
		id_set[id] = true
	}

	var cards []Card
	for _, card := range g.Hands[seat] {
		if id_set[card.Id] {
			cards = append(cards, card)
			delete(id_set, card.Id)
		}
	}

	if len(id_set) > 0 {
		return nil
	}

	return cards
}

func (g *Game_State) Get_Card_By_Id(seat int, id int) *Card {
	for i := range g.Hands[seat] {
		if g.Hands[seat][i].Id == id {
			return &g.Hands[seat][i]
		}
	}
	return nil
}

func (g *Game_State) Remove_Cards(seat int, ids []int) {
	id_set := make(map[int]bool)
	for _, id := range ids {
		id_set[id] = true
	}

	var remaining []Card
	for _, card := range g.Hands[seat] {
		if !id_set[card.Id] {
			remaining = append(remaining, card)
		}
	}
	g.Hands[seat] = remaining
}

func (g *Game_State) Setup_Tributes() {
	g.Tributes = nil

	if len(g.Finish_Order) < 2 {
		return
	}

	first := g.Finish_Order[0]
	winning_team := first % 2

	var last_loser int = -1
	var second_last_loser int = -1

	for i := 3; i >= 0; i-- {
		seat := -1
		if i < len(g.Finish_Order) {
			seat = g.Finish_Order[i]
		} else {
			for s := 0; s < 4; s++ {
				found := false
				for _, f := range g.Finish_Order {
					if f == s {
						found = true
						break
					}
				}
				if !found {
					seat = s
					break
				}
			}
		}

		if seat%2 != winning_team {
			if last_loser == -1 {
				last_loser = seat
			} else if second_last_loser == -1 {
				second_last_loser = seat
				break
			}
		}
	}

	first_winner := g.Finish_Order[0]
	var second_winner int = -1
	for _, seat := range g.Finish_Order {
		if seat%2 == winning_team && seat != first_winner {
			second_winner = seat
			break
		}
	}

	if last_loser != -1 {
		g.Tributes = append(g.Tributes, Tribute_Info{
			From_Seat: last_loser,
			To_Seat:   first_winner,
		})
	}

	if g.is_double_win() && second_last_loser != -1 && second_winner != -1 {
		g.Tributes = append(g.Tributes, Tribute_Info{
			From_Seat: second_last_loser,
			To_Seat:   second_winner,
		})
	}

	g.Tribute_Leader = first_winner
}

func (g *Game_State) is_double_win() bool {
	if len(g.Finish_Order) < 2 {
		return false
	}
	return g.Finish_Order[0]%2 == g.Finish_Order[1]%2
}

func (g *Game_State) Get_Tribute_Info(seat int) *Tribute_Info {
	for i := range g.Tributes {
		if g.Tributes[i].From_Seat == seat && !g.Tributes[i].Done {
			return &g.Tributes[i]
		}
	}
	return nil
}

func (g *Game_State) Mark_Tribute_Done(seat int) {
	for i := range g.Tributes {
		if g.Tributes[i].From_Seat == seat {
			g.Tributes[i].Done = true
			break
		}
	}
}

func (g *Game_State) All_Tributes_Done() bool {
	for _, t := range g.Tributes {
		if !t.Done {
			return false
		}
	}
	return true
}

func (g *Game_State) Reset_Hand() {
	g.Current_Lead = Combination{Type: Comb_Invalid}
	g.Lead_Player = 0
	g.Pass_Count = 0
	g.Finish_Order = g.Finish_Order[:0]
	g.Tributes = nil

	winning_team := g.Tribute_Leader % 2
	g.Level = Rank(g.Team_Levels[winning_team])
}
