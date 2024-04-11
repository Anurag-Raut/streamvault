"use client"
import React, { Ref, useEffect } from 'react';

// This imports the functional component from the previous sample.
import VideoJS from '../_components/player';
import videojs from 'video.js';
import Player from 'video.js/dist/types/player';

export default function Watch({ params }: { params: { videoId: string } }) {
  console.log(params.videoId, "params ");
  const playerRef = React.useRef<Player | null>(null);

  const [data, setData] = React.useState<{
    title: string|null;
    description: string|null;
    category: string|null;
    likes: number|null;
    comments: number|null;
    createdAt: string|null;
  
  }>({
    title: null,
    description: null,
    category: null,
    likes: null,
    comments: null,
    createdAt: null,
  });

  const videoJsOptions = {
    autoplay: true,
    controls: true,
    responsive: true,
    aspectRatio: '16:9',


    fluid: true,
    sources: [{
      src: `http://localhost:8080/hls/${params.videoId}/${params.videoId}.m3u8`,

    }]
  };
  // useEffect(() => {
  //   async function fetchData() {
  //       const response=await fetch(`http://localhost:8080/getVideoData`,{
  //           method:"POST",
  //           headers:{
  //               "Content-Type":"application/json"
  //           },
  //           body:JSON.stringify(params.videoId),
  //           credentials:'include'
          
  //       })
  //       const data=await response.json()
  //       console.log(data)
  //       setData(data)

  //   }
  //   fetchData()
  // }, []);

  const handlePlayerReady = (player: any) => {
    playerRef.current = player;

    // You can handle player events here, for example:
    player.on('waiting', () => {
      videojs.log('player is waiting');
    });

    player.on('dispose', () => {
      videojs.log('player will dispose');
    });
  };

  return (
    // <div className='w-full h-full p-9 '>
    //   <div className='w-full   h-full' >
        <div className='w-[55%] rounded-xl overflow-hidden'>
          <VideoJS options={videoJsOptions} onReady={handlePlayerReady} />
        </div>
      //   <div className='mt-3'>
      //       <div className='text-lg '>{data.title}</div>

      //   </div>
      // </div>


    // </div>
  );
}