# Changelog

## Unreleased

* Add `process` and `go` metrics.

## 0.2.0

* Fix wallet ballance to use a single metric with labels
  `wallet_balance_satoshis`, with `confirmed` and `unconfirmed` labels
* Change default prefix to `lnd`
* Remove unused flags
* Add `channels` metric, with `pending`, `active` and `inactive` labels
* Add `block_height` metric
* Add `synced_to_chain` metric
* Add channel related metrics
  * `channels_limbo_balance_satoshis`
  * `channels_pending`
  * `channels_waiting_close`

## 0.1.0

* Initial release.
* Support to wallet balance metrics
  * `TotalBallance`
  * `ConfirmedBalance`
  * `UnconfirmedBalance`