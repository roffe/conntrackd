# Conntrackd

Prometheus exporter that gathers `conntrack -S` statistics and exposes

Uses port 2112

Metrics exposed behind prometheus default:

    conntrack_found
    conntrack_invalid
    conntrack_ignore
    conntrack_insert
    conntrack_insert_faile,
    conntrack_drop
    conntrack_early_drop
    conntrack_error
    conntrack_search_start
