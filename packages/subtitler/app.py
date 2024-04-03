from faster_whisper import WhisperModel

model_size = "tiny"

# Run on GPU with FP16
model = WhisperModel(model_size, device="cpu", compute_type="int8")

def generate_webvtt(cues):
    webvtt_content = "WEBVTT\n\n"
    
    for index, cue in enumerate(cues, start=1):
        start_time = cue["start"]
        end_time = cue["end"]
        text = cue["text"]
        
        webvtt_content += f"{index}\n"
        webvtt_content += f"{start_time} --> {end_time}\n"
        webvtt_content += f"{text}\n\n"
    
    return webvtt_content


# or run on GPU with INT8
# model = WhisperModel(model_size, device="cuda", compute_type="int8_float16")
# or run on CPU with INT8
# model = WhisperModel(model_size, device="cpu", compute_type="int8")

segments, info = model.transcribe("/home/anurag/s3mnt/stream1/stream10.ts", beam_size=5)

print("Detected language '%s' with probability %f" % (info.language, info.language_probability))

for segment in segments:
    
