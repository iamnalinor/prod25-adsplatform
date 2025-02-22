FROM grafana/grafana:latest

ENV GF_AUTH_ANONYMOUS_ENABLED=true

COPY provisioning /etc/grafana/provisioning
COPY dashboards /var/lib/grafana/dashboards
COPY grafana.ini /etc/grafana/grafana.ini

HEALTHCHECK --interval=5s --timeout=5s --retries=5 CMD curl --fail http://localhost:3000/api/health
