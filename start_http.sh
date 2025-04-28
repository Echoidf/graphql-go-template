#!/bin/bash
socat TCP-LISTEN:10000,reuseaddr,fork UNIX-CONNECT:/Users/macbookpro/.uds/gqlexample.sock

