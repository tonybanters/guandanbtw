import { motion } from 'framer-motion'
import { Card as Card_Type, get_suit_symbol, get_rank_symbol, is_red_suit, is_wild, Rank, Rank_Red_Joker, Suit_Joker } from '../game/types'

interface Card_Props {
  card: Card_Type
  level: Rank
  selected: boolean
  on_click: () => void
}

export function Card({ card, level, selected, on_click }: Card_Props) {
  const is_joker = card.Suit === Suit_Joker
  const is_red = is_joker ? card.Rank === Rank_Red_Joker : is_red_suit(card.Suit)
  const is_wild_card = is_wild(card, level)

  return (
    <motion.div
      onClick={on_click}
      animate={{
        y: selected ? -20 : 0,
        scale: selected ? 1.05 : 1,
      }}
      whileHover={{ scale: 1.08 }}
      transition={{ type: 'spring', stiffness: 400, damping: 25 }}
      style={{
        width: 70,
        height: 100,
        backgroundColor: is_wild_card ? '#fff3cd' : '#fff',
        border: is_wild_card ? '3px solid #ffc107' : '2px solid #333',
        borderRadius: 8,
        position: 'relative',
        cursor: 'pointer',
        userSelect: 'none',
        boxShadow: selected ? '0 8px 16px rgba(0,0,0,0.3)' : '0 2px 4px rgba(0,0,0,0.1)',
        color: is_red ? '#dc3545' : '#000',
        fontWeight: 'bold',
      }}
    >
      <div style={{
        position: 'absolute',
        top: 4,
        left: 6,
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        lineHeight: 1,
      }}>
        <span style={{ fontSize: 16, fontWeight: 'bold' }}>{is_joker ? (card.Rank === Rank_Red_Joker ? 'R' : 'B') : get_rank_symbol(card.Rank)}</span>
        <span style={{ fontSize: 18 }}>{is_joker ? '🃏' : get_suit_symbol(card.Suit)}</span>
      </div>
      <div style={{
        position: 'absolute',
        bottom: 4,
        right: 6,
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        lineHeight: 1,
        transform: 'rotate(180deg)',
      }}>
        <span style={{ fontSize: 14 }}>{is_joker ? (card.Rank === Rank_Red_Joker ? 'R' : 'B') : get_rank_symbol(card.Rank)}</span>
        <span style={{ fontSize: 12 }}>{is_joker ? '🃏' : get_suit_symbol(card.Suit)}</span>
      </div>
      <div style={{
        position: 'absolute',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
        fontSize: is_joker ? 28 : 24,
      }}>
        {is_joker ? '🃏' : get_suit_symbol(card.Suit)}
      </div>
    </motion.div>
  )
}

export function Card_Back() {
  return (
    <div
      style={{
        width: 70,
        height: 100,
        backgroundColor: '#1e3a5f',
        border: '2px solid #0d1b2a',
        borderRadius: 8,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundImage: 'repeating-linear-gradient(45deg, transparent, transparent 5px, rgba(255,255,255,0.1) 5px, rgba(255,255,255,0.1) 10px)',
      }}
    >
      <div style={{ color: '#fff', fontSize: 24 }}>🀄</div>
    </div>
  )
}
