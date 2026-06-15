import {getJsonAsDataUrl, getXmlAsDataUrl} from "../../services/data.js"
import { useState } from 'react'

function ExportFiles(){

    const [error,setError]=useState("")
    const [ok,setOk]=useState(true)
    const [format,setFormat]=useState("application/xml,text/xml,.xml")
    const [category,setCategory]=useState("attending")
    const [files, setFiles]=useState([])

    const downloadClick = async (e)=>{
        if(format.includes("json")){
            setOk(true)
            setError("Downloading JSON...")
            let res=await getJsonAsDataUrl(category)
            setOk(res.success)
            if(ok){
                setError("File obtained, download it from the top of the list below.")
                setFiles([{href:res.href, name: category+'.json', timestamp: Date.now()}, ...files])
            }
            else setError(res.msg)
        }
        else if(format.includes("xml")){
            setOk(true)
            setError("Downloading XML...")
            let res=await getXmlAsDataUrl(category)
            setOk(res.success)
            if(ok){
                setError("File obtained, download it from the top of the list below.")
                setFiles([{href:res.href, name: category+'.xml', timestamp: Date.now()}, ...files])
            }
            else setError(res.msg)
        }
        else{
            setOk(false)
            setError("Unknown file extension")
        }
    }

    return(
        <div>
            <h2>Export from the database</h2>
            <div className="row"><label className='form-label offset-1 col-2' htmlFor='cat2'>Type of data:</label>
            <div className='col-8'><select id='cat2' className='form-input' height="3" value={category} onChange={(e)=>setCategory(e.target.value)}>
            <option value="attending">Attending students</option>
            <option value="graduating">Graduating students</option>
            <option value="averageWage">Average wage</option></select></div></div>
            <div className="row"><label className='form-label offset-1 col-2' htmlFor='format2'>Format:</label>
            <div className='col-8'><select id='format2' value={format} onChange={(e)=>setFormat(e.target.value)}><option value="application/xml,text/xml,.xml">XML</option><option value="application/json,.json">JSON</option></select></div></div>
            <button className='btn btn-primary mainbtn' onClick={downloadClick}>Get file</button>
            <p className={ok? "text-success":"text-danger"}>{error}</p>
            <h4>Exported files</h4>
            <ul>
                {files.length==0 ? <li>no files yet</li>:files.map(f=><li key={f.timestamp}><a download={f.name} href={f.href}>{f.name} at {new Date(f.timestamp).toISOString()}</a></li>)}
            </ul>
        </div>
    )
}

export default ExportFiles