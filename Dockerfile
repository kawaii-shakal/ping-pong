FROM golang:1.22.8-bookworm AS builed

WORKDIR /src
COPY app.go go.mod /src/

RUN CGO_ENABLED=0 GOOS=linux  go build -o /bin/pingpong 

FROM scratch

COPY --from=builed /bin/pingpong /pingpong
ENTRYPOINT ["/pingpong"]