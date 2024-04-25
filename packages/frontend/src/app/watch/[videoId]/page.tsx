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
    title: string | null;
    description: string | null;
    category: string | null;
    likes: number | null;
    comments: number | null;
    createdAt: string | null;

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

    },],
    subtitles :[{
      kind: 'captions',
      srclang: 'en',
      label: 'English',
      src: '/test.vtt',
      mode:"showing",
      // default:true/

    }]

  };


  let captionOption = {
    kind: 'captions',
    srclang: 'en',
    label: 'English',
    src: `http://localhost:8080/hls/subtitle/${params.videoId}.vtt`,
    mode:"showing",
    default:true

  }

  const handlePlayerReady = (player: any) => {
    playerRef.current = player;

    setInterval(() => {
    player?.addRemoteTextTrack(captionOption);
    },2000)
    

    // You can handle player events here, for example:
    player.on('waiting', () => {
      videojs.log('player is waiting');
    });

    player.on('dispose', () => {
      videojs.log('player will dispose');
    });
  };

  return (

    <div className='w-[700px] min-h-[80%] bg-card rounded-xl overflow-hidden'>
      <VideoJS options={videoJsOptions}  onReady={handlePlayerReady} />
    </div>
  );
}