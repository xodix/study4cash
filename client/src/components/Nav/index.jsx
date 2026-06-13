import {logout,deleteCurrent} from "../../services/auth.js"
import { useState } from 'react'
import "./Nav.css"

function Nav({username,refreshFunction}){
    const [error,setError]=useState("")

    const logoutClick=(e)=>{
        e.preventDefault()
        logout()
        refreshFunction()
        window.location="/"
    }
    const deleteClick=async (e)=>{
        e.preventDefault()
        if(window.confirm("Are you sure you want to delete your account?"))
        {
            const result=await deleteCurrent()
            console.log("result:",result)
            if(result.success) {
                refreshFunction();
                window.location="/"
            }
            else{
                setError("Error: "+result.msg.toString())
            }
        }
    }

    return(
        <>
        <nav className='navbar'>
            <h1 className="navbar-brand">Study4Cash</h1>
            <div className='col-4'><a href="/" className="nav-item btn btn-secondary">View charts</a><a href="/import" className="nav-item btn btn-secondary">Load data</a></div>
            <div className='col-4'><span>Current user: {username}</span><button className='nav-item btn btn-info' onClick={logoutClick}>Log out</button>
            <button className='nav-item btn btn-danger' onClick={deleteClick}>Delete account</button></div>
        </nav>
        <p className="text-danger">{error}</p>
        </>
    )
}

export default Nav