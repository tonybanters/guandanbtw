import { motion, AnimatePresence } from 'framer-motion'
import { Card as Card_Type, Rank } from '../game/types'
import { Card } from './Card'
import { use_is_mobile } from '../hooks/use_is_mobile'

interface Hand_Props {
  cards: Card_Type[]
  level: Rank
  selected_ids: Set<number>
  on_card_click: (id: number) => void
}

export function Hand({ cards, level, selected_ids, on_card_click }: Hand_Props) {
  const is_mobile = use_is_mobile()
  const card_width = is_mobile ? 56 : 70
  const card_height = is_mobile ? 80 : 100
  const overlap = is_mobile ? 28 : 35

  if (is_mobile) {
    return (
      <div
        style={{
          width: '100%',
          overflowX: 'auto',
          overflowY: 'hidden',
          WebkitOverflowScrolling: 'touch',
          padding: '8px 0',
          minHeight: card_height + 30,
        }}
      >
        <div
          style={{
            display: 'flex',
            position: 'relative',
            width: cards.length > 0 ? card_width + (cards.length - 1) * overlap + 16 : 0,
            height: card_height + 20,
            margin: '0 auto',
            paddingLeft: 8,
            paddingRight: 8,
          }}
        >
          <AnimatePresence>
            {cards.map((card, index) => (
              <motion.div
                key={card.Id}
                initial={{ opacity: 0, y: 30, scale: 0.8 }}
                animate={{ opacity: 1, y: 0, scale: 1 }}
                exit={{ opacity: 0, y: -30, scale: 0.8 }}
                transition={{ delay: index * 0.015 }}
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
                  size={is_mobile ? 'small' : 'normal'}
                />
              </motion.div>
            ))}
          </AnimatePresence>
        </div>
      </div>
    )
  }

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
          height: card_height,
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
