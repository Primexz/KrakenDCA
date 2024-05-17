# 🐙 Kraken-DCA

**Automated dollar cost averaging on the Kraken exchange**

![Buid](https://img.shields.io/github/actions/workflow/status/primexz/krakendca/release.yml)
![License](https://img.shields.io/github/license/primexz/krakendca)

## Table of Contents
1. ➤ [About the project](#-about-the-project)
    - [Orders](#-orders)
2. ➤ [Configuration](#-configuration)
3. ➤ [Push Notifications](#-push-notifications)
4. ➤ [Run with Docker](#-run-with-docker)
    - [Docker-CLI](#-docker-cli)
    - [Docker-Compose](#-docker-compose)
5. ➤ [Run without Docker](#-run-without-docker)

## 🔍 About the project

Since Kraken has extremely high fees for its crypto savings plan, this tool was developed to invest regularly in cryptocurrencies. Since Kraken offers two platforms: ‘Kraken’ and ‘Kraken Pro’, we take advantage of the fact that Kraken Pro has very low fees. All you have to do is regularly deposit fiat currencies on Kraken Pro and the bot does everything else for you.

### 💰 Orders

#### When will orders be placed?

Bitcoin orders are placed as often as possible. To illustrate the behaviour in more detail, let's look at the following example:
You deposit €500 per month on Kraken-Pro.
The bot calculates how often you can buy Bitcoin this month, as Kraken has a minimum purchase limit of 0.0001 BTC. Your orders are executed as often as possible throughout the month, thereby achieving dollar-cost averaging.

## ⚙️ Configuration

This tool is configured via environment variables. Some environment variables are required and some activate additional functionalities.


| Variable | Description | Required | Default |
| --- | --- | --- | --- |
| `KRAKEN_PUBLIC_KEY` | Your Kraken API public key | ✅ | |
| `KRAKEN_PRIVATE_KEY` | Your Kraken API private key | ✅ | |
| `CURRENCY` | Your fiat currency to be used, e.g. USD or EUR | ❌  | `USD` |
| `KRAKEN_ORDER_SIZE` | The order size to be used. This value should only be edited if you know exactly what you are doing. | ❌ | `0.0001` |
| `LIMIT_ORDER_MODE` | If set to true, limit orders are placed. With a normal monthly volume, you only pay 0.25% fees per purchase instead of 0.4%. | ❌ | `false` |
| `CHECK_DELAY` | How often the algorithm should be executed, in seconds. | ❌ | `60` |
| `GOTIFY_URL` | URL to your Gotify server | ❌ |  |
| `GOTIFY_APP_TOKEN` | App token for the app on the Gotify server | ❌ |  |

## 📱 Push Notifications

The environment variables `GOTIFY_URL` and `GOTIFY_APP_TOKEN` can be used to activate Gotify Push Notifications. As soon as a purchase has been made, you will immediately receive a notification so that you always have a full overview of your purchases.


### 🐳 Run with Docker

###  Docker-CLI

```bash
docker run -d --name kraken_dca \
  -e KRAKEN_PUBLIC_KEY=your-public-key \
  -e KRAKEN_PRIVATE=your-private-key \
  -e CURRENCY=EUR \
  ghcr.io/primexz/kraken_dca:latest

```


### 🚀 Docker-Compose

```bash
vim docker-compose.yml
```

```yaml
version: "3.8"
services:
  kraken_dca:
    image: ghcr.io/primexz/kraken_dca:latest
    environment:
      - KRAKEN_PUBLIC_KEY=your-public-key
      - KRAKEN_PRIVATE_KEY=your-private-key
      - CURRENCY=EUR
    restart: always
```

```bash
docker-compose up -d
```


## 💻 Run without Docker

This tool can be run directly with Go for development.

```bash
KRAKEN_PUBLIC_KEY=your-public-key KRAKEN_PRIVATE_KEY=your-private-key go run .
```