# Choria Fleet Tally

In large dynamic fleets it's hard to keep track of counts and versions of nodes. This is a tool that can observe 
a running network and gather versions of a specific component. The results are exposed as Prometheus metrics.

```nohighlight
$ tally --component '*' --port 8080 --prefix choria_tally
```

For this to work it uses the normal Choria client configuration to connect to the right middleware using TLS and
listen there. The certificates needed does not need to match the `.+.mcollective` pattern, meaning tally will never
make RPC requests to your fleet - it passively listens for advisories.

This will listen on port 8080 for /metrics, it will observe events from the server component and expose metrics
as below:

| Metric                       | Description                                                  |
|------------------------------|--------------------------------------------------------------|
| choria_tally_good_events     | Events processed successfully                                |
| choria_tally_process_errors  | The number of events received that failed to process         |
| choria_tally_event_types     | The number of events received by type                        |
| choria_tally_versions        | Gauge indicating the number of running components by version |
| choria_tally_processing_time | The time taken to process events                             |
| choria_tally_nodes_expired   | The number of nodes removed during maintenance runs          |

Additionally, this tool can also watch [Choria Autonomous Agent](https://choria.io/docs/autoagents/) events, today it supports transition events and exec watchers:

| Metric                             | Description                                                             |
|------------------------------------|-------------------------------------------------------------------------|
| choria_tally_machine_transition    | Information about transition events handled by Choria Autonomous Agents |
| choria_tally_exec_watcher_success  | Machine exec watcher success runs                                       |
| choria_tally_exec_watcher_failures | Number of exec watcher executions that failed                           |
| choria_tally_exec_watcher_runtime  | Exec watcher runtime                                                    |


If you have any [Choria Governors](https://choria.io/docs/streams/governor/) the tool can listen to the events
these emit and report on those.

| Metric                | Description                                                |
|-----------------------|------------------------------------------------------------|
| choria_tally_governor | Events and their types seen per Governor and per Component |

## Clustered Deployment

Clustered deployments are problematic since from the perspective of Prometheus the stats will be multiplied by the number of running instances.

When ran with the `--election TALLY` flag it will activate leader election against Choria Streams.

Regardless all metrics will have a `active` label set which will be `1` on the active node (or only node when leader election is not configured).