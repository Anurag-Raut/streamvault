import Card from "~/app/_components/card";




export default function Channel({ params }: { params: { id: string } }) {

    return (
        <div className="w-full h-full p-5">
            <h1 className="text-xl my-3">
                Channel Dashboard
            </h1>
            <div className="flex  w-full">

                <Card classname={'m-3 min-w-[300px] min-h-[400px]'} >
                    <div className="">
                        latest videos
                    </div>
                </Card>
                <Card classname={'m-3 min-w-[300px] min-h-[400px]'}>
                    <div className="">
                        <div className="mb-2">
                            <div className="text-md opacity-80">

                                Current Subscribers
                            </div>
                            <div className="text-xl mt-2 font-bold">
                                100
                            </div>
                        </div>
                        <div>
                            <div className="mb-2 text-lg">
                                Summary
                            </div>
                            <div className="mb-2 text-sm flex justify-between">
                                <h3>

                                    Views
                                </h3>
                                <h3>
                                    200
                                </h3>
                            </div>
                            <div className="mb-2 text-sm flex justify-between">
                                <h3>

                                    Watch time
                                </h3>
                                <h3>
                                    200
                                </h3>
                            </div>
                        </div>

                    </div>

                </Card>
            </div>
        </div>
    )

}