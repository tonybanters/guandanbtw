import { motion, AnimatePresence } from 'framer-motion'
import { Card as Card_Type, Rank } from '../game/types'
import { Card } from './Card'

interface Table_Props {
  cards: Card_Type[]
  level: Rank
  combo_type: string
}

export function Table({ cards, level, combo_type }: Table_Props) {
  const card_width = 70
  const overlap = 40

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: 180,
        padding: 20,
      }}
    >
      <div
        style={{
          display: 'flex',
          position: 'relative',
          width: cards.length > 0 ? card_width + (cards.length - 1) * overlap : 100,
          height: 100,
          justifyContent: 'center',
        }}
      >
        <AnimatePresence mode="wait">
          {cards.length > 0 ? (
            cards.map((card, index) => (
              <motion.div
                key={`table-${card.Id}`}
                initial={{ opacity: 0, scale: 0.5, y: 100 }}
                animate={{ opacity: 1, scale: 1, y: 0 }}
                exit={{ opacity: 0, scale: 0.5, y: -100 }}
                transition={{ delay: index * 0.05, type: 'spring', stiffness: 300 }}
                style={{
                  position: 'absolute',
                  left: index * overlap,
                  zIndex: index,
                }}
              >
                <Card
                  card={card}
                  level={level}
                  selected={false}
                  on_click={() => {}}
                />
              </motion.div>
            ))
          ) : (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 0.5 }}
              style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: '#666',
                fontSize: 14,
              }}
            >
              No cards played
            </motion.div>
          )}
        </AnimatePresence>
      </div>
      {combo_type && (
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          style={{
            marginTop: 12,
            padding: '4px 12px',
            backgroundColor: '#333',
            color: '#fff',
            borderRadius: 4,
            fontSize: 12,
            textTransform: 'uppercase',
          }}
        >
          {combo_type}
        </motion.div>
      )}
    </div>
  )
}
