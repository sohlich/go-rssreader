# go-rssreader
[![Build Status](https://travis-ci.org/sohlich/go-rssreader.svg?branch=master)](https://travis-ci.org/sohlich/go-rssreader)

Simple RSS reader written in GO. Basic projec to learn how to use channels.

Usage:
Modify the rss.source file to add urls for rss feed source, or use --url switch
to read one url from url given in argument.


Read all sources from rss.source file:
```
./rssreader
```

Read one url from command line:
```
./rssreader --url http://servis.idnes.cz/rss.aspx?c=zpravodaj
```