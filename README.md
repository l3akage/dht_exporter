# dht_exporter

Prometheus exporter for temperature/humidity data provided by DHT11/DHT22 sensors connected to a Raspberry PI GPIO

# Usage

```
./dht_exporter
```

# Options

Name     | Default | Description
---------|-------------|----
--version || Print version information
--listen-address | :9330 | Address on which to expose metrics.
--path | /metrics | Path under which to expose metrics.
--device | 22 | Sensor type, either 11 or 22 for DHT11/DHT22
--names | names.yaml | File mapping GPIOs to names

# Configuration

The names.yaml files contains the gpio and name to use
```
names:
  gpio: dht_sensor_01
```

# Example output
```
# HELP dht_humidity_percent Humidity
# TYPE dht_humidity_percent gauge
dht_humidity_percent{location="dht_sensor_01"} 99.1
# HELP dht_temp Air temperature (in degrees C)
# TYPE dht_temp gauge
dht_temp{location="dht_sensor_01"} 24.6
# HELP dht_up Scrape was successful
# TYPE dht_up gauge
dht_up{location="dht_sensor_01"} 1
```

# Circuit
Connect pin1 to 3.3V, pin2 with a pull-up resistor (4.7k – 10kΩ) to a gpio pin (i used 4.7k and gpio 4) and pin4 to GND on the PI.

![Circuit](https://m0u.de/assets/dht22_circuit.svg)

