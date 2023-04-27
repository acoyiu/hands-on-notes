from prometheus_client import start_http_server, Summary, Counter
import random
import time

# Create a metric to track time spent and requests made.
REQUEST_TIME = Summary('request_processing_seconds',
                       'Time spent processing request')


@REQUEST_TIME.time()  # Decorate function with metric.
def process_request(t):
    """A dummy function that takes some time."""
    time.sleep(t)


if __name__ == '__main__':
    # Start up the server to expose the metrics.
    start_http_server(8000)

# way to add Counter
c = Counter('my_failures', 'Description of counter')
c.inc()     # Increment by 1
c.inc(1.6)  # Increment by given value

print('metrics server is started')
print('get metrics by: curl localhost:8000')

# Generate some requests.
while True:
    process_request(random.random())
