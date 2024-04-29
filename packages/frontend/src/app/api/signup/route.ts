import { NextResponse } from 'next/server';
import { cookies, headers } from 'next/headers';
import cookie from "cookie"

export async function POST(req: Request){
  
  try {
    const {username,password}:{
      username:string,
      password:string
    }=await req.json()
    // const response = await axios.get('${process.env.NEXT_PUBLIC_BACKEND_URL}/signOut')
    const response =await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/signup`,{
        method: "POST",
        headers: {
            "Content-Type": "application/json",

        },
        body: JSON.stringify({
          username:username,
          password:password
        }),
        credentials: 'include',
        cache: "no-store",
        
      
    })

    console.log(response,"  response")
    if (response.ok) {
      console.log(response.headers,"headerssss  ")
        const backendCookies = response.headers.get('set-cookie');
        if (backendCookies) {
            // You can now work with the cookies
            console.log('Cookies:', backendCookies);
            

            const backendCookie = cookie.parse(backendCookies);
            console.log(backendCookie,backendCookie.expires,"expires")
            for (const [name, value] of Object.entries(backendCookie)) {
              console.log('name:', name, 'value:', value);
              if(name && value && name==="jwt"){
                cookies().set(name, value,{
                  path: "/",
                  expires: new Date(backendCookie?.expires??new Date().getTime() + 1000*60*60*24*7),
                  sameSite: "none",
                  secure: true,
                  httpOnly: true
                });
              }
            }
          

            
            
        }
        
        return NextResponse.json({ message: 'Logged out successfully' });
    } else {
        
        return NextResponse.json({ error: 'failed auth' },{
          status: 401
        })
    }
    // redirect("")
    //  return NextResponse.redirect(new URL("/",req.url))r
    return NextResponse.json({ message: 'Logged In successfully' })
  } catch (err:any) {
    console.error(err.toString())
    return NextResponse.json({ error: 'Internal server error' },{
      status: 500
    })
  }
}