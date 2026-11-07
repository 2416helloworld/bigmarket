FROM alpine:3.19
RUN apk add --no-cache tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
WORKDIR /app
COPY big-market-go /app/big-market-go
RUN chmod +x /app/big-market-go
ENV TZ=Asia/Shanghai
EXPOSE 8080
CMD ["/app/big-market-go"]
