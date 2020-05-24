FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY chat .
COPY web web
EXPOSE 80
CMD ["./chat"]

