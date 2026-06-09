import { useState } from 'react'
import {register} from "../../services/auth.js"
import "../Auth.css"

const dateAsInputValue=(date)=>{
        return date.getFullYear()+"-"+date.getMonth().toString().padStart(2,'0')+"-"+date.getDate().toString().padStart(2,'0')
}

function Register({refreshFunction}) {
    const [errors, setErrors] = useState({})
    const [userData, setUserData] = useState({email:"",password:"",password2:"",name:"",surname:"",birthdate:dateAsInputValue(new Date(Date.now()))})

    const inputChanged = (e)=>{setUserData({...userData, [e.currentTarget.id]:e.currentTarget.value})}

    const onsubmit= async (e)=>{
        e.preventDefault()
        let err={}, ok=true
        if(userData.password!=userData.password2) {
            ok=false
            err.password="The two passwords are not identical"
        }
        const dateParsed=Date.parse(userData.birthdate)
        if(isNaN(dateParsed)){
            ok=false
            err.birthdate="Invalid date"
        }else userData.birthdate=new Date(dateParsed)
        if(ok){
            const {password2, ...userWithOnePwd}=userData
            const err= await register(userWithOnePwd)
            if(err===null) {
                refreshFunction()
                window.location="/"
            }
            setErrors(err ?? {})
        } else setErrors(err)
    }

    return (
        <form onSubmit={onsubmit} className="authMainDiv">
            <h2>Register</h2>
            <p className="text-danger offset-1">{errors.other ?? ""}</p>
            <div className="row"><label className="form-label offset-1 col-2" htmlFor="email">Email:</label>
            <div className="col-8"><input className='form-control' id="email" type="email" required value={userData.email} onChange={inputChanged}/></div></div>
            <p className="text-danger offset-1">{errors.email ?? ""}</p>
            <div className="row"><label className="form-label offset-1 col-2" htmlFor="name">First name:</label>
            <div className="col-8"><input className='form-control' id="name" type="text" minLength="3" maxLength="60" required value={userData.name} onChange={inputChanged}/></div></div>
            <p className="text-danger offset-1">{errors.name ?? ""}</p>
            <div className="row"><label className="form-label offset-1 col-2" htmlFor="surname">Surname:</label>
            <div className="col-8"><input className='form-control' id="surname" type="text" minLength="3" maxLength="60" required value={userData.surname} onChange={inputChanged}/></div></div>
            <p className="text-danger offset-1">{errors.surname ?? ""}</p>
            <div className="row"><label className="form-label offset-1 col-2" htmlFor="birthdate">Date of birth:</label>
            <div className="col-8"><input className='form-control' id="birthdate" type="date" required value={userData.birthdate} onChange={inputChanged}/></div></div>
            <p className="text-danger">{errors.birthdate ?? ""}</p>
            <div className="row"><label className="form-label offset-1 col-2" htmlFor="password">Password:</label>
            <div className="col-8"><input className='form-control' id="password" type="password" minLength="8" maxLength="32" required value={userData.password} onChange={inputChanged}/></div></div>
            <p className="text-danger offset-1">{errors.password ?? ""}</p>
            <div className="row"><label className="form-label offset-1 col-2" htmlFor="password2">Repeat password:</label>
            <div className="col-8"><input className='form-control' id="password2" type="password" minLength="8" maxLength="32" required value={userData.password2} onChange={inputChanged}/></div></div>
            <div className='offset-2 col-8'><input className="btn btn-primary mainbtn" type="submit" value="Register"/></div>
            <div className='offset-2 col-8'><a href="login">I already have an account</a></div>
        </form>
    )
}

export default Register