# TXT-Summarizer

Summarize a value key pairs in txt files structured as:

```txt
1 foo
12 bar
```

## Running on windows system

Copy ready made exe to the folder of the files and:


```ps
summarize file1.txt file2.txt filex.txt
```

## Create windows executable

On windows:

```bash
go build -o summarize.exe cmd/main.go
```

on sane operating system:

```bash
GOOS=windows GOARCH=amd64 go build -o summarize.exe cmd/main.go
```

## Compile

```bash
go build -o summarize cmd/main.go
```

Run

```bash
./summarize file1.txt file2.txt filex.txt
```

## Tests

```bash
go test ./... -v
```

