import { useState } from 'react'
import {login} from "../../services/auth.js"
import "../Auth.css"

function Login({refreshFunction}) {
    const [errors, setErrors] = useState({})
    const [userData, setUserData] = useState({email:"",password:""})

    const inputChanged = (e)=>{setUserData({...userData, [e.currentTarget.id]:e.currentTarget.value})}

    const onsubmit= async (e)=>{
        e.preventDefault()
        const err= await login(userData)
        if(err===null) {
            refreshFunction()
            window.location="/"
        }
        setErrors(err ?? {})
    }

    return (
        <form onSubmit={onsubmit} className="authMainDiv">
            <h2>Log in</h2>
            <p className="text-danger offset-1">{errors.other ?? ""}</p>
            <div className="row"><label className="form-label offset-1 col-2" htmlFor="email">Email:</label>
            <div className="col-8"><input className='form-control' id="email" type="email" required value={userData.email} onChange={inputChanged}/></div></div>
            <p className="text-danger offset-1">{errors.email ?? ""}</p>
            <div className="row"><label className="form-label offset-1 col-2" htmlFor="password">Password:</label>
            <div className="col-8"><input className='form-control' id="password" type="password" minLength="8" maxLength="32" required value={userData.password} onChange={inputChanged}/></div></div>
            <p className="text-danger offset-1">{errors.password ?? ""}</p>
            <div className='offset-2 col-8'><input className="btn btn-primary mainbtn" type="submit" value="Log in"/></div>
            <div className='offset-2 col-8'><a href="register">Create an account</a></div>
        </form>
    )
}

export default Login