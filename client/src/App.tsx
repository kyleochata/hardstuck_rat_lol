import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  const [count, setCount] = useState(0)
  const fetchData = async () => {
    try {
      const response = await fetch(`api/lol/kyleochata/NA1`)
      if (response.ok) {
        const data = await response.json()
        console.log(data)
      } else {
        console.error('Fetch error:', response.status, response.statusText)
      }
    } catch (e) {
      console.error('Network error:', e)
    }
  }

  const fetchContext = async () => {
    const r = await fetch(`api/lol/kyleochata/NA1/matches`)
    if (r.ok) {
      console.log(r.json())
    } else {
      console.log('error2', r.status)
    }
  }

  const getboth = async () => {
    const sum = await fetchData()
    const fetch = await fetchContext()

    return
  }
  getboth()
  return (
    <>
      <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  )
}

export default App
