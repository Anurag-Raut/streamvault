from faster_whisper import WhisperModel
import os

model_size = "tiny"

# Run on GPU with FP16
model = WhisperModel(model_size, device="cpu", compute_type="int8")
file_handles = {}  

def seconds_to_formatted_time(seconds):
    hours = int(seconds // 3600)
    minutes = int((seconds % 3600) // 60)
    secs = int(seconds % 60)
    millis = int((seconds % 1) * 1000)
    return f"{hours:02d}:{minutes:02d}:{secs:02d}.{millis:03d}"


def start_transcription(streamId):
    try:
        # Directory where subtitle files are stored
        directory = '/home/anurag/s3mnt/subtitle'
        
        # Create the directory if it doesn't exist
        if not os.path.exists(directory):
            os.makedirs(directory)
        
        # Full file path for the subtitle file
        file_path = os.path.join(directory, f'{streamId}.vtt')
        
        # Check if the file handle already exists in the dictionary
        if streamId in file_handles:
            f = file_handles[streamId]
        else:
            # Open the file in append mode, creating it if it doesn't exist
            f = open(file_path, "a")
            file_handles[streamId] = f
        
        # Write the VTT header to the file
        f.write("WEBVTT\n\n")
        f.flush()
        
        print(f"Transcription started for stream {streamId}")
    except Exception as e:
        print(f"Error: {str(e)}")

def stop_transcription(streamId):
    try:
        if streamId in file_handles:
            f = file_handles[streamId]
        else:
            f = open(f'/home/anurag/s3mnt/subtitle/{streamId}.vtt', "a")
            file_handles[streamId] = f

        print(f"Transcription stopped for stream {streamId}")
        
        f.close()
    except Exception as e:
        print(f"Error: {str(e)}")

def generate_webvtt(start,end,text,index,streamId):
    try:
        start=seconds_to_formatted_time(start)
        end=seconds_to_formatted_time(end)
        webvtt_content = ""
        if streamId in file_handles:
            f = file_handles[streamId]
        else:
            f = open(f'/home/anurag/s3mnt/subtitle/{streamId}.vtt', "a")
            file_handles[streamId] = f
        
        with open('/home/anurag/s3mnt/subtitle/'+streamId+'.vtt', "a") as f:
            webvtt_content += f"{index}\n"
            webvtt_content += f"{start} --> {end}\n"
            webvtt_content += f"{text}\n\n"
            f.write(webvtt_content)
            f.flush()
    except Exception as e:
        print(f"Error: {str(e)}")


def transcribe_audio(audio_path="stream1/stream10.ts", duration=0, totalDuration=0, segmentNumber=0,streamId="stream1"):
    audio_path_str = audio_path

    segments, info = model.transcribe("/home/anurag/s3mnt/"+audio_path_str, beam_size=5)
    text=""
    for segment in segments:
        text+=segment.text+" "
    print(text)

   
 
    generate_webvtt(totalDuration, totalDuration+duration, text,segmentNumber,streamId)
    print(f"Transcription completed for segment {segmentNumber} of stream {streamId}")

#transcribe_audio()
