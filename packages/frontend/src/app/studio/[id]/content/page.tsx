import Card from "~/app/_components/card";


export default async function Content({ params }: { params: { id: string } }) {
    const wa= await fetch(`http://localhost:8080/getUserId`,{
        method:"GET",
        headers:{
            "Content-Type":"application/json"
        },
        credentials:"include"
    },)
    const d=await wa.text()
    console.log(d,'name ')
    const data=[1,2,3,4]
    return (
        <div className="w-full h-full p-5">
            <h1 className="text-xl my-3">
                Channel Content
            </h1>
            <Card>
                <div className="overflow-x-auto">
                    <table className="table">
                        {/* head */}
                        <thead>
                            <tr>
                                <th>
                                    <label>
                                        <input type="checkbox" className="checkbox" />
                                    </label>
                                </th>
                                <th>Video</th>
                                <th>Date</th>
                                <th>Likes </th>
                                <th>Comments</th>

                                <th></th>
                            </tr>
                        </thead>
                        <tbody>
                            {
                                data.map((item,index)=>{
                                    return (
                                        <tr className="hover:opacity-50 hover:bg-black">
                                        <th>
                                            <label>
                                                <input type="checkbox" className="checkbox" />
                                            </label>
                                        </th>
                                        <td>
                                            <div className="flex items-center gap-3">
                                                <div className="avatar">
                                                    <div className="mask mask-squircle w-12 h-12">
                                                        <img src="/tailwind-css-component-profile-2@56w.png" alt="Avatar Tailwind CSS Component" />
                                                    </div>
                                                </div>
                                                <div>
                                                    <div className="font-bold">Hart Hagerty</div>
                                                    <div className="text-sm opacity-50">United States</div>
                                                </div>
                                            </div>
                                        </td>
                                        <td>
                                            Zemlak, Daniel and Leannon
                                            <br />
                                            <span className="badge badge-ghost badge-sm">Desktop Support Technician</span>
                                        </td>
                                        <td>Purple</td>
                                        <th>
                                            <button className="btn btn-ghost btn-xs">details</button>
                                        </th>
                                    </tr>
                                    )
                                })
                                
                            }
                           
                        </tbody>
                        {/* foot */}
                        <tfoot>
                            <tr>
                                <th></th>
                                <th>Name</th>
                                <th>Job</th>
                                <th>Favorite Color</th>
                                <th></th>
                            </tr>
                        </tfoot>

                    </table>
                </div>
            </Card>
        </div>
    )

}