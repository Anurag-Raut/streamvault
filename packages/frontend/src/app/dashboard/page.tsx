'use client';
import { useRef, useState } from "react";
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



export default function DashBoard() {
    const mediaRecorderRef = useRef<MediaRecorder | null>(null);
    let socket: WebSocket | null = null;
    const [streamId, setStreamId] = useState<string>('');

    const startStreaming = (streamId: string) => {
        socket = new WebSocket('ws://localhost:8080/ws', ["streamId", streamId]);

        socket.onopen = () => {
            console.log('WebSocket connection established');
        };

        mediaRecorderRef.current?.start(500);

        mediaRecorderRef.current?.addEventListener('dataavailable', (event) => {
            console.log(event.data);
            if (socket && socket.readyState === WebSocket.OPEN) {
                socket.send(event.data);
            }
        });
    };

    const [streamInfoState, _] = useRecoilState(streamInfo)

    const stream = async () => {
        // const csrfToken = await getCsrfToken();
        // console.log(csrfToken, "tokennn ")
        // const formData = new FormData();
        // formData.append("title", streamInfoState.title);
        // formData.append("description", streamInfoState.description);
        // formData.append("category", streamInfoState.category);
        // formData.append("thumbnail", streamInfoState.thumbnail);
        const response = await fetch("http://localhost:8080/startStream", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                // "Authorization": `Bearer ${csrfToken}` ,
           
            },
            body: JSON.stringify(streamInfoState),
            credentials: 'include'

        });
        const data: {
            streamId: string
        } = await response.json();

        startStreaming(data.streamId);



    }


    return (

        <div className="w-[100%] h-screen p-5  ">
            <Card classname={"flex flex-row w-[100%]"} >
                <VideoComponent mediaRecorderRef={mediaRecorderRef} />
                <InfoComponent />
            </Card>


            <div className="mt-3">
                <button onClick={stream} className="btn btn-primary">Stream</button>
            </div>
            <Modal />
        </div>
    );
}
