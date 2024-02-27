# 1BRCgo

Study project Based on [The One Billion Row Challenge](https://github.com/gunnarmorling/1brc)
and [This Repository](https://github.com/shraddhaag/1brc)

## Run

Generate a CSV file with 1.000.000 rows:

```bash
go run cmd/generate/main.go -size=1_000_000 -output=1brc.csv
```

Parse the csv file

```bash
go run cmd/parse/main.go 1brc.csv
```
