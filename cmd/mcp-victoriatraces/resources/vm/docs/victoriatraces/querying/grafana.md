---
weight: 4
title: Visualization in Grafana
disableToc: true
menu:
  docs:
    parent: "victoriatraces-querying"
    weight: 4
tags:
  - traces
aliases:
  - /victoriatraces/querying/grafana.html
---

[Grafana Jaeger Datasource](https://grafana.com/docs/grafana/latest/datasources/jaeger/) allows you to query and visualize VictoriaTraces data in Grafana.

![Visualization with Grafana](grafana-jaeger.webp)

Simply click "Add new data source" on Grafana, and then fill your VictoriaTraces URL to "Connection.URL". 

The URL format for VictoriaTraces single-node is:
```
http://<victoria-traces>:10428/select/jaeger
```

Finally, click "Save & Test" at the bottom to complete the process.