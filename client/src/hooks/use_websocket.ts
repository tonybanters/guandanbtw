import { useCallback, useEffect, useRef, useState } from 'react'
import { Message } from '../game/types'

type Message_Handler = (msg: Message) => void

export function use_websocket(url: string) {
  const [connected, set_connected] = useState(false)
  const ws_ref = useRef<WebSocket | null>(null)
  const handlers_ref = useRef<Map<string, Message_Handler>>(new Map())

  useEffect(() => {
    const ws = new WebSocket(url)
    ws_ref.current = ws

    ws.onopen = () => {
      set_connected(true)
    }

    ws.onclose = () => {
      set_connected(false)
    }

    ws.onmessage = (event) => {
      const msg: Message = JSON.parse(event.data)
      const handler = handlers_ref.current.get(msg.type)
      if (handler) {
        handler(msg)
      }
    }

    return () => {
      ws.close()
    }
  }, [url])

  const send = useCallback((msg: Message) => {
    if (ws_ref.current && ws_ref.current.readyState === WebSocket.OPEN) {
      ws_ref.current.send(JSON.stringify(msg))
    }
  }, [])

  const on = useCallback((type: string, handler: Message_Handler) => {
    handlers_ref.current.set(type, handler)
    return () => {
      handlers_ref.current.delete(type)
    }
  }, [])

  return { connected, send, on }
}
