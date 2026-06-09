import Nav from "../Nav"
import ImportFiles from "../ImportFiles"
import ExportFiles from "../ExportFiles"

function ImportPage({username,refreshFunction}){

    return(
        <>
        <Nav username={username} refreshFunction={refreshFunction}/>
        <ImportFiles/>
        <ExportFiles/>
        </>
    )
}

export default ImportPage