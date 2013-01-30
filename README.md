# shawty
A URL shortener (visit http://luu.bz)

## Build status
[![Build Status](https://secure.travis-ci.org/3fps/shawty.png)](http://travis-ci.org/3fps/shawty)

## Get started
1. Code setup
  1. Make sure you have [GOPATH](http://golang.org/cmd/go/#GOPATH_environment_variable) setup correctly.
  2. Put `shawty` at `$GOPATH/src/github.com/3fps/shawty`
  3. Make configure changes to `shawty/run`. For DB configuration, please see the examples at the [Go-MySQL-Driver](https://github.com/Go-SQL-Driver/MySQL#examples) page.
2. Database setup (MySQL)
  1. Create a MySQL database and run the setup script at `shawty/install/db.sql`
3. Run project
  1. Go to the shawty diretory: `cd $GOPATH/src/github.com/3fps/shawty`
  2. Execute the run script (might need `sudo` if the port is set to 80): `./run`

## History

### Version 1
* 1.1.3: Updated to use new version of logger ([Issue #6]()(https://github.com/3fps/shawty/issues/6))
* 1.1.2: Fixed [issue 5](https://github.com/3fps/shawty/issues/5)
* 1.1.1: SSL support for bookmarklet ([Issue #4](https://github.com/3fps/shawty/issues/4))
* 1.1.0: Rate limit
* 1.0.1
  * GOMAXPROCS support
  * Various bug fixes and enhancements
* 1.0: Initial release
