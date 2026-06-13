import Nav from "../Nav"
import ImportFiles from "../ImportFiles"
import { useState, useEffect } from 'react'
import {getAll} from "../../services/data.js"
import {LineChart, Line, XAxis, YAxis, Legend, Tooltip, BarChart, Bar, CartesianGrid} from 'recharts'
import { RechartsDevtools } from '@recharts/devtools';
import "./chart.css"

const initialData=await getAll()
const startV=(initialData.success&&initialData.data.length>0 ? initialData.data[0].voivodeship:"")
const startY=(initialData.success&&initialData.data.length>0 ? initialData.data.reduce((m,x)=>Math.max(m,x.year),0):0)
const AllYears=(initialData.success&&initialData.data.length>0 ? initialData.data.reduce((m,x)=>{if(!m.includes(x.year)) m.push(x.year);return m;},[]) :[])
const AllRegions=(initialData.success&&initialData.data.length>0 ? initialData.data.reduce((m,x)=>{if(!m.includes(x.voivodeship)) m.push(x.voivodeship);return m;},[]) :[])

function ChartPage({username,refreshFunction}){

    const [data,setData]=useState(initialData)
    const [years,setYears]=useState(AllYears)
    const [regions,setRegions]=useState(AllRegions)
    const [chartType,setChartType]=useState('v')
    const [year,setYear]=useState(startY)
    const [region,setRegion]=useState(startV)
    const [visibleData,setVisibleData]=useState({type:"",data:[]})

    const refresh=async(e)=>{
        e.preventDefault()
        const res=await getAll()
        setData(res)
        setRegion(res.success&&initialData.data.length>0 ? res.data[0].voivodeship:"")
        setYear(res.success&&initialData.data.length>0 ? res.data.reduce((m,x)=>Math.max(m,x.year),0):0)
        setYears(res.success&&initialData.data.length>0 ? res.data.reduce((m,x)=>{if(!m.includes(x.year)) m.push(x.year);return m;},[]) :[])
        setRegions(res.success&&initialData.data.length>0 ? res.data.reduce((m,x)=>{if(!m.includes(x.voivodeship)) m.push(x.voivodeship);return m;},[]) :[])
    }
    const prepareChart=(e)=>{
        if(chartType=='v'){
            let vd=data.data.filter(w=>w.voivodeship==region)
            setVisibleData({type:'v',data:vd})
        }else{
            let vd=data.data.filter(w=>w.year==year)
            setVisibleData({type:'y',data:vd})
        }
    }
    
    return(
        <>
        <Nav username={username} refreshFunction={refreshFunction}/>
        {
            data.success && data.data.length>0 ?  
            <div><h2>Comparing average wages to the number of students</h2>
            <div className='col-8 offset-2'>Chart type:<select value={chartType} onChange={(e)=>{setChartType(e.target.value)}}>
                <option value='y'>All voivodeships in a specific year</option>
                <option value='v'>All years in a specific voivodeship</option>
                </select>
                {chartType=='v'?
                <>, voivodeship: <select value={region} onChange={(e)=>{setRegion(e.target.value)}}>{regions.map(m=><option key={m} value={m}>{m}</option>)}</select></>
                :
                <>, year: <select value={year} onChange={(e)=>{setYear(e.target.value)}}>{years.map(m=><option key={m} value={m}>{m}</option>)}</select></>}
                <button onClick={prepareChart}>Show</button></div>
                {visibleData.type=='v'&&visibleData.data.length>0?
                <div className='chart'><LineChart style={{width:'100%', aspectRatio: 1.618}} data={visibleData.data} responsive>
                    <Line dataKey='attending' name='Attending students per 10 000 residents' stroke='red' yAxisId="students"/>
                    <Line dataKey='graduating' name='Graduating students per 10 000 residents' stroke='blue' yAxisId="students"/>
                    <Line dataKey='wage' name='Average wage in PLN' stroke='green' yAxisId="money"/>
                    <XAxis dataKey='year'/>
                    <YAxis dataKey='attending' yAxisId='students' name='Number per 10 000' orientation="left" label={{value:"Number per 10 000", position:"insideLeft", angle: -90}}/>
                    <YAxis dataKey='wage' yAxisId='money' name='PLN' orientation="right" label={{value:"PLN", position:"insideRight", angle: -90}}/>
                    <Legend verticalAlign="top"/>
                    <Tooltip/>
                    <CartesianGrid/>
                    <RechartsDevtools/>
                    </LineChart></div>
                :(visibleData.type=='y'&&visibleData.data.length>0 ? <div className='chart'>
                <BarChart style={{width:'100%', aspectRatio: 1.618}} data={visibleData.data} responsive margin={{bottom: 160, left:3, right:3}}>
                    <Bar dataKey='attending' name='Attending students per 10 000 residents' fill='red' yAxisId="students"/>
                    <Bar dataKey='graduating' name='Graduating students per 10 000 residents' fill='blue' yAxisId="students"/>
                    <Bar dataKey='wage' name='Average wage in PLN' fill='green' yAxisId="money"/>
                    <XAxis type='category' dataKey='voivodeship' angle="70" tickMargin="100"/>
                    <YAxis dataKey='attending' yAxisId='students' name='Number per 10 000' orientation="left" label={{value:"Number per 10 000", position:"insideLeft", angle: -90}}/>
                    <YAxis dataKey='wage' yAxisId='money' name='PLN' orientation="right" label={{value:"PLN", position:"insideRight", angle: -90}}/>
                    <Legend verticalAlign="top"/>
                    <Tooltip/>
                </BarChart></div>:<div className="chart">No chart yet. Choose a chart type, then click the "show" button.</div>)}
                </div>:
            <div><h2>No data found!</h2>
            <p className="highlighted">{data.missingText}</p>
            <p>If this is because of an error, try logging in again or <a className="btn-link" onClick={refresh}>refreshing the page</a>.<br/>
            This could also be because there is no data in the database for some categories. You can load your own data in the <a href='/import'>Load Data tab</a>.<br/>
                The following errors occured while trying to load data:
            </p>
            <ul>{
            data.errors && data.errors.length>0 ? data.errors.map((x,i)=><li key={i}>{x}</li>) : <li>No errors.</li>
            }</ul>
            </div>
        }
        </>
    )
}

export default ChartPage