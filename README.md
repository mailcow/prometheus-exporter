# mailcow exporter

[Prometheus](https://prometheus.io) exporter for information about a
[mailcow](https://github.com/mailcow/mailcow-dockerized) instance.

## Usage

```bash
# As a docker container
$ docker run -p '9099:9099' thej6s/mailcow-exporter

# Natively
$ ./mailcow-exporter
```

The `/metrics` endpoint requires `host` and `apiKey` URL parameters. `host` is the
hostname of your mailcow instance. `apiKey` should be a readonly API key that can
be generated by logging into the mailcow management interface and navigating to
'Access > API'.

![Visualization of where to find the API Key](./.readme/api-key)

The following prometheus configuration can be used in order to pass these information
to the endpoint:

```yaml
scrape_configs:
  - job_name: 'mailcow'
    static_configs:
      - targets: [ 'mailcow_exporter:9099' ]
    params:
      host: [ 'mail.example.com' ]
      apiKey: [ 'YOUR-APIKEY-HERE' ]
```

## Example metrics

```
mailcow_mailbox_last_login{host="mail.example.com",mailbox="foo@bar.com"} 1.599255303e+09
mailcow_mailbox_last_login{host="mail.example.com",mailbox="test@bar.com"} 1.599247706e+09

mailcow_mailbox_messages{host="mail.example.com",mailbox="foo@bar.com"} 23476
mailcow_mailbox_messages{host="mail.example.com",mailbox="test@bar.com"} 1891

mailcow_mailbox_quota_allowed{host="mail.example.com",mailbox="foo@bar.com"} 3.221225472e+09
mailcow_mailbox_quota_allowed{host="mail.example.com",mailbox="test@bar.com"} 3.221225472e+09

mailcow_mailbox_quota_used{host="mail.example.com",mailbox="foo@bar.com"} 1.919023167e+09
mailcow_mailbox_quota_used{host="mail.example.com",mailbox="test@bar.com"} 1.844312552e+09

mailcow_mailq{host="mail.example.com",queue="deferred",sender="foo@bar.com"} 2
mailcow_mailq{host="mail.example.com",queue="deferred",sender="test@bar.com"} 1

# HELP mailcow_quarantine_age Age of quarantined items in seconds
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="foo@bar.com",le="10800"} 0
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="foo@bar.com",le="43200"} 0
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="foo@bar.com",le="86400"} 0
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="foo@bar.com",le="259200"} 3
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="foo@bar.com",le="604800"} 9
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="foo@bar.com",le="1.2096e+06"} 12
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="foo@bar.com",le="2.592e+06"} 41
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="foo@bar.com",le="+Inf"} 147
mailcow_quarantine_age_sum{host="mail.example.com",recipient="foo@bar.com"} 1.301292926e+09
mailcow_quarantine_age_count{host="mail.example.com",recipient="foo@bar.com"} 147
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="test@bar.com",le="10800"} 0
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="test@bar.com",le="43200"} 0
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="test@bar.com",le="86400"} 0
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="test@bar.com",le="259200"} 0
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="test@bar.com",le="604800"} 0
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="test@bar.com",le="1.2096e+06"} 0
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="test@bar.com",le="2.592e+06"} 0
mailcow_quarantine_age_bucket{host="mail.example.com",recipient="test@bar.com",le="+Inf"} 2
mailcow_quarantine_age_sum{host="mail.example.com",recipient="test@bar.com"} 2.7138547e+07
mailcow_quarantine_age_count{host="mail.example.com",recipient="test@bar.com"} 2

mailcow_quarantine_score_bucket{host="mail.example.com",recipient="foo@bar.com",le="0"} 0
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="foo@bar.com",le="10"} 0
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="foo@bar.com",le="20"} 41
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="foo@bar.com",le="40"} 122
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="foo@bar.com",le="60"} 136
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="foo@bar.com",le="80"} 141
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="foo@bar.com",le="100"} 141
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="foo@bar.com",le="+Inf"} 147
mailcow_quarantine_score_sum{host="mail.example.com",recipient="foo@bar.com"} 16225.91000000001
mailcow_quarantine_score_count{host="mail.example.com",recipient="foo@bar.com"} 147
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="test@bar.com",le="0"} 0
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="test@bar.com",le="10"} 0
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="test@bar.com",le="20"} 0
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="test@bar.com",le="40"} 0
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="test@bar.com",le="60"} 0
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="test@bar.com",le="80"} 0
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="test@bar.com",le="100"} 0
mailcow_quarantine_score_bucket{host="mail.example.com",recipient="test@bar.com",le="+Inf"} 2
mailcow_quarantine_score_sum{host="mail.example.com",recipient="test@bar.com"} 3988.03
mailcow_quarantine_score_count{host="mail.example.com",recipient="test@bar.com"} 2

mailcow_quarantine_count{host="mail.example.com",is_virus="0",recipient="foo@bar.com"} 147
mailcow_quarantine_count{host="mail.example.com",is_virus="0",recipient="test@bar.com"} 2
mailcow_quarantine_count{host="mail.example.com",is_virus="1",recipient="foo@bar.com"} 0
mailcow_quarantine_count{host="mail.example.com",is_virus="1",recipient="test@bar.com"} 0
```
