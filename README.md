# conval - The CONtest eVALuator

This little Go library helps to evaluate the log files from amateur radio contests in several ways:

- calculate the score of a log file based on a given rule set (work in progress)
- show the performance over time during the contest (not yet implemented)
- calculate statistics (not yet implemented)
- compare the performance of two log files (not yet implemented)

Log files can be provided in [ADIF](https://www.adif.org/) or [Cabrillo](https://wwrof.org/cabrillo/) format. The results are either provided as plain text, CSV, or JSON.

## Use as a Go Library

To include `conval` into your own projects as a library:

```shell
go get github.com/ftl/conval
```

## Use as a CLI Tool (work in progress)

`conval` also includes a simple CLI tool that is mainly used to demonstrate the integration of the library.

Build it:

```shell
go build -o conval ./cmd
```

Simply run it:

```shell
go run ./cmd
```

## License
This software is published under the [MIT License](https://www.tldrlegal.com/l/mit).

Copyright [Florian Thienel](http://thecodingflow.com/)
