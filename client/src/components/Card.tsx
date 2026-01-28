import { motion } from 'framer-motion'
import { Card as Card_Type, get_suit_symbol, get_rank_symbol, is_red_suit, is_wild, Rank, Rank_Red_Joker, Suit_Joker } from '../game/types'

type Card_Size = 'small' | 'normal'

interface Card_Props {
  card: Card_Type
  level: Rank
  selected: boolean
  on_click: () => void
  size?: Card_Size
}

const SIZE_CONFIG = {
  small: { width: 56, height: 80, rank_font: 14, suit_font: 16, center_font: 20, corner_rank: 12, corner_suit: 10 },
  normal: { width: 70, height: 100, rank_font: 16, suit_font: 18, center_font: 24, corner_rank: 14, corner_suit: 12 },
}

export function Card({ card, level, selected, on_click, size = 'normal' }: Card_Props) {
  const is_joker = card.Suit === Suit_Joker
  const is_red = is_joker ? card.Rank === Rank_Red_Joker : is_red_suit(card.Suit)
  const is_wild_card = is_wild(card, level)
  const cfg = SIZE_CONFIG[size]

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
        width: cfg.width,
        height: cfg.height,
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
        <span style={{ fontSize: cfg.rank_font, fontWeight: 'bold' }}>{is_joker ? (card.Rank === Rank_Red_Joker ? 'R' : 'B') : get_rank_symbol(card.Rank)}</span>
        <span style={{ fontSize: cfg.suit_font }}>{is_joker ? '🃏' : get_suit_symbol(card.Suit)}</span>
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
        <span style={{ fontSize: cfg.corner_rank }}>{is_joker ? (card.Rank === Rank_Red_Joker ? 'R' : 'B') : get_rank_symbol(card.Rank)}</span>
        <span style={{ fontSize: cfg.corner_suit }}>{is_joker ? '🃏' : get_suit_symbol(card.Suit)}</span>
      </div>
      <div style={{
        position: 'absolute',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
        fontSize: is_joker ? cfg.center_font + 4 : cfg.center_font,
      }}>
        {is_joker ? '🃏' : get_suit_symbol(card.Suit)}
      </div>
    </motion.div>
  )
}

interface Card_Back_Props {
  size?: Card_Size
}

export function Card_Back({ size = 'normal' }: Card_Back_Props) {
  const cfg = SIZE_CONFIG[size]

  return (
    <div
      style={{
        width: cfg.width,
        height: cfg.height,
        backgroundColor: '#1e3a5f',
        border: '2px solid #0d1b2a',
        borderRadius: 8,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundImage: 'repeating-linear-gradient(45deg, transparent, transparent 5px, rgba(255,255,255,0.1) 5px, rgba(255,255,255,0.1) 10px)',
      }}
    >
      <div style={{ color: '#fff', fontSize: size === 'small' ? 20 : 24 }}>🀄</div>
    </div>
  )
}
