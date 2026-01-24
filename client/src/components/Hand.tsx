import { motion, AnimatePresence } from 'framer-motion'
import { Card as Card_Type, Rank } from '../game/types'
import { Card } from './Card'

interface Hand_Props {
  cards: Card_Type[]
  level: Rank
  selected_ids: Set<number>
  on_card_click: (id: number) => void
}

export function Hand({ cards, level, selected_ids, on_card_click }: Hand_Props) {
  const card_width = 70
  const overlap = 35

  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'center',
        padding: 20,
        minHeight: 140,
      }}
    >
      <div
        style={{
          display: 'flex',
          position: 'relative',
          width: cards.length > 0 ? card_width + (cards.length - 1) * overlap : 0,
          height: 100,
        }}
      >
        <AnimatePresence>
          {cards.map((card, index) => (
            <motion.div
              key={card.Id}
              initial={{ opacity: 0, y: 50, scale: 0.8 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: -50, scale: 0.8 }}
              transition={{ delay: index * 0.02 }}
              style={{
                position: 'absolute',
                left: index * overlap,
                zIndex: index,
              }}
            >
              <Card
                card={card}
                level={level}
                selected={selected_ids.has(card.Id)}
                on_click={() => on_card_click(card.Id)}
              />
            </motion.div>
          ))}
        </AnimatePresence>
      </div>
    </div>
  )
}
