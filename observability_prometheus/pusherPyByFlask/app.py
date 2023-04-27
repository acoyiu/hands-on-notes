import os
import yaml
import time
from pythonping import ping
from prometheus_client import CollectorRegistry, Gauge, push_to_gateway, generate_latest
from pathlib import Path
from flask import Flask
from threading import Timer


# =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-


ips = None
pushgateway_location = None

# Get the current working directory
with open(f"{Path(__file__).parent}/ep.yaml") as file:
    ymlData = yaml.load(file, Loader=yaml.FullLoader)
    ips = ymlData['ips']
    pushgateway_location = ymlData['pushgateway_location']

job_name = 'service_ping'
registry = CollectorRegistry()

# Declare prometheus gauge
g = Gauge(
    'ping_time_to_service',
    'the ms taken to reach services',
    labelnames=['ip', 'endpoint'],
    registry=registry
)


# =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-


temporaryDown = False

# function to ping all endpoints
def fetch():
    print('fetching...', flush=True)
    for ip, endpoint in ips:
        # the ping function
        Response = ping(ip, timeout=2, count=2, verbose=False)
        # set metrics one by one
        g.labels(ip, endpoint).set(2000 if temporaryDown else Response.rtt_avg_ms)
    print(f"\x1b[6;30;42mAll metrics updated!\x1b[0m", flush=True)


# =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-


app = Flask(__name__)


@app.route("/temporary_down")
def toggle_down():
    print('toggling temporary_down...', flush=True)
    global temporaryDown
    temporaryDown = not temporaryDown
    return str(temporaryDown)


@app.route("/metrics")
def metrics():
    fetch()
    t = generate_latest(registry)
    return t


@app.route("/stop_app")
def stop_app():
    print('stopping app...', flush=True)
    # set long run timer
    r = Timer(1.0, graceExiting)
    r.start()
    return "true"
    

def graceExiting():
    time.sleep(5)
    print('exiting...', flush=True)
    os._exit(0)