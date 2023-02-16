local grafana = import 'grafonnet/grafana.libsonnet';

local singlestatHeight = 100;
local singlestatGuageHeight = 150;

grafana.dashboard.new(
  'Eth Balance Proxy Dashboard',
  description='Metrics about the proxy and its upstreams',
  tags=['kubernetes'],
  time_from='now-1h',
)
.addTemplate(
  grafana.template.datasource(
    'datasource',
    'prometheus',
    '',
  )
)
.addPanel(
    grafana.singlestat.new(
      'CPU Usage',
      datasource='$datasource',
      format='percent',
      gaugeShow=true,
      height=singlestatGuageHeight,
      span=3,
      thresholds='60,80',
    )
    .addTarget(
      grafana.prometheus.target(
        'sum (rate (container_cpu_usage_seconds_total{}[1m])) / sum (machine_cpu_cores) * 100',
      )
    )
)
.addPanel(
    grafana.singlestat.new(
      'Memory Usage',
      datasource='$datasource',
      format='percent',
      gaugeShow=true,
      height=singlestatGuageHeight,
      span=3,
      thresholds='80,90',
    )
    .addTarget(
      grafana.prometheus.target(
        ' sum(container_memory_usage_bytes{})',
      )
    )
)
.addPanel(
    grafana.singlestat.new(
      'Disk Usage',
      datasource='$datasource',
      format='percentunit',
      gaugeShow=true,
      height=singlestatGuageHeight,
      span=3,
      thresholds='80,90',
    )
    .addTarget(
      grafana.prometheus.target(
        '(sum (node_filesystem_size_bytes) - sum (node_filesystem_free_bytes)) / sum (node_filesystem_size_bytes)',
      )
    )
)
.addPanel(
  grafana.graphPanel.new(
    'Http Status Codes',
    datasource='$datasource',
    span=3,
  )
 .addTarget(
   grafana.prometheus.target(
    'echo_request_duration_seconds_count{code="2*"}',
   )
  )
  .addTarget(
    grafana.prometheus.target(
    'echo_request_duration_seconds_count{code="4*"}',
    )
  )
  .addTarget(
    grafana.prometheus.target(
    'echo_request_duration_seconds_count{code="5*"}',
    )
  )
)
.addPanel(
  grafana.graphPanel.new(
    'Block Height',
    datasource='$datasource',
    span=1,
  )
  .addTarget(
    grafana.prometheus.target(
      'max_ethereum_block',
    )
  )
)
.addPanel(
  grafana.graphPanel.new(
    'Start Up Time',
    datasource='$datasource',
    span=1,
  )
  .addTarget(
    grafana.prometheus.target(
        'start_up_time_ms',
        )
    )
)
.addPanel(
  grafana.graphPanel.new(
    'Upstreams Metrics',
    datasource='$datasource',
    span=3,
  )
  .addTarget(
    grafana.prometheus.target(
      'upstreams_healthy',
    )
  )
  .addTarget(
    grafana.prometheus.target(
      'upstreams_archive',
    )
  )
  .addTarget(
    grafana.prometheus.target(
      'upstreams_configured',
    )
  )
)
.addPanel(
  grafana.graphPanel.new(
    'Latency',
    datasource='$datasource',
    span=3,
  )
  .addTarget(
    grafana.prometheus.target(
      'latency_eth_syncing',
    )
  )
  .addTarget(
    grafana.prometheus.target(
      'latency_eth_get_block_number',
    )
  )
  .addTarget(
    grafana.prometheus.target(
      'latency_eth_get_balance',
    )
  )
)