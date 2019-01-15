# Subprocess

Subprocess abstracts running processes with restarts.

## Install

`go get -u github.com/jakubdal/subprocess`

## Usage

### Run a process that restarts on error

```go
proc, err := subprocess.NewProcess(context.Background(), nil, "sleep", nil, "3s")
```

### Stop the process

```go
proc.Stop()
```

### Send a signal to process

```go
proc.Signal(os.Interrupt)
```
