'use client'
import React, { useRef, useEffect } from 'react';

const CameraField = ({ mediaRecorderRef }: any) => {
  const videoRef = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    const constraints: MediaStreamConstraints = { video: true, audio: true };

    const enableCamera = async () => {
      try {
        const stream = await navigator.mediaDevices.getUserMedia(constraints);
        if (videoRef.current) {
          videoRef.current.srcObject = stream;
          mediaRecorderRef.current = new MediaRecorder(stream, {
            videoBitsPerSecond: 3000000,
            audioBitsPerSecond: 64000,
          });
        }
      } catch (err) {
        console.error('Error accessing the camera:', err);
      }
    };

    enableCamera();

    return () => {
      if (videoRef.current) {
        const stream = videoRef.current.srcObject as MediaStream;
        if (stream) {
          const tracks = stream.getTracks();
          tracks.forEach(track => track.stop());
        }
      }
    };
  }, []);

  return (
    <div
   
      className='w-[550px] h-[225px] bg-black flex justify-center items-center  '
    >
      <video
        ref={videoRef}
        width={300}
        height={225} // 16:9 aspect ratio (400/225 = 16/9)
        autoPlay
        playsInline
        style={{
          objectFit:"cover"
        }}
  
      ></video>
    </div>
  );
};

export default CameraField;
