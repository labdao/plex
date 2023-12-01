'use client'

import backendUrl from 'lib/backendUrl'
import React, { useEffect, useState } from 'react'
import { useSelector } from 'react-redux'

import {
  selectJobDetail,
} from '@/lib/redux'

const LogViewer = () => {
  const [logs, setLogs] = useState('')
  const job = useSelector(selectJobDetail)

  useEffect(() => {
    const BacalhauJobId = job.BacalhauJobID || window.location.href.split("/").pop()
    setLogs(`Connecting to stream with Bacalhau Job Id ${BacalhauJobId}`)

    let formattedBackendUrl = backendUrl().replace('http://', '').replace('https://', '')
    let wsProtocol = backendUrl().startsWith('https://') ? 'wss' : 'ws';

    console.log(formattedBackendUrl)
    const ws = new WebSocket(`${wsProtocol}://${formattedBackendUrl}/jobs/${BacalhauJobId}/logs`)

    ws.onopen = () => {
      console.log('connected')
    }

    ws.onmessage = (event) => {
      // Handle incoming message
      console.log(event.data);
      setLogs((prevLogs) => `${prevLogs}\n${event.data}`);
    };

    ws.onclose = () => {
      console.log('disconnected')
    }

    return () => {
      ws.close()
    }
  }, [job])

  return (
    <pre>
      {logs}
    </pre>
  )
}

export default LogViewer
