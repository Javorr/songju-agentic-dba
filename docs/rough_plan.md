# Telemetry Platform

An end-to-end telemetry platform consisting of STM32 firmware and a backend service.

---

# Goals

- Design a communication protocol between an embedded device and a backend.
- Model and store telemetry efficiently.
- Handle unreliable networks gracefully.
- Produce clean documentation and architecture diagrams.

# High Level Architecture

```
                +----------------+
                | STM32 BlackPill|
                +--------+-------+
                         |
                  HTTP / Serial
                         |
                         v
                +----------------+
                | Backend API    |
                +--------+-------+
                         |
                  PostgreSQL
                         |
                         v
                Dashboard / API
```

---

# Version 1 Scope

## Firmware

The firmware should:

- Generate a unique device ID.
- Track uptime.
- Read MCU internal temperature.
- Read MCU supply voltage (or simulated battery voltage).
- Collect telemetry every configurable interval.
- Store configuration in memory.
- Retry failed uploads.
- Cache unsent telemetry in RAM.
- Log important events over UART.

### Telemetry Payload (this is subject to change)

```json
{
  "deviceId": "...",
  "timestamp": "...",
  "temperature": 24.6,
  "voltage": 3.29,
  "uptimeSeconds": 531,
  "firmwareVersion": "1.0.0"
}
```

---

## Backend

The backend should provide a REST API.

### Device Registration

```
POST /devices/register
```

Registers a device and returns:

- Authentication token
- Reporting interval
- Configuration

---

### Upload Telemetry

```
POST /telemetry
```

Receives telemetry from devices.

Responsibilities:

- Authenticate device
- Validate payload
- Store measurement
- Update device status

---

### Device Endpoints

```
GET /devices

GET /devices/{id}

GET /devices/{id}/telemetry
```

---

### Health Tracking

Track:

- Last seen
- Firmware version
- Online/offline status
- Number of received measurements

---

## Database

Possible structure (this is subject to change):

### devices

```
id
device_id
firmware_version
registered_at
last_seen
status
report_interval
```

### telemetry

```
id
device_id
timestamp
temperature
voltage
uptime
```

---

# Dashboard

Things that I could display maybe:

- Registered devices
- Online devices
- Offline devices
- Last seen
- Latest telemetry
- Historical temperature graph
- Historical voltage graph

---

# Communication

For Version 1:

- HTTP
- JSON
- Simple Bearer token authentication

Later versions may explore binary protocols.

---

# Non-Functional Requirements

The project should prioritize:

- Readability
- Maintainability
- Simplicity
- Good documentation
- Modular architecture
- Proper error handling
- Logging
- Unit tests where practical

---

# Technologies

## Firmware

- C17
- STM32 HAL
- STM32CubeMX

## Backend

- Go
- PostgreSQL
- Docker

---

# Future Improvements

These are intentionally out of scope until Version 1 is complete.

## Device Configuration

```
POST /devices/{id}/configuration
```

Update:

- Reporting interval
- Enabled sensors
- Device name

---

## OTA Firmware Updates

```
GET /firmware/latest
```

Device:

- Checks version
- Downloads firmware
- Verifies checksum
- Performs update

---

## Persistent Offline Storage

Instead of storing unsent telemetry in RAM:

- External EEPROM
- SD card
- Flash circular buffer

Here I could use something similar to what I used for the POS devices at work.

---

## Alerts

Notify when:

- Device offline
- Temperature exceeds threshold
- Voltage too low

---

## Multiple Devices

Support many STM32 boards simultaneously.

---

## Binary Protocol

Replace JSON with a compact binary format to reduce bandwidth.

---

## Authentication Improvements

Replace static tokens with:

- Device certificates
- Signed requests
- Key rotation

---

# Success Criteria

The project is considered complete when:

- A STM32 board successfully registers itself.
- Telemetry is periodically uploaded.
- Failed uploads recover automatically.
- Data is stored in PostgreSQL.
- Historical telemetry can be viewed.
- Device health is visible.
- The codebase is documented and easy to understand.
