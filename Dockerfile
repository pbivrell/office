FROM golang:1.15-alpine AS build

WORKDIR /src/
COPY main.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/office main.go

FROM scratch
COPY --from=build /bin/office /bin/office
COPY html /bin/html/
ENTRYPOINT ["/bin/office"]
CMD ["-html","/bin/html" ]
