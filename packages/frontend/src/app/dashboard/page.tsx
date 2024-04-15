'use client';
import { useEffect, useRef, useState } from "react";
import VideoComponent from "./_components/videoComponent";
import axios from 'axios';
import Card from "../_components/card";
import InfoComponent from "./_components/infoComponent";
import {
    RecoilRoot,
    atom,
    selector,
    useRecoilState,
    useRecoilValue,
} from 'recoil';
import Modal from "./_components/modal";
import { streamInfo } from "~/recoil/atom/streamInfo";
import { get, post } from "~/api";
import Chat from "../_components/chat";



export default function DashBoard() {
    const mediaRecorderRef = useRef<MediaRecorder | null>(null);
    let socket=useRef<WebSocket | null>(null);
    const [streamId, setStreamId] = useState<string>('');
    const [isStreaming,setIsStreaming] = useState(false)

    useEffect(()=>{

        return ()=>{
            stopStreaming()
        }
    },[])

    const startStreaming = (streamId: string) => {
        socket?.current?.close()
        socket.current = new WebSocket('ws://localhost:8080/ws', ["streamId", streamId]);

        socket.current.onopen = () => {
            console.log('WebSocket connection established');
        };

        mediaRecorderRef.current?.start(500);
        setIsStreaming(true)

        mediaRecorderRef.current?.addEventListener('dataavailable', (event) => {
            console.log(event.data);
            if (socket.current && socket.current.readyState === WebSocket.OPEN) {
                socket.current.send(event.data);
            }
        });
    };

    const stopStreaming = () => {
        console.log('stop streaming')
        if (socket.current) {
                console.log('closing socket')
            socket.current.close(      );
            socket.current = null; // Reset socket variable after closing
        }
    
        mediaRecorderRef.current?.stop();
        setIsStreaming(false)
    }

    const [streamInfoState, _] = useRecoilState(streamInfo)

    const stream = async () => {


        const data: {
            streamId: string
        } = await post('startStream', JSON.stringify(streamInfoState))

        setStreamId(data.streamId)
        startStreaming(data.streamId);



    }


    return (

        <div className="w-[100%] h-full p-5  flex ">
            <div className="flex h-full w-full ">
                <div className="w-full pr-5">
                    <Card classname={"flex flex-row w-[100%] h-fit"} >
                        <VideoComponent mediaRecorderRef={mediaRecorderRef} />
                        <InfoComponent />
                    </Card>


                    <div className="mt-3">
                        <button onClick={()=>{isStreaming?stopStreaming():stream()}} className={`btn ${isStreaming?"btn-error":"btn-primary"}`}>{isStreaming?"Stop Streaming":"Stream"}</button>
                    </div>
                </div>


                <div className="w-[500px] h-full " >
                    <Chat streamId={streamId} />
                </div>
            </div>
            <Modal />

        </div>
    );
}
