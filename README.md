# shawty
A URL shortener (visit http://luu.bz)

## Build status
[![Build Status](https://secure.travis-ci.org/3fps/shawty.png)](http://travis-ci.org/3fps/shawty)

## Get started
1. Code setup
  1. Make sure you have [GOPATH](http://golang.org/cmd/go/#hdr-GOPATH_environment_variable) setup correctly.
  2. Put `shawty` at `$GOPATH/src/github.com/3fps/shawty`
  3. [Set configurations](https://github.com/3fps/shawty/wiki/Configuration). For DB configuration, please see the examples at the [Go-MySQL-Driver](https://github.com/go-sql-driver/MySQL#examples) or [Postgres driver](https://github.com/lib/pq#use) page.
2. Database setup
  1. Create a database and run the setup script at `shawty/install/<your db type>/db.sql`
3. Run project
  1. Go to the shawty diretory: `cd $GOPATH/src/github.com/3fps/shawty`
  2. Execute the run script (might need `sudo` if the port is set to 80): `source run`

## History

### Version 1
* 1.2.2: Fixed run script
* 1.2.1: Minor fixes
* 1.2.0: Added support for Postgres (**This update requires configuration change**).
* 1.1.8: Give run script execute permission from the repository
* 1.1.7: Fixed run script not calling config file properly
* 1.1.6: Separated run script from the configuration
* 1.1.5: Index displays version number
* 1.1.4: Updated import paths
* 1.1.3: Updated to use new version of logger ([Issue #6](https://github.com/3fps/shawty/issues/6))
* 1.1.2: Fixed [Issue #5](https://github.com/3fps/shawty/issues/5)
* 1.1.1: SSL support for bookmarklet ([Issue #4](https://github.com/3fps/shawty/issues/4))
* 1.1.0: Rate limit
* 1.0.1
  * GOMAXPROCS support
  * Various bug fixes and enhancements
* 1.0: Initial release
