from flask import Flask, request, jsonify
import pika
import threading
from transcriber import transcribe_audio,stop_transcription,start_transcription
import json


app = Flask(__name__)

# RabbitMQ connection parameters
RABBITMQ_HOST = 'localhost'
RABBITMQ_PORT = 5672
RABBITMQ_USERNAME = 'guest'
RABBITMQ_PASSWORD = 'guest'

# List to store queue names
queues = []

def push_to_rabbitmq(message, streamId):
    connection = pika.BlockingConnection(pika.ConnectionParameters(
        host=RABBITMQ_HOST, port=RABBITMQ_PORT, credentials=pika.PlainCredentials(RABBITMQ_USERNAME, RABBITMQ_PASSWORD)))
    channel = connection.channel()
    channel.queue_declare(queue=streamId)
    channel.basic_publish(exchange='', routing_key=streamId, body=message)
    connection.close()

def consume_from_queue(queue_name):
    connection = pika.BlockingConnection(pika.ConnectionParameters(
        host=RABBITMQ_HOST, port=RABBITMQ_PORT, credentials=pika.PlainCredentials(RABBITMQ_USERNAME, RABBITMQ_PASSWORD)))
    channel = connection.channel()
    channel.queue_declare(queue=queue_name)

    def callback(ch, method, properties, body):
   
        print(f"Received message from queue {queue_name}: {body}")
        decoded_data = json.loads(body)
        print(decoded_data,"decoded_data")
        if 'stop' in decoded_data:
            stop_transcription(queue_name)
            return 

        transcribe_audio(decoded_data['message'], decoded_data['duration'], decoded_data['totalDuration'], decoded_data['segmentNumber'],queue_name)
        ch.basic_ack(delivery_tag=method.delivery_tag)  # Acknowledge the message

    channel.basic_consume(queue=queue_name, on_message_callback=callback, auto_ack=False)
    print(f"Consuming from queue {queue_name}")
    channel.start_consuming()

@app.route('/receive_text', methods=['POST'])
def receive_text():
    try:
        data = request.get_json()
        message = data['message']
        streamId = data['streamId']
        duration = data['duration']
        totalDuration = data['totalDuration']
        segmentNumber = data['segmentNumber']
        message_data = {
            'message': message,
            'duration': duration,
            'totalDuration': totalDuration,
            'segmentNumber': segmentNumber
        }
        message_json = json.dumps(message_data)
        print(f"Received text: {message_json}")



        # Add the queue to the list if it's not already there
        if streamId not in queues:
            queues.append(streamId)
            # Start a new thread to consume from the queue
            consumer_thread = threading.Thread(target=consume_from_queue, args=(streamId,))
            consumer_thread.start()

        # Push the received text to RabbitMQ queue
        push_to_rabbitmq(message_json, streamId)
        return jsonify({'success': True, 'message': 'Text received and pushed to RabbitMQ queue'})
    except Exception as e:
        print(f"Error: {str(e)}")
        return jsonify({'success': False, 'error': str(e)})
    

@app.route('/stop_transcription', methods=['POST'])
def stopTranscription():
    try:
        data = request.get_json()
       
        streamId = data['streamId']
        jsonData={
            "stop":True
        }
        jsonMessage=json.dumps(jsonData)

        push_to_rabbitmq(jsonMessage, streamId)
        # stop_transcription(streamId)
      
       

        
        return jsonify({'success': True, })
    except Exception as e:
        print(f"Error: {str(e)}")
        return jsonify({'success': False})


@app.route('/start_transcription', methods=['POST'])
def startTranscription():
    try:
        data = request.get_json()
       
        streamId = data['streamId']
        start_transcription(streamId)
      
       

        
        return jsonify({'success': True, })
    except Exception as e:
        print(f"Error: {str(e)}")
        return jsonify({'success': False})

if __name__ == '__main__':
    app.run(debug=True,port=5000,host='0.0.0.0')