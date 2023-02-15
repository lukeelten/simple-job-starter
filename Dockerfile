FROM registry.access.redhat.com/ubi9/go-toolset:1.18 as builder

WORKDIR /app
COPY . .
RUN go build -o bash-job-starter . 

FROM registry.access.redhat.com/ubi9/ubi:latest

WORKDIR /
COPY --from=builder /app/bash-job-starter .

ENTRYPOINT [ "/bash-job-starter" ]
