package room

import (
	"time"

	"guandanbtw/game"
	"guandanbtw/protocol"
)

type Play_Action struct {
	client   *Client
	card_ids []int
}

type Tribute_Action struct {
	client  *Client
	card_id int
}

type Room struct {
	id        string
	clients   [4]*Client
	game      *game.Game_State
	join      chan *Client
	leave     chan *Client
	play      chan Play_Action
	pass      chan *Client
	tribute   chan Tribute_Action
	fill_bots chan *Client
	bot_turn  chan int
}

func new_room(id string) *Room {
	return &Room{
		id:        id,
		join:      make(chan *Client),
		leave:    make(chan *Client),
		play:      make(chan Play_Action),
		pass:      make(chan *Client),
		tribute:   make(chan Tribute_Action),
		fill_bots: make(chan *Client),
		bot_turn:  make(chan int),
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.join:
			r.handle_join(client)
		case client := <-r.leave:
			r.handle_leave(client)
		case action := <-r.play:
			r.handle_play(action)
		case client := <-r.pass:
			r.handle_pass(client)
		case action := <-r.tribute:
			r.handle_tribute(action)
		case <-r.fill_bots:
			r.handle_fill_bots()
		case seat := <-r.bot_turn:
			r.handle_bot_turn(seat)
		}
	}
}

func (r *Room) handle_join(client *Client) {
	seat := r.find_empty_seat()
	if seat == -1 {
		client.send_error("room is full")
		return
	}

	r.clients[seat] = client
	client.room = r

	r.broadcast_room_state()

	if r.is_full() {
		r.start_game()
	}
}

func (r *Room) handle_leave(client *Client) {
	for i := range 4 {
		if r.clients[i] == client {
			r.clients[i] = nil
			break
		}
	}
	client.room = nil

	r.broadcast(&protocol.Message{
		Type: protocol.Msg_Player_Left,
		Payload: protocol.Player_Info{
			Id:   client.id,
			Name: client.name,
		},
	})

	r.broadcast_room_state()
}

func (r *Room) handle_play(action Play_Action) {
	if r.game == nil {
		return
	}

	seat := r.get_seat(action.client)
	if seat == -1 || seat != r.game.Current_Turn {
		action.client.send_error("not your turn")
		return
	}

	cards := r.game.Get_Cards_By_Id(seat, action.card_ids)
	if cards == nil {
		action.client.send_error("invalid cards")
		return
	}

	combo := game.Detect_Combination(cards, r.game.Level)
	if combo.Type == game.Comb_Invalid {
		action.client.send_error("invalid combination")
		return
	}

	if r.game.Current_Lead.Type != game.Comb_Invalid {
		if !game.Can_Beat(combo, r.game.Current_Lead) {
			action.client.send_error("cannot beat current play")
			return
		}
	}

	r.game.Remove_Cards(seat, action.card_ids)
	r.game.Current_Lead = combo
	r.game.Lead_Player = seat
	r.game.Pass_Count = 0

	r.broadcast(&protocol.Message{
		Type: protocol.Msg_Play_Made,
		Payload: protocol.Play_Made_Payload{
			Player_Id:  action.client.id,
			Seat:       seat,
			Cards:      cards,
			Combo_Type: combo_type_name(combo.Type),
			Is_Pass:    false,
		},
	})

	if len(r.game.Hands[seat]) == 0 {
		r.game.Finish_Order = append(r.game.Finish_Order, seat)
		if r.check_hand_end() {
			return
		}
	}

	r.advance_turn()
}

func (r *Room) handle_pass(client *Client) {
	if r.game == nil {
		return
	}

	seat := r.get_seat(client)
	if seat == -1 || seat != r.game.Current_Turn {
		client.send_error("not your turn")
		return
	}

	if r.game.Current_Lead.Type == game.Comb_Invalid {
		client.send_error("cannot pass when leading")
		return
	}

	r.game.Pass_Count++

	r.broadcast(&protocol.Message{
		Type: protocol.Msg_Play_Made,
		Payload: protocol.Play_Made_Payload{
			Player_Id: client.id,
			Seat:      seat,
			Is_Pass:   true,
		},
	})

	if r.game.Pass_Count >= 3 {
		r.game.Current_Lead = game.Combination{Type: game.Comb_Invalid}
		next_leader := r.game.Lead_Player
		if r.is_finished(next_leader) {
			teammate := (next_leader + 2) % 4
			next_leader = teammate
		}
		r.game.Current_Turn = next_leader
		r.game.Pass_Count = 0
		r.send_turn_notification()
		r.trigger_bot_turn_if_needed()
		return
	}

	r.advance_turn()
}

func (r *Room) handle_tribute(action Tribute_Action) {
	if r.game == nil || r.game.Phase != game.Phase_Tribute {
		return
	}

	seat := r.get_seat(action.client)
	tribute_info := r.game.Get_Tribute_Info(seat)
	if tribute_info == nil {
		action.client.send_error("you don't need to give tribute")
		return
	}

	card := r.game.Get_Card_By_Id(seat, action.card_id)
	if card == nil {
		action.client.send_error("invalid card")
		return
	}

	if game.Is_Wild(*card, r.game.Level) {
		action.client.send_error("cannot tribute wild cards")
		return
	}

	r.game.Remove_Cards(seat, []int{action.card_id})
	r.game.Hands[tribute_info.To_Seat] = append(r.game.Hands[tribute_info.To_Seat], *card)

	if r.clients[tribute_info.To_Seat] != nil {
		r.clients[tribute_info.To_Seat].send_message(&protocol.Message{
			Type: protocol.Msg_Tribute_Recv,
			Payload: protocol.Tribute_Recv_Payload{
				Card: *card,
			},
		})
	}

	r.game.Mark_Tribute_Done(seat)

	if r.game.All_Tributes_Done() {
		r.game.Phase = game.Phase_Play
		r.game.Current_Turn = r.game.Tribute_Leader
		r.send_turn_notification()
		r.trigger_bot_turn_if_needed()
	}
}

func (r *Room) start_game() {
	r.game = game.New_Game_State()

	deck := game.New_Deck()
	deck.Shuffle()
	hands := deck.Deal()

	for i := 0; i < 4; i++ {
		r.game.Hands[i] = hands[i]
	}

	for i := 0; i < 4; i++ {
		if r.clients[i] != nil {
			r.clients[i].send_message(&protocol.Message{
				Type: protocol.Msg_Deal_Cards,
				Payload: protocol.Deal_Cards_Payload{
					Cards: r.game.Hands[i],
					Level: r.game.Level,
				},
			})
		}
	}

	r.game.Phase = game.Phase_Play
	r.game.Current_Turn = 0
	r.send_turn_notification()
	r.trigger_bot_turn_if_needed()
}

func (r *Room) check_hand_end() bool {
	if len(r.game.Finish_Order) < 2 {
		return false
	}

	first := r.game.Finish_Order[0]
	second := r.game.Finish_Order[1]

	first_team := first % 2
	second_team := second % 2

	if first_team == second_team {
		r.end_hand(first_team, r.calculate_level_advance())
		return true
	}

	if len(r.game.Finish_Order) >= 3 {
		third := r.game.Finish_Order[2]
		third_team := third % 2

		winning_team := first_team
		r.end_hand(winning_team, r.calculate_level_advance())
		return true
		_ = third_team
	}

	return false
}

func (r *Room) calculate_level_advance() int {
	if len(r.game.Finish_Order) < 2 {
		return 0
	}

	first := r.game.Finish_Order[0]
	first_team := first % 2

	for i, seat := range r.game.Finish_Order {
		if seat%2 != first_team {
			partner_pos := -1
			for j, s := range r.game.Finish_Order {
				if s%2 == first_team && j != 0 {
					partner_pos = j
					break
				}
			}

			if partner_pos == -1 {
				partner_pos = 3
			}

			switch {
			case i == 3 && partner_pos == 1:
				return 4
			case i == 2 && partner_pos == 1:
				return 2
			default:
				return 1
			}
		}
	}

	return 4
}

func (r *Room) end_hand(winning_team int, level_advance int) {
	old_level := r.game.Team_Levels[winning_team]
	new_level := old_level + level_advance
	if new_level > 12 {
		new_level = 12
	}
	r.game.Team_Levels[winning_team] = new_level

	r.broadcast(&protocol.Message{
		Type: protocol.Msg_Hand_End,
		Payload: protocol.Hand_End_Payload{
			Finish_Order:  seats_to_ids(r.game.Finish_Order, r.clients),
			Winning_Team:  winning_team,
			Level_Advance: level_advance,
			New_Levels:    r.game.Team_Levels,
		},
	})

	if new_level >= 12 && game.Rank(old_level) == game.Rank_Ace {
		r.broadcast(&protocol.Message{
			Type: protocol.Msg_Game_End,
			Payload: protocol.Game_End_Payload{
				Winning_Team: winning_team,
				Final_Levels: r.game.Team_Levels,
			},
		})
		return
	}

	r.setup_tribute()
}

func (r *Room) setup_tribute() {
	r.game.Setup_Tributes()

	if len(r.game.Tributes) == 0 {
		r.start_new_hand()
		return
	}

	for _, t := range r.game.Tributes {
		if r.clients[t.From_Seat] != nil {
			r.clients[t.From_Seat].send_message(&protocol.Message{
				Type: protocol.Msg_Tribute,
				Payload: protocol.Tribute_Payload{
					From_Seat: t.From_Seat,
					To_Seat:   t.To_Seat,
				},
			})
		}
	}

	r.game.Phase = game.Phase_Tribute
}

func (r *Room) start_new_hand() {
	r.game.Reset_Hand()

	deck := game.New_Deck()
	deck.Shuffle()
	hands := deck.Deal()

	for i := 0; i < 4; i++ {
		r.game.Hands[i] = hands[i]
	}

	for i := 0; i < 4; i++ {
		if r.clients[i] != nil {
			r.clients[i].send_message(&protocol.Message{
				Type: protocol.Msg_Deal_Cards,
				Payload: protocol.Deal_Cards_Payload{
					Cards: r.game.Hands[i],
					Level: r.game.Level,
				},
			})
		}
	}

	r.game.Phase = game.Phase_Play
	r.game.Current_Turn = r.game.Tribute_Leader
	r.send_turn_notification()
	r.trigger_bot_turn_if_needed()
}

func (r *Room) advance_turn() {
	for i := 1; i <= 4; i++ {
		next := (r.game.Current_Turn + i) % 4
		if !r.is_finished(next) {
			r.game.Current_Turn = next
			r.send_turn_notification()
			r.trigger_bot_turn_if_needed()
			return
		}
	}
}

func (r *Room) is_finished(seat int) bool {
	for _, s := range r.game.Finish_Order {
		if s == seat {
			return true
		}
	}
	return false
}

func (r *Room) send_turn_notification() {
	can_pass := r.game.Current_Lead.Type != game.Comb_Invalid

	r.broadcast(&protocol.Message{
		Type: protocol.Msg_Turn,
		Payload: protocol.Turn_Payload{
			Player_Id: r.clients[r.game.Current_Turn].id,
			Seat:      r.game.Current_Turn,
			Can_Pass:  can_pass,
		},
	})
}

func (r *Room) broadcast(msg *protocol.Message) {
	for _, client := range r.clients {
		if client != nil {
			client.send_message(msg)
		}
	}
}

func (r *Room) broadcast_room_state() {
	players := make([]protocol.Player_Info, 0)
	for i, c := range r.clients {
		if c != nil {
			players = append(players, protocol.Player_Info{
				Id:   c.id,
				Name: c.name,
				Seat: i,
				Team: i % 2,
			})
		}
	}

	for _, client := range r.clients {
		if client != nil {
			client.send_message(&protocol.Message{
				Type: protocol.Msg_Room_State,
				Payload: protocol.Room_State_Payload{
					Room_Id:     r.id,
					Players:     players,
					Game_Active: r.game != nil,
					Your_Id:     client.id,
				},
			})
		}
	}
}

func (r *Room) find_empty_seat() int {
	for i := 0; i < 4; i++ {
		if r.clients[i] == nil {
			return i
		}
	}
	return -1
}

func (r *Room) is_full() bool {
	for _, c := range r.clients {
		if c == nil {
			return false
		}
	}
	return true
}

func (r *Room) get_seat(client *Client) int {
	for i, c := range r.clients {
		if c == client {
			return i
		}
	}
	return -1
}

func combo_type_name(t game.Combination_Type) string {
	names := map[game.Combination_Type]string{
		game.Comb_Single:     "single",
		game.Comb_Pair:       "pair",
		game.Comb_Triple:     "triple",
		game.Comb_Full_House: "full_house",
		game.Comb_Straight:   "straight",
		game.Comb_Tube:       "tube",
		game.Comb_Plate:      "plate",
		game.Comb_Bomb:       "bomb",
	}
	return names[t]
}

func seats_to_ids(seats []int, clients [4]*Client) []string {
	ids := make([]string, len(seats))
	for i, seat := range seats {
		if clients[seat] != nil {
			ids[i] = clients[seat].id
		}
	}
	return ids
}

func (r *Room) handle_fill_bots() {
	bot_idx := 0
	for i := 0; i < 4; i++ {
		if r.clients[i] == nil && bot_idx < len(bot_names) {
			bot := new_bot(generate_id(), bot_names[bot_idx])
			bot.room = r
			r.clients[i] = bot
			bot_idx++
		}
	}

	r.broadcast_room_state()

	if r.is_full() {
		r.start_game()
	}
}

func (r *Room) handle_bot_turn(seat int) {
	if r.game == nil || r.game.Current_Turn != seat {
		return
	}

	client := r.clients[seat]
	if client == nil || !client.is_bot {
		return
	}

	time.Sleep(1500 * time.Millisecond)

	hand := r.game.Hands[seat]
	if len(hand) == 0 {
		return
	}

	if r.game.Current_Lead.Type != game.Comb_Invalid {
		lead_team := r.game.Lead_Player % 2
		bot_team := seat % 2
		if lead_team == bot_team {
			r.handle_pass(client)
			return
		}

		play := r.find_lowest_valid_play(seat)
		if play == nil {
			r.handle_pass(client)
			return
		}

		r.handle_play(Play_Action{
			client:   client,
			card_ids: play,
		})
		return
	}

	card_ids := []int{hand[0].Id}
	r.handle_play(Play_Action{
		client:   client,
		card_ids: card_ids,
	})
}

func (r *Room) find_lowest_valid_play(seat int) []int {
	hand := r.game.Hands[seat]
	lead := r.game.Current_Lead

	sorted_hand := make([]game.Card, len(hand))
	copy(sorted_hand, hand)
	for i := 0; i < len(sorted_hand)-1; i++ {
		for j := i + 1; j < len(sorted_hand); j++ {
			vi := game.Card_Value(sorted_hand[i], r.game.Level)
			vj := game.Card_Value(sorted_hand[j], r.game.Level)
			if vi > vj {
				sorted_hand[i], sorted_hand[j] = sorted_hand[j], sorted_hand[i]
			}
		}
	}

	if lead.Type == game.Comb_Single {
		for _, card := range sorted_hand {
			combo := game.Detect_Combination([]game.Card{card}, r.game.Level)
			if game.Can_Beat(combo, lead) {
				return []int{card.Id}
			}
		}
	}

	if lead.Type == game.Comb_Pair {
		rank_cards := make(map[game.Rank][]game.Card)
		for _, card := range sorted_hand {
			rank_cards[card.Rank] = append(rank_cards[card.Rank], card)
		}

		var valid_pairs [][]game.Card
		for _, cards := range rank_cards {
			if len(cards) >= 2 {
				combo := game.Detect_Combination(cards[:2], r.game.Level)
				if game.Can_Beat(combo, lead) {
					valid_pairs = append(valid_pairs, cards[:2])
				}
			}
		}

		if len(valid_pairs) > 0 {
			lowest := valid_pairs[0]
			lowest_val := game.Card_Value(lowest[0], r.game.Level)
			for _, pair := range valid_pairs[1:] {
				val := game.Card_Value(pair[0], r.game.Level)
				if val < lowest_val {
					lowest = pair
					lowest_val = val
				}
			}
			return []int{lowest[0].Id, lowest[1].Id}
		}
	}

	if lead.Type == game.Comb_Triple {
		rank_cards := make(map[game.Rank][]game.Card)
		for _, card := range sorted_hand {
			rank_cards[card.Rank] = append(rank_cards[card.Rank], card)
		}

		var valid_triples [][]game.Card
		for _, cards := range rank_cards {
			if len(cards) >= 3 {
				combo := game.Detect_Combination(cards[:3], r.game.Level)
				if game.Can_Beat(combo, lead) {
					valid_triples = append(valid_triples, cards[:3])
				}
			}
		}

		if len(valid_triples) > 0 {
			lowest := valid_triples[0]
			lowest_val := game.Card_Value(lowest[0], r.game.Level)
			for _, triple := range valid_triples[1:] {
				val := game.Card_Value(triple[0], r.game.Level)
				if val < lowest_val {
					lowest = triple
					lowest_val = val
				}
			}
			return []int{lowest[0].Id, lowest[1].Id, lowest[2].Id}
		}
	}

	return nil
}

func (r *Room) trigger_bot_turn_if_needed() {
	if r.game == nil {
		return
	}

	seat := r.game.Current_Turn
	client := r.clients[seat]
	if client != nil && client.is_bot {
		go func() {
			r.bot_turn <- seat
		}()
	}
}
