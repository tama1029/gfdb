## What's this?

`gfdb` is a code-generation tool for Go.
Gen From Database(now only mysql).

Output gorm struct and dto.

## Install

`go get -u github.com/tama1029/gfdb`

## Example

* create output directory

`mkdir example`

* gen from database

`gfdb struct --host 127.0.0.1 --port 3306 --database development --user root`

## Acknowledgments

Inspired by `Shelnutt2/db2struct`.