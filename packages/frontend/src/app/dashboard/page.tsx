'use client';
import { useRef, useState } from "react";
import VideoComponent from "./_components/videoComponent";
import axios from 'axios';

export default function DashBoard() {
    const mediaRecorderRef = useRef<MediaRecorder | null>(null);
    let socket: WebSocket | null = null;
    const [streamId, setStreamId] = useState<string>('');

    // const startStreaming = async () => {


    //     // const streamId = await axios.post('http://localhost:8080/startStream', {

    //     //     "title": "anurag"

    //     // }, {
    //     //     headers: {
    //     //         'Content-Type': 'application/json',

    //     //     }
    //     // }).then((response) => {
    //     //     console.log(response);
    //     //     setStreamId(response.data.streamId);
    //     //     console.log(response.data.streamId)


       
    //         const socket = new WebSocket('ws://localhost:8080/ws');
    //         socket.onopen = () => {
    //             console.log('WebSocket connection established');
    //         }

    //         mediaRecorderRef.current?.start(500);

    //         mediaRecorderRef.current?.addEventListener('dataavailable', (event) => {
    //             console.log(event.data);
    //             if (socket && socket.readyState === WebSocket.OPEN) {
    //                 socket.send(event.data);
    //             }
    //             // });





    //         }
    //         );





        // });

        // socket = new WebSocket('ws://localhost:8080/ws');

        // socket.onopen = () => {
        //     console.log('WebSocket connection established');
        // };
        // axios.post('http://localhost:8080/ws', {
        //     streamId: streamId
        // }).then((response) => {
        //     console.log(response);


        // }
        // );








    // };



    const startStreaming = () => {
        socket = new WebSocket('ws://localhost:8080/ws',["Bearer","token"]);

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

    return (
        <div className="w-[100%] h-screen p-5 ">
            <div>
                <VideoComponent mediaRecorderRef={mediaRecorderRef} />
            </div>
            <div>
                <button onClick={startStreaming} className="btn btn-primary">Stream</button>
            </div>
        </div>
    );
}
