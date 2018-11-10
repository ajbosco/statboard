# statboard

[![Travis CI](https://img.shields.io/travis/ajbosco/statboard.svg?style=flat-square)](https://travis-ci.org/ajbosco/statboard)
[![Go Report Card](https://goreportcard.com/badge/github.com/ajbosco/statboard?style=flat-square)](https://goreportcard.com/report/github.com/ajbosco/statboard)

Personal dashboard and metrics collector

- [Supported Metrics](#supported-metrics)
- [Components](#components)
- [Deployment](#deployment)
- [Setup](#setup)
  * [Configuration](#configuration)
  * [Environment Variables](#environment-variables)

![screenshot](/img/screenshot_v1.png)

### Supported Metrics

* [Github](https://developer.github.com/v3/) Contributions
* [Fitbit](https://dev.fitbit.com/build/reference/web-api/) Steps
* [Goodreads](https://www.goodreads.com/api) Books Read
* [Goodreads](https://www.goodreads.com/api) Pages Read

### Components

`collector` - an application that parses data from a number of APIs, aggregates it into monthly metrics, and stores the results in the `store`

`reporter` - a http service that reads data from `store`, generates [chart.js](https://www.chartjs.org/) charts, and serves a dashboard  

`store` - data is stored in [BoltDB](https://github.com/etcd-io/bbolt) using [Storm](https://github.com/asdine/storm)

### Deployment

This project is intended to be deployed via Docker with two containers (`collector` and `reporter`) and a shared volume for the backing database. The the `collector` application should be a scheduled job such as a [CronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) in Kubernetes.

### Setup

#### Configuration

Populate a configuration yaml file ([example](/config.example.yml)) with the credentials required for the metrics you are interested in collecting.

* Goodreads - create a Developer Key [here](https://www.goodreads.com/api/keys)
* Fitbit - register your application [here](https://dev.fitbit.com/apps/new)
* Github - create a Personal Token [here](https://github.com/settings/tokens)

#### Environment Variables

Statboard requires two environment variables to be set:

* `STATBOARD_CONFIGFILEPATH` - path to your yaml configuration file described above
* `STATBOARD_DBFILEPATH` - path to store `BoltDB` database file


