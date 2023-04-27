cat <<EOF | curl --data-binary @- http://localhost:9091/metrics/job/some_job/instance/some_instance
# TYPE some_metric counter
some_metric{label="val1"} 42
# TYPE some_another_metric gauge
# HELP some_another_metric Just an example.
some_another_metric 2398.283
EOF