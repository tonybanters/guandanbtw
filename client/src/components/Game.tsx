import { motion } from 'framer-motion'
import { Card as Card_Type, Rank, get_rank_symbol, get_suit_symbol, is_red_suit, Suit_Joker, Rank_Red_Joker } from '../game/types'
import { Hand } from './Hand'
import { Table } from './Table'
import { Card_Back } from './Card'

interface Play_Log_Entry {
  seat: number
  cards: Card_Type[]
  combo_type: string
  is_pass: boolean
}

interface Game_Props {
  hand: Card_Type[]
  level: Rank
  selected_ids: Set<number>
  on_card_click: (id: number) => void
  on_play: () => void
  on_pass: () => void
  table_cards: Card_Type[]
  combo_type: string
  current_turn: number
  my_seat: number
  can_pass: boolean
  player_card_counts: number[]
  team_levels: [number, number]
  play_log: Play_Log_Entry[]
  players_map: Record<number, string>
  last_play_seat: number | null
}

export function Game({
  hand,
  level,
  selected_ids,
  on_card_click,
  on_play,
  on_pass,
  table_cards,
  combo_type,
  current_turn,
  my_seat,
  can_pass,
  player_card_counts,
  team_levels,
  play_log,
  players_map,
  last_play_seat,
}: Game_Props) {
  const is_my_turn = current_turn === my_seat
  const relative_positions = get_relative_positions(my_seat)

  return (
    <div style={styles.container}>
      <div style={styles.info_bar}>
        <div style={styles.level_badge}>
          Level: {get_rank_symbol(level)}
        </div>
        <div style={styles.team_scores}>
          <span style={{ color: '#2196f3' }}>Team 1: {get_rank_symbol(team_levels[0] as Rank)}</span>
          <span style={{ marginLeft: 16, color: '#e91e63' }}>Team 2: {get_rank_symbol(team_levels[1] as Rank)}</span>
        </div>
      </div>

      <div style={styles.main_layout}>
        <div style={styles.game_area}>
          <div style={styles.opponent_top}>
            <Opponent_Hand
              count={player_card_counts[relative_positions.top]}
              is_turn={current_turn === relative_positions.top}
              just_played={last_play_seat === relative_positions.top}
              seat={relative_positions.top}
              name={players_map[relative_positions.top]}
            />
          </div>

          <div style={styles.middle_row}>
            <div style={styles.opponent_side}>
              <Opponent_Hand
                count={player_card_counts[relative_positions.left]}
                is_turn={current_turn === relative_positions.left}
                just_played={last_play_seat === relative_positions.left}
                seat={relative_positions.left}
                vertical
                name={players_map[relative_positions.left]}
              />
            </div>

            <div style={styles.table_area}>
              <Table cards={table_cards} level={level} combo_type={combo_type} last_play_seat={last_play_seat} />
            </div>

            <div style={styles.opponent_side}>
            <Opponent_Hand
              count={player_card_counts[relative_positions.right]}
              is_turn={current_turn === relative_positions.right}
              just_played={last_play_seat === relative_positions.right}
              seat={relative_positions.right}
              vertical
              name={players_map[relative_positions.right]}
            />
          </div>
          </div>

          <div style={styles.my_area}>
          <Hand
            cards={hand}
            level={level}
            selected_ids={selected_ids}
            on_card_click={on_card_click}
          />

          <div style={styles.actions}>
            <motion.button
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={on_play}
              disabled={!is_my_turn || selected_ids.size === 0 || hand.length === 0}
              style={{
                ...styles.action_button,
                backgroundColor: is_my_turn && selected_ids.size > 0 && hand.length > 0 ? '#28a745' : '#444',
              }}
            >
              Play
            </motion.button>
            <motion.button
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={on_pass}
              disabled={!is_my_turn || !can_pass || hand.length === 0}
              style={{
                ...styles.action_button,
                backgroundColor: is_my_turn && can_pass && hand.length > 0 ? '#dc3545' : '#444',
              }}
            >
              Pass
            </motion.button>
          </div>

          {is_my_turn && hand.length > 0 && (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              style={styles.turn_indicator}
            >
              Your turn!
            </motion.div>
          )}
          {hand.length === 0 && (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              style={{ ...styles.turn_indicator, backgroundColor: '#28a745' }}
            >
              You finished!
            </motion.div>
          )}
          </div>
        </div>

        <div style={styles.play_log}>
          <div style={styles.play_log_title}>Play Log</div>
          {play_log.map((entry, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              style={styles.play_log_entry}
            >
              <span style={{ color: entry.seat % 2 === 0 ? '#2196f3' : '#e91e63', fontWeight: 'bold' }}>
                {players_map[entry.seat] || `Seat ${entry.seat + 1}`}:
              </span>{' '}
              {entry.is_pass ? (
                <span style={{ color: '#888' }}>Pass</span>
              ) : (
                <span>
                  {entry.cards.map((c, j) => (
                    <span key={j} style={{ color: c.Suit === Suit_Joker ? (c.Rank === Rank_Red_Joker ? '#dc3545' : '#000') : is_red_suit(c.Suit) ? '#dc3545' : '#000' }}>
                      {get_rank_symbol(c.Rank)}{get_suit_symbol(c.Suit)}{j < entry.cards.length - 1 ? ' ' : ''}
                    </span>
                  ))}
                  {entry.combo_type && <span style={{ color: '#888', marginLeft: 4 }}>({entry.combo_type})</span>}
                </span>
              )}
            </motion.div>
          ))}
        </div>
      </div>
    </div>
  )
}

interface Opponent_Hand_Props {
  count: number
  is_turn: boolean
  just_played?: boolean
  seat: number
  vertical?: boolean
  name?: string
}

function Opponent_Hand({ count, is_turn, just_played, seat, vertical, name }: Opponent_Hand_Props) {
  const display_count = Math.min(count, 10)
  const overlap = vertical ? 15 : 20

  const get_highlight_style = () => {
    if (just_played) return { backgroundColor: 'rgba(76, 175, 80, 0.3)', border: '2px solid #4caf50' }
    if (is_turn) return { backgroundColor: 'rgba(255,193,7,0.2)', border: '2px solid #ffc107' }
    return { backgroundColor: 'transparent', border: '2px solid transparent' }
  }

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: vertical ? 'column' : 'row',
        alignItems: 'center',
        gap: 8,
        padding: 8,
        borderRadius: 8,
        transition: 'all 0.2s ease',
        ...get_highlight_style(),
      }}
    >
      <div
        style={{
          display: 'flex',
          flexDirection: vertical ? 'column' : 'row',
          position: 'relative',
          width: vertical ? 50 : 50 + (display_count - 1) * overlap,
          height: vertical ? 70 + (display_count - 1) * overlap : 70,
        }}
      >
        {Array.from({ length: display_count }).map((_, i) => (
          <div
            key={i}
            style={{
              position: 'absolute',
              left: vertical ? 0 : i * overlap,
              top: vertical ? i * overlap : 0,
              transform: 'scale(0.7)',
              transformOrigin: 'top left',
            }}
          >
            <Card_Back />
          </div>
        ))}
      </div>
      <div style={{ color: '#fff', fontSize: 12 }}>
        {name || `Seat ${seat + 1}`}: {count}
      </div>
    </div>
  )
}

function get_relative_positions(my_seat: number) {
  return {
    top: (my_seat + 2) % 4,
    left: (my_seat + 1) % 4,
    right: (my_seat + 3) % 4,
  }
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    display: 'flex',
    flexDirection: 'column',
    height: '100vh',
    backgroundColor: '#0f3460',
  },
  info_bar: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '12px 24px',
    backgroundColor: '#16213e',
  },
  level_badge: {
    padding: '8px 16px',
    backgroundColor: '#ffc107',
    color: '#000',
    borderRadius: 8,
    fontWeight: 'bold',
  },
  team_scores: {
    color: '#fff',
    fontSize: 14,
  },
  game_area: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    padding: 20,
  },
  opponent_top: {
    display: 'flex',
    justifyContent: 'center',
    marginBottom: 20,
  },
  middle_row: {
    flex: 1,
    display: 'flex',
    alignItems: 'center',
  },
  opponent_side: {
    width: 120,
    display: 'flex',
    justifyContent: 'center',
  },
  table_area: {
    flex: 1,
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: 'rgba(0,0,0,0.2)',
    borderRadius: 16,
    margin: '0 20px',
  },
  my_area: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    paddingTop: 20,
    borderTop: '2px solid #333',
  },
  actions: {
    display: 'flex',
    gap: 16,
    marginTop: 16,
  },
  action_button: {
    padding: '12px 32px',
    fontSize: 16,
    border: 'none',
    borderRadius: 8,
    color: '#fff',
    cursor: 'pointer',
  },
  turn_indicator: {
    marginTop: 12,
    padding: '8px 16px',
    backgroundColor: '#ffc107',
    color: '#000',
    borderRadius: 8,
    fontWeight: 'bold',
  },
  main_layout: {
    display: 'flex',
    flex: 1,
    overflow: 'hidden',
  },
  play_log: {
    width: 220,
    backgroundColor: '#16213e',
    padding: 12,
    overflowY: 'auto',
    borderLeft: '2px solid #333',
  },
  play_log_title: {
    color: '#fff',
    fontSize: 14,
    fontWeight: 'bold',
    marginBottom: 12,
    paddingBottom: 8,
    borderBottom: '1px solid #333',
  },
  play_log_entry: {
    fontSize: 12,
    color: '#fff',
    padding: '6px 0',
    borderBottom: '1px solid rgba(255,255,255,0.1)',
  },
}
