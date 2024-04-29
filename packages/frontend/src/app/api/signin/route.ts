import { NextResponse } from 'next/server';
import { cookies, headers } from 'next/headers';


export async function POST(req: Request){
  
  try {
    const {username,password}:{
      username:string,
      password:string
    }=await req.json()
    // const response = await axios.get('${process.env.NEXT_PUBLIC_BACKEND_URL}/signOut')
    const response =await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/signin`,{
        method: "POST",
        headers: {
            "Content-Type": "application/json",
       

        },
        credentials: 'include',
        cache: "no-store",
        body: JSON.stringify({
          username:username,
          password:password
        }),

      
    })

    console.log(response,"  response")
    if (response.ok) {
        const backendCookies = response.headers.get('set-cookie');
        if (backendCookies) {
            // You can now work with the cookies
            console.log('Cookies:', backendCookies);
            backendCookies.split(';').forEach((cookie) => {
                const [name, value] = cookie?.split('=');
                console.log('name:', name, 'value:', value);
                if(name && value)
                    cookies().set(name, value);

            }
            )
            
        }
        
        return NextResponse.json({ message: 'Logged out successfully' });
    } else {
        return NextResponse.json({ message: 'failed auth' })
    }
    // redirect("")
    //  return NextResponse.redirect(new URL("/",req.url))r
    return NextResponse.json({ message: 'Logged In successfully' })
  } catch (err:any) {
    console.error(err.toString())
    return NextResponse.json({ message: 'Internal server error' })
  }
}