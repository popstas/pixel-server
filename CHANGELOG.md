## NEXT VERSION

BREAKING CHANGES:
  * default port changed from 8080 to 8246

FEATURES:
  * AnyBar support for MacOS
  * add default pixel brightness to settings

IMPROVEMENTS:

BUG FIXES:
  * pixel state not saved from /kapacitor signals

## v0.2.1 (31.08.2016)

BUG FIXES:
  * fixed broken /kapacitor handler

## v0.2 (31.08.2016)

IMPROVEMENTS:
  * code decoupled to packages
  * server rewrited from net/http to gin-gonic/gin
  * add release binaries for windows and linux

## v0.1 (30.08.2016)

FEATURES:
  * receive status to `/status?value=50&message=msg\second line&blink=1&brightness=100`
  * receive status from kapacitor post to `/kapacitor`
