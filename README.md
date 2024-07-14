<p align="center">
    <h1 align="center">go-therm</h1>
</p>
<p align="center">
    <em>Simple Tasmota-based heating control</em>
</p>
<p align="center">
	<img src="https://img.shields.io/github/license/monishth/go-therm?style=default&logo=opensourceinitiative&logoColor=white&color=0080ff" alt="license">
	<img src="https://img.shields.io/github/last-commit/monishth/go-therm?style=default&logo=git&logoColor=white&color=0080ff" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/monishth/go-therm?style=default&color=0080ff" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/monishth/go-therm?style=default&color=0080ff" alt="repo-language-count">
<p>
<p align="center">
	<!-- default option, no dependency badges. -->
</p>

## Overview

go-therm is a project designed to take existing Tasmota enabled temperature sensors and heating control valves and use this information to _smartly_ control heating based on user-defined targets.
It uses InfluxDB2 to store historical temperature data as there is some value to having this data queryable in the future.

## Getting Started

**System Requirements:**

- **Go**
- **InfluxDB2** (Docker or otherwise)

### Installation

<h4>From <code>source</code></h4>

> 1. Clone the go-therm repository:
>
> ```console
> $ git clone https://github.com/monishth/go-therm
> ```
>
> 2. Change to the project directory:
>
> ```console
> $ cd go-therm
> ```
>
> 3. Install the dependencies:
>
> ```console
> $ go build -o go-therm
> ```

### Usage

<h4>From <code>source</code></h4>

> Run go-therm using the command below:
>
> ```console
> $ ./go-therm
> ```

### Tests

> Run the test suite using the command below:
>
> ```console
> $ go test
> ```

---

## Project Roadmap

- [x] Tasmota state handlers
- [x] PID controllers
- [x] Basic frontend for setting target per zone
- [ ] Set valve state based on PID output
- [ ] Schedule system (core + frontend)
