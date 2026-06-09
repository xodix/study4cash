import loadConfig from "react-dynamic-env"
import axios from "axios"
import {getUserFromCookie} from "./auth.js"
import XMLBuilder from "fast-xml-builder"

//glowny element root, w nim wiele <attending><Kod>0200000</Kod><Nazwa>DOLNOŚLĄSKIE</Nazwa><Y2002>523</Y2002></attending> lub jakakolwiek inna kategoria

async function sendXml(file, category){
    var formData = new FormData();
    formData.append("file", file);
    const env=await loadConfig()
    try{
        const user=getUserFromCookie()
        const res=await axios.put(env.API_URL+'/data/'+category+"/XML", formData, {headers: {'Content-Type': 'multipart/form-data','Authorization':"Bearer "+user.token}})
        return {success: true, msg: "Data added successfully."}
    }catch(e){
        console.log("error",e)
        if(e.response){
            if(e.response.status==400||e.response.status==401||e.response.status==404||e.response.status==500)
                return {success: false, msg: e.res}
            else return {success: false, msg: "Unknown error."}
        }else return {success: false, msg: "Unknown error."}
    }
    
}
async function sendJson(file,category){
    let fileContents=await new Promise((resolve)=>{
        try {
            var reader= new FileReader()
            reader.onload = (e) => resolve(reader.result);
            reader.readAsText(file)
        } catch (error) {}
    })
    if(fileContents==undefined) return {success: false, msg: "There was an error while reading the file."}
    try {
        const objs=JSON.parse(fileContents)
        const xml=ExportObjectToXml(objs,category)
        const xmlFile=new Blob([xml],{type:"text/xml"})
        return sendXml(xmlFile,category)
    } catch (error) {
        return {success: false, msg: "There was an error while parsing the file. Are you sure the JSON is correct?"}
    }
}

async function getOneCategory(category){
    const env=await loadConfig()
    try{
        const user=getUserFromCookie()
        const {data: res}=await axios.get(env.API_URL+'/data/'+category, {headers: {'Authorization':"Bearer "+user.token}})
        return {success: true, data: res}
    }catch(e){
        if(e.response){
            if(e.response.status==400||e.response.status==401||e.response.status==404||e.response.status==500)
                return {success: false, msg: e.res}
            else return {success: false, msg: "Unknown error."}
        }else return {success: false, msg: "Unknown error."}
    }
}

async function getExports(category){
    const fieldName=(category=="attending"?"StudentsAttending":(category=="graduating"?"StudentsGraduating":"AverageWage"))
    const downloaded=await getOneCategory(category)
    if(downloaded.success==false) return downloaded
    var result=[];
    for(let i=0;i<downloaded.data.length;i++){
        let x=result.findIndex(w=>(w.Nazwa==downloaded.data[i].Voivodeship))
        if(x===-1){
            result.push({Nazwa:downloaded.data[i].Voivodeship})
        }
        else{
            result[x]["Y"+downloaded.data[i].Year] = downloaded.data[i][fieldName].toString().replace('.',',')
        }
    }
    return {success:true, data:result}
}

async function getJsonAsDataUrl(category){
    var objs=await getExports(category)
    if(objs.success==false) return objs
    return{success:true, href:"data:application/json;charset=utf-8,"+encodeURI(JSON.stringify(objs.data))}
}

function ExportObjectToXml(data, category){
    const fieldName=category=="graduating"?"graduated":(category=="averageWage"?"averagewage":category)
    const builder=new XMLBuilder({arrayNodeName:fieldName, format:true})
    const txt='<?xml version="1.0" encoding="utf-8"?>\n<root>'+builder.build(data)+"</root>"
    return txt
}

async function getXmlAsDataUrl(category){
    var objs=await getExports(category)
    if(objs.success==false) return objs
    const txt=ExportObjectToXml(objs.data, category)
    return{success:true, href:"data:application/xml;charset=utf-8,"+encodeURI(txt)}
}

async function getAll(){
    const env=await loadConfig()
    try{
        const user=getUserFromCookie()
        const config={headers: {'Authorization':"Bearer "+user.token}}
        const {data: att}=await axios.get(env.API_URL+'/data/attending',config)
        const {data: grad}=await axios.get(env.API_URL+'/data/graduating', config)
        const {data: wage}=await axios.get(env.API_URL+'/data/averageWage', config)
        let result=wage.map(w=>{
            let res={voivodeship:w.Voivodeship, year:w.Year, wage: w.AverageWage}
            let foundGrad=grad.find(a=>(a.Year==w.Year&&a.Voivodeship==w.Voivodeship)), foundAtt=att.find(a=>(a.Year==w.Year&&a.Voivodeship==w.Voivodeship))
            res.attending= foundAtt==undefined? undefined:foundAtt.StudentsAttending
            res.graduating= foundGrad==undefined? undefined:foundGrad.StudentsGraduating
            return res

        })
        return {success: true, data: result.filter(w=>(w.attending!=undefined && w.graduating!=undefined))}
    }catch(e){
        console.log("error",e)
        if(e.response){
            if(e.response.status==400||e.response.status==401||e.response.status==404||e.response.status==500)
                return {success: false, msg: e.res}
            else return {success: false, msg: "Unknown error."}
        }else return {success: false, msg: "Unknown error."}
    }
}

export {sendXml, sendJson, getAll, getJsonAsDataUrl, getXmlAsDataUrl}