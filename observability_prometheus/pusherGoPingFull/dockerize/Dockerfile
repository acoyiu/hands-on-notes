FROM alpine:latest  

WORKDIR /

COPY ./etable /

# if need debug, ENTRYPOINT [ "sh", "-c", "tail -f /dev/null" ]
ENTRYPOINT [ "/etable" ]