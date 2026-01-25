import { Card as Card_Type, Rank } from '../game/types'
import { Card } from './Card'

interface Table_Props {
  cards: Card_Type[]
  level: Rank
  combo_type: string
  last_play_seat?: number | null
}

export function Table({ cards, level, combo_type, last_play_seat }: Table_Props) {
  const card_width = 70
  const overlap = 40
  const show_highlight = last_play_seat !== null && last_play_seat !== undefined && cards.length > 0

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: 180,
        padding: 20,
        borderRadius: 12,
        transition: 'background-color 0.2s ease, box-shadow 0.2s ease',
        backgroundColor: show_highlight ? 'rgba(76, 175, 80, 0.15)' : 'transparent',
        boxShadow: show_highlight ? '0 0 20px rgba(76, 175, 80, 0.4)' : 'none',
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
        {cards.length > 0 ? (
          cards.map((card, index) => (
            <div
              key={`table-${card.Id}`}
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
            </div>
          ))
        ) : (
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: '#666',
              fontSize: 14,
              opacity: 0.5,
            }}
          >
            No cards played
          </div>
        )}
      </div>
      {combo_type && (
        <div
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
        </div>
      )}
    </div>
  )
}
