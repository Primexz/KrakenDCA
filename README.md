# 🐙 Kraken-DCA

**Automated dollar cost averaging on the Kraken exchange**

![Buid](https://img.shields.io/github/actions/workflow/status/primexz/krakendca/release.yml)
![License](https://img.shields.io/github/license/primexz/krakendca)

## Table of Contents
1. ➤ [About the project](#-about-the-project)
    - [Orders](#-orders)
2. ➤ [Configuration](#-configuration)
3. ➤ [Run with Docker](#-run-with-docker)
4. ➤ [Run without Docker](#-run-without-docker)

## 🔍 About the project

FooBar

### 💰 Orders

#### When will orders be placed?

## ⚙️ Configuration

This tool is configured via environment variables. Some environment variables are required and some activate additional functionalities.


| Variable | Description | Required | Default |
| --- | --- | --- | --- |
| `KRAKEN_PUBLIC_KEY` | Your Kraken API public key | ✅ | |
| `KRAKEN_PRIVATE_KEY` | Your Kraken API private key | ✅ | |
| `CURRENCY` | Your fiat currency to be used, e.g. USD or EUR | ❌  | `USD` |
| `KRAKEN_ORDER_SIZE` | The order size to be used. This value should only be edited if you know exactly what you are doing.r | ❌ | `0.0001` |
| `CHECK_DELAY` | How often the algorithm should be executed, in seconds. | ❌ | `60` |


## 🐳 Run with Docker

## 💻 Run without Docker