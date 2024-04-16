


export const post = async (url:string,data:any,headers={},serverHeaders:any=null)=>{
    const res = await fetch(`http://localhost:8080/${url}`, {
        method: "POST",
        headers: serverHeaders?serverHeaders:{
            ...headers
        
        },
        body: data,
        credentials:'include',
        cache:"no-store"

    })
    return await res.json()
}

export const get = async (url:string,headers={},serverHeaders:any=null)=>{
    const res = await fetch(`http://localhost:8080/${url}`, {
        method: "GET",
        headers: serverHeaders?serverHeaders:{
            "Content-Type": "application/json",
            ...headers
        
        },
        credentials:'include',
        cache:"no-cache"
    })
    return await res.json()
}