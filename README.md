# CPUIdle exporter

Prometheus exporter for CPUIdle driver metrics exposed by **Linux** kernel. \
(See Also: https://www.kernel.org/doc/Documentation/cpuidle/sysfs.txt)

## Installation and Usage

### using [Docker image](https://github.com/masan4444/cpuidle_exporter/pkgs/container/cpuidle-exporter)

```bash
docker run -d \
  -p 9975:9975 \
  ghcr.io/masan4444/cpuidle-exporter:latest
```

### using docker-compose

#### docker-compose.yaml

```yaml
version: "3"
services:
  cpuidle-exporter:
    image: ghcr.io/masan4444/cpuidle-exporter:latest
    container_name: cpuidle-exporter
    ports:
      - 9975:9975
    restart: always
```
