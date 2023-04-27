const data = `
# TYPE some_metric counter
# HELP some_metric Just an example.
some_metric{label="val1"} 41
some_metric{label="val2"} 42
some_metric{label="val3"} 43
some_metric{label="val4"} 44
# TYPE some_another_metric gauge
# HELP some_another_metric Just another example.
some_another_metric 2398.283
`; // <- Remember to line change here, otherwise will error

fetch(
    'http://localhost:9091/metrics/job/job_name/instance/instance_name',
    {
        method: 'POST',
        headers: { 'Content-Type': 'text/plain' },
        body: data,
    },
);