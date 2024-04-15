// "use client";
import Link from "next/link";
// import { authOptions } from "~/server/auth";
import Sidebar from "./_components/sidebar/sidebar";
import Home from "./_components/home";


export default async function HomePage({ params }: {
  params: {
    id: string,

  }
}) {
  // const data = await getServerSession(authOptions);
  // console.log(data,'data');  

  return (
    <main className="w-full h-full flex flex-row ">
      <Sidebar id={params.id} />
      <Home />
    </main>
  );
}
