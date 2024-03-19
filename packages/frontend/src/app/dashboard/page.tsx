'use client';
import { useRef } from "react";
import VideoComponent from "./_components/videoComponent";

export default function DashBoard() {
    const mediaRecorderRef = useRef<MediaRecorder | null>(null);
    let socket: WebSocket | null = null;

    const startStreaming = () => {
        socket = new WebSocket('ws://localhost:8080/ws');

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
