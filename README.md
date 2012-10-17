[![Build Status](https://secure.travis-ci.org/3fps/shawty.png)](http://travis-ci.org/3fps/shawty)

# shawty
A URL shortener (visit http://luu.bz) 

## Get started
1. Code setup
  1. Make sure you have [GOPATH](http://golang.org/cmd/go/#GOPATH_environment_variable) setup correctly.
  2. Put `shawty` at `$GOPATH/src/go.3fps.com/shawty`
  3. Make configure changes to `shawty/run`. For DB configuration, please see the examples at the [go-mysql-driver](http://code.google.com/p/go-mysql-driver/#Examples) page.
2. Database setup (MySQL)
  1. Create a MySQL database and run the setup script at `shawty/install/db.sql`
3. Run project
  1. Go to the shawty diretory: `cd $GOPATH/src/go.3fps.com/shawty`
  2. Execute the run script (might need `sudo` if the port is set to 80): `./run`