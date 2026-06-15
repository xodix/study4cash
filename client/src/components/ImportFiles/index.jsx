import {sendFile} from "../../services/data.js"
import { useState } from 'react'

function ImportFiles(){

    const [error,setError]=useState("")
    const [ok,setOk]=useState(true)
    const [format,setFormat]=useState("application/xml,text/xml,.xml")
    const [file,setFile]=useState()
    const [category,setCategory]=useState("attending")

    const fileTypeChanged = (e)=>{
        e.preventDefault()
        setFormat(e.target.value)
        document.getElementById("file").value=null
        setFile()
    }

    const categoryChanged = (e)=>{
        e.preventDefault()
        setCategory(e.target.value)
        document.getElementById("file").value=null
        setFile()
    }

    const sendClick = async (e)=>{
        e.preventDefault()
        if(file===undefined){
            setOk(false)
            setError("No file")
            return
        }
        if(format.includes("xml")){
            setOk(true)
            setError("Sending...")
            let res=await sendFile(file,category,"XML")
            setOk(res.success)
            setError(res.msg)
        }
        else if(format.includes("json")){
            setOk(true)
            setError("Sending...")
            let res=await sendFile(file,category, "JSON")
            setOk(res.success)
            setError(res.msg)
        }
        else{
            setOk(false)
            setError("Unknown file extension")
        }
    }

    return(
        <div>
            <h2>Import your data</h2>
            <div className="row"><label className='form-label offset-1 col-2' htmlFor='cat'>Type of data:</label>
            <div className='col-8'><select id='cat' className='form-input' height="3" value={category} onChange={categoryChanged}>
            <option value="attending">Attending students</option>
            <option value="graduating">Graduating students</option>
            <option value="averageWage">Average wage</option></select></div></div>
            <div className="row"><label className='form-label offset-1 col-2' htmlFor='format'>Format:</label>
            <div className='col-8'><select id='format' value={format} onChange={fileTypeChanged}><option value="application/xml,text/xml,.xml">XML</option><option value="application/json,.json">JSON</option></select></div></div>
            <div className="row"><label className='form-label offset-1 col-2' htmlFor="file">File:</label>
            <div className='col-8'><input id="file" type='file' accept={format} onChange={(e)=>setFile(e.target.files[0])}/></div></div>
            <button className='btn btn-primary mainbtn' onClick={sendClick}>Send</button>
            <p className={ok? "text-success":"text-danger"}>{error}</p>
        </div>
    )
}

export default ImportFiles