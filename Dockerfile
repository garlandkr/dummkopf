FROM golang:1.14-alpine AS build
RUN apk --no-cache add make
WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .
RUN make OUT=/bin/dummkopf

FROM scratch
COPY --from=build /bin/dummkopf /bin/dummkopf
EXPOSE 9000
CMD ["/bin/dummkopf"]
