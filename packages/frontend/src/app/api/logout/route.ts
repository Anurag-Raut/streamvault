  import { NextResponse } from 'next/server';
  import axios from 'axios';
  import { get } from '~/api';
  import { cookies, headers } from 'next/headers';


  export async function GET(req: Request){
    
    try {
      // const response = await axios.get('${process.env.NEXT_PUBLIC_BACKEND_URL}/signOut')
      const response =await get('signOut',{
        Cookie:cookies().toString(),

      })
      console.log(response,"  response")
      // redirect("")
      //  return NextResponse.redirect(new URL("/",req.url))r
      return NextResponse.json({ message: 'Logged out successfully' })
    } catch (err:any) {
      console.error(err.toString())
      return NextResponse.json({ message: 'Internal server error' })
    }
  }