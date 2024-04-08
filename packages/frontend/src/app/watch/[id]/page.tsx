"use client"
import React, { Ref } from 'react';

// This imports the functional component from the previous sample.
import VideoJS from '../_components/player';
import videojs from 'video.js';
import Player from 'video.js/dist/types/player';

export default function Watch({params}: {params: {id: string}}) {
  console.log(params.id,"params ");
  const playerRef = React.useRef<Player|null>(null);

  const videoJsOptions = {
    autoplay: true,
    controls: true,
    responsive: true,
    fluid: true,
    sources: [{
      src: `http://localhost:8080/hls/${params.id}/${params.id}.m3u8`,
    
    }]
  };

  const handlePlayerReady = (player:any) => {
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
    <>
      <VideoJS options={videoJsOptions} onReady={handlePlayerReady} />
    </>
  );
}