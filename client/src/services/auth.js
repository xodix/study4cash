import {jwtDecode} from "jwt-decode"
import axios from "axios"
import loadConfig from "react-dynamic-env"

async function setTokenCookie(token, username){
    const decodedToken=jwtDecode(token)
    const expires=new Date(decodedToken.exp*1000).toUTCString()
    if(username===undefined){
        try{
            const env=await loadConfig()
            const {data: d} = await axios({url:env.API_URL+"/user",method:"get",headers:{'Authorization':"Bearer "+token}}) //Wymienic na env przed wrzuceniem do kontenera
            username=d.name+" "+d.surname
            return null
        }catch(e){
            console.log(e.response)
            username="Unknown"
        }
    }
    document.cookie = `token=${token}; expires=${expires}`
    document.cookie = `username=${username}; expires=${expires}`
}

function getUserFromCookie(){
    const cookies=decodeURIComponent(document.cookie).split(";")
    let user={}, findtoken=true, findname=true
    for(let i=0;(findtoken||findname)&&i<cookies.length;i++){
        if(cookies[i].includes("token=")) {
            user.token=cookies[i].substring(cookies[i].indexOf("token=")+6,cookies[i].length)
            findtoken=false
        }else if(cookies[i].includes("username=")) {
            user.username=cookies[i].substring(cookies[i].indexOf("username=")+9,cookies[i].length)
            findname=false
        }
    }
    if(!findtoken) return user
    else return undefined
}

function decodeError(plaintext){
    let errors={}
    const text=plaintext.split("\n")
    for(const line of text){
        if(line.includes("RegisterRequest.")){
            const startOfFieldName=line.indexOf("RegisterRequest.")+16
            const endOfFieldName=line.indexOf("'",startOfFieldName)
            errors[line.substring(startOfFieldName,endOfFieldName).toLowerCase()]=line.substring(line.indexOf("Error:")+6)
        }
        else if(line.includes("parsing time")) errors.birthdate = "This is not a valid date"
        else errors.other = (errors.other ?? "")+line+"\n"
    } 
    return errors
}

async function register(data){
    try{
        const env=await loadConfig()
        const {data: d} = await axios.post(env.API_URL+"/user/register",data) //Wymienic na env przed wrzuceniem do kontenera
        await setTokenCookie(d.token,data.name+" "+data.surname)
        return null
    }catch(e){
        if(e.response&&e.response.status==400){
            return decodeError(e.response.data)
        }else return {other: "Unknown error."}
    }
}

async function login(data){
    try{
        const env=await loadConfig()
        const {data: d} = await axios.post(env.API_URL+"/user/login",data) //Wymienic na env przed wrzuceniem do kontenera
        await setTokenCookie(d.token)
        return null
    }catch(e){
        if(e.response){
            if(e.response.status==400)
                return decodeError(e.response.data)
            if(e.response.status==401)
                return {password: "Incorrect password."}
        }else return {other: "Unknown error."}
    }
}

async function deleteCurrent(){
    console.log("entered")
    const user=getUserFromCookie();
    if(user===undefined) return {success: false, msg: "You are not currently logged in."}
    try{
        const env=await loadConfig()
        const d = await axios({url:env.API_URL+"/user",method:"delete",headers:{'Authorization':"Bearer "+user.token}}) //Wymienic na env przed wrzuceniem do kontenera
        logout()
        return {success: d.status==200, msg: e.response.data.message}
    }catch(e){
        if(e.response&&e.response.status>400&&e.response.status<=404) return {success: false, msg: e.response.data.message}
    }
}

function logout(){
    const past=new Date(Date.now()-1000).toUTCString()
    document.cookie = `token=; expires=${past}`
    document.cookie = `username=; expires=${past}`
}

export {register, login, logout, getUserFromCookie, deleteCurrent}
