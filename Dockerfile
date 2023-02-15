FROM golang:1.20 as builder

WORKDIR /app
COPY . .
RUN go build -o bash-job-starter . 

FROM centos:7

WORKDIR /
COPY --from=builder /app/bash-job-starter .

ENTRYPOINT [ "/bash-job-starter" ]
