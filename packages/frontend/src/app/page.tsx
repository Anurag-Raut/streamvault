
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
    <main className="w-full h-[calc(100vh-82px)] min-h-[calc(100vh-82px)]  flex flex-row ">
      {/* <Sidebar id={params.id} /> */}
      <Home />
    </main>
  );
}