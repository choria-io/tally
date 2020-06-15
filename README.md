# Choria Fleet Tally

In large dynamic fleets it's hard to keep track of counts and versions of nodes. This is a tool that can observe 
a running network and gather versions of a specific component. The results are exposed as Prometheus metrics.

```nohighlight
$ tally --component server --port 8080 --prefix choria_tally
```

For this to work it uses the normal Choria client configuration to connect to the right middleware using TLS and
listen there. The certificates needed does not need to match the `.+.mcollective` pattern, meaning tally will never
make RPC requests to your fleet - it passively listens for advisories.

This will listen on port 8080 for /metrics, it will observe events from the server component and expose metrics
as below:

|Metric|Description|
|------|-----------|
|choria_tally_good_events|Events processed successfully|
|choria_tally_process_errors|The number of events received that failed to process|
|choria_tally_event_types|The number of events received by type|
|choria_tally_versions|Gauge indicating the number of running components by version|
|choria_tally_maintenance_time|Time spent doing regular maintenance on the stored data|
|choria_tally_processing_time|The time taken to process events|
|choria_tally_nodes_expired|The number of nodes removed during maintenance runs|

Additionally, this tool can also watch Choria Autonomous Agent events, today it supports transition events only:

|Metric|Description|
|------|-----------|
|choria_tally_machine_transition|Information about transition events handled by Choria Autonomous Agents|

Here the prefix - `choria_tally` - is what would be the default if you didn't specify `--prefix`.
