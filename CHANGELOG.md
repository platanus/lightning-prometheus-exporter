# Changelog

## Unreleased

* Fix wallet ballance to use a single metric with labels
  `wallet_balance_satoshis_total`, with `confirmed` and `unconfirmed` labels
* Change default prefix to `lnd`
* Remove unused flags

## 0.1.0

* Initial release.
* Support to wallet balance metrics
  * `TotalBallance`
  * `ConfirmedBalance`
  * `UnconfirmedBalance`