import { useState } from 'react'
import './App.css'
import Login from "./components/Login/index.jsx"
import Register from "./components/Register"
import Nav from "./components/Nav"
import ImportPage from "./components/ImportPage"
import ChartPage from "./components/ChartPage"
import {getUserFromCookie} from "./services/auth.js"
import {Routes, Route, BrowserRouter, Navigate} from "react-router-dom"

function App() {
  const [user, setUser] = useState(getUserFromCookie)
  const refreshUser = () => {setUser(getUserFromCookie())}

  return (
    <BrowserRouter>
    {
      user ?
      <Routes>
          <Route path="/import" element={<ImportPage username={user.username} refreshFunction={refreshUser}/>}/>
          <Route path="/" element={<ChartPage username={user.username} refreshFunction={refreshUser}/>}/>
          <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>:
      <Routes>
          <Route path="/login" exact element={<Login refreshFunction={refreshUser}/>}/>
          <Route path="/register" exact element={<Register refreshFunction={refreshUser}/>}/>
          <Route path="*" element={<Navigate to="/login" replace />}
          />
      </Routes>
    }
    </BrowserRouter>
  )
}

export default App
