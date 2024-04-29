import { NextResponse } from 'next/server';
import axios from 'axios';
import { get } from '~/api';
import { cookies, headers } from 'next/headers';
import cookie from 'cookie';

export async function POST(req: Request){
  
  try {
    const code:string = await req.json()
    console.log(code,"  code")
    // const response = await axios.get('${process.env.NEXT_PUBLIC_BACKEND_URL}/signOut')
    const response =await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/signinWithGoogle`,{
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            ...headers

        },
        body: JSON.stringify(code),
        credentials: 'include',
        cache: "no-store"
      
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
        console.log(response,"  response")
        return NextResponse.json({ error: 'failed auth' },{
          status: 401
        })
    }
    // redirect("")
    //  return NextResponse.redirect(new URL("/",req.url))r
    return NextResponse.json({ message: 'Logged out successfully' })
  } catch (err:any) {
    console.error(err.toString())
    return NextResponse.json({ error: 'Internal server error' },{
      status: 500
    })
  }
}