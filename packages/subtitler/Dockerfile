FROM python:3.10.12

WORKDIR /

COPY requirements.txt .

RUN pip3 install --no-cache-dir -r requirements.txt 

COPY . .

EXPOSE 5000

CMD ["python3", "server.py"]
