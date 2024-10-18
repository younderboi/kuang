# Kuang

## Building binaries

Install garble for obfuscation

```sh
go install mvdan.cc/garble@latest
```

### **Building the agent binary**

**Windows**

```sh
GOOS=windows GOARCH=amd64 garble build -tiny -o bin/agent/win/cute_poppy.exe cmd/tcp_agent/main.go
```

**Linux**

```sh
GOOS=linux GOARCH=amd64 garble build -tiny -o bin/agent/linux/cute_poppy cmd/tcp/agent/main.go
```

## Tests

I write totally perfect, absolutely bug free code, so why bother??
