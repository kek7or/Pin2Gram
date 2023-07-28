FROM golang:1.19 as build 
  
WORKDIR /app 
COPY go.mod . 
COPY go.sum . 
COPY Makefile . 
RUN go mod download 

COPY . . 
RUN make 

FROM alpine:3.16 
COPY --from=build /app/bin/example /usr/bin/ 

EXPOSE 8080 

CMD ["/usr/bin/example"]
