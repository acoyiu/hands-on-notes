FROM python:alpine3.15

WORKDIR /app 

RUN cd /app

ADD ./app.py /app
ADD ./ep.yaml /app

EXPOSE 4000

RUN pip install pythonping prometheus_client Flask pyyaml

CMD flask run --host=0.0.0.0 --port=4000