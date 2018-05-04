FROM golang:1.8.3 as builder
LABEL builder=true
COPY src /go/src/
WORKDIR /go/src/dynatrace/cmd/visitflow
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o visitflow .

FROM scratch
LABEL maintainer="kin.wai.koo@dynatrace.com"
LABEL builder=false
COPY --from=builder /go/src/dynatrace/cmd/visitflow/visitflow /
COPY html /html/
EXPOSE 8080

ENTRYPOINT ["/visitflow", "-port", "8080", "-dir", "/html"]

