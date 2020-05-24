FROM alpine:3.11
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY chat .
COPY web web
EXPOSE 80
CMD ["./chat"]

