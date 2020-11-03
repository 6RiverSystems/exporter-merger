# exporter-merger

[![Build Status](https://travis-ci.org/rebuy-de/exporter-merger.svg?branch=master)](https://travis-ci.org/rebuy-de/exporter-merger)
[![license](https://img.shields.io/github/license/rebuy-de/exporter-merger.svg)]()
[![GitHub release](https://img.shields.io/github/release/rebuy-de/exporter-merger.svg)]()

Merges Prometheus metrics from multiple sources.

> **Development Status** *exporter-merger* is in an early development phase.
> Expect incompatible changes and abandoment at any time.

## But Why?!

> [prometheus/prometheus#3756](https://github.com/prometheus/prometheus/issues/3756)

## Usage

*exporter-merger* needs a configuration file. Currently, nothing but URLs are accepted:

```yaml
exporters:
- url: http://localhost:9100/metrics
  addLabels:
  - name: "job"
    value: "someDownstreamThing"
  - name: "anotherLabel"
    value: "anotherLabelValue!"
- url: http://localhost:9101/metrics
```

To start the exporter:

```
exporter-merger --config-path merger.yaml --listen-ip 10.100.1.2 --listen-port 8080
```

### Environment variables

Alternatively configuration can be passed via environment variables, here is relevant part of `exporter-merger -h` output:
```
      --listen-port int      Listen port for the HTTP server. (ENV:MERGER_PORT) (default 8080)
```

## Kubernetes

The exporter-merger is supposed to run as a sidecar. Here is an example config with [nginx-exporter](https://github.com/rebuy-de/nginx-exporter):

```yaml
apiVersion: apps/v1
kind: Deployment

metadata:
  name: my-nginx
  labels:
    app: my-nginx

spec:
  selector:
    matchLabels:
      app: my-nginx

  template:
    metadata:
      name: my-nginx
      labels:
        app: my-nginx
      annotations:
        prometheus.io/should-be-scraped: "true"
        prometheus.io/scrape-port: "8080"
    spec:
      containers:
      - name: "nginx"
        image: "my-nginx" # nginx image with modified config file

        volumeMounts:
        - name: mtail
          mountPath: /var/log/nginx/mtail

      - name: nginx-exporter
        image: quay.io/rebuy/nginx-exporter:v1.1.0
        ports:
        - containerPort: 9397
        env:
        - name: NGINX_ACCESS_LOGS
          value: /var/log/nginx/mtail/access.log
        - name: NGINX_STATUS_URI
          value: http://localhost:8888/nginx_status
        volumeMounts:
        - name: mtail
          mountPath: /var/log/nginx/mtail

      - name: exporter-merger
        image: gcr.io/plasma-column-128721/exporter-merger:latest
        volumeMounts:
          - mountPath: /etc/exporter-merger/
            name: example-exporter-merger
            readOnly: true
     volumes:
       - name: example-exporter-merger
         configMap:
           name: example-exporter-merger

```

## Planned Features

* Allow transforming of metrics from backend exporters.
  * eg add a prefix to the metric names
* Allow dynamic adding of exporters.
