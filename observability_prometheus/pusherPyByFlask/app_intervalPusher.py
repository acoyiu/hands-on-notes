import os
import yaml
import time
from pprint import pprint
from pythonping import ping
from prometheus_client import CollectorRegistry, Gauge, push_to_gateway, generate_latest
from pathlib import Path
from flask import Flask
from threading import Thread, Event, Timer


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

def setInterval(interval):
    def decorator(function):
        def wrapper(*args, **kwargs):
            stopped = Event()
            def loop():                            # executed in another thread
                while not stopped.wait(interval):  # until stopped
                    function(*args, **kwargs)
            t = Thread(target=loop)
            t.daemon = True                        # stop if the program exits
            t.start()
            return stopped
        return wrapper
    return decorator

# =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

temporaryDown = False

def fetch():
    print('fetching...', flush=True)
    for ip, endpoint in ips:
        # the ping function
        Response = ping(
            ip,
            timeout=2,
            count=2,
            verbose=False
        )
        # pprint(dir(Response))
        print(f"Average in ms: {Response.rtt_avg_ms}", flush=True)
        for index, res in enumerate(Response):
            print(f"Try {index}: Pinged to {ip} as {endpoint}, {res.time_elapsed_ms}ms", flush=True)
        # set metrics one by one
        g.labels(ip, endpoint).set(2000 if temporaryDown else Response.rtt_avg_ms)
    # push to gateway
    push_to_gateway(pushgateway_location, job=job_name, registry=registry)
    print(f"\x1b[6;30;42mAll metrics sent!\x1b[0m", flush=True)

# =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

# function to ping all endpoints


@setInterval(15)
def to_fetch():
    fetch()


stop = to_fetch()

# =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

app = Flask(__name__)


@app.route("/stop_app")
def stop_app():
    print('stop_app...', flush=True)
    # stop the pusher
    stop.set()
    # set long run timer
    r = Timer(1.0, graceExiting)
    r.start()
    return "true"


@app.route("/manual_push_metrics")
def manual_push_metrics():
    print('manual_push_metrics...', flush=True)
    fetch()
    return "true"


@app.route("/temporary_down")
def toggle_down():
    print('toggling temporary_down...', flush=True)
    global temporaryDown
    temporaryDown = not temporaryDown
    return str(temporaryDown)

@app.route("/out_default")
def out_default():
    t = generate_latest()
    print(t, flush=True)
    return str(t)

@app.route("/out_registry")
def out_registry():
    t = generate_latest(registry)
    print(t, flush=True)
    return str(t)

@app.route("/metrics")
def metrics():
    t = generate_latest(registry)
    return str(t)


def graceExiting():
    time.sleep(5)
    print('exiting...', flush=True)
    os._exit(0)

# =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

fetch()