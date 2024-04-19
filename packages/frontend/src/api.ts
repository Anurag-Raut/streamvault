


export const post = async (url: string, data: any, headers = {}, serverHeaders: any = null) => {
    try {


        const res = await fetch(`http://localhost:8080/${url}`, {
            method: "POST",
            headers: serverHeaders ? serverHeaders : {
                ...headers

            },
            body: data,
            credentials: 'include',
            cache: "no-store"

        })
        
        const resp= await res.json()
        if(resp.error){
            throw new Error(resp.error)
        }
        return resp
    }
    catch (error : any) {
        throw error
    }
}

export const get = async (url: string, headers = {}, serverHeaders: any = null) => {
    try{
        const res = await fetch(`http://localhost:8080/${url}`, {
            method: "GET",
            headers: serverHeaders ? serverHeaders : {
                "Content-Type": "application/json",
                ...headers
    
            },
            credentials: 'include',
            cache: "no-cache"
        })
        const resp= await res.json()
        if(resp.error){
            throw new Error(resp.error)
        }
        return resp;
    }
    catch(error:any){
        throw error
    }
   
}