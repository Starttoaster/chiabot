FROM alpine:latest

COPY GOBINARY /

CMD ["/GOBINARY --token $TOKEN --interval $CHECK_INTERVAL --channel-id $CHANNEL_ID"]