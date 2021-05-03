# N26 plugin for MoneyLover CLI

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/moneylovercli-plugin-n26)](https://github.com/nhatthm/moneylovercli-plugin-n26/releases/latest)
[![Build Status](https://github.com/nhatthm/moneylovercli-plugin-n26/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/moneylovercli-plugin-n26/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/moneylovercli-plugin-n26/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/moneylovercli-plugin-n26)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/moneylovercli-plugin-n26)](https://goreportcard.com/report/github.com/nhatthm/moneylovercli-plugin-n26)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/nhatthm/moneylovercli-plugin-n26)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

A plugin for [moneylovercli](https://github.com/nhatthm/moneylovercl) to convert n26 transactions and import to MoneyLover.

## Prerequisites

- `Go >= 1.15`

## Install

```bash
moneylovercli plugin install github.com/nhatthm/moneylovercli-plugin-n26
```

## Usage

If you have a list of n26 transactions, you can either pipe it or read it, like this

```bash
n26 transactions --from 2020-01-01 --to 2020-02-01 > transactions.json
moneylovercli transactions import --provider n26 < transactions.json
```

or

```bash
n26 transactions --from 2020-01-01 --to 2020-02-01 | moneylovercli transactions import --provider n26
```

## Donation

If you like this project, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
