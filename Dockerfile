FROM debian

RUN mkdir -p /doc/page/ \
    && touch /doc/page/test.md \
    && echo "# Test Page" >> /doc/page/test.md \
    && echo "Hello world" >> /doc/page/test.md

COPY go-cms /usr/local/bin/go-cms

EXPOSE 7335

CMD ["go-cms", "-port=7335", "-path=/doc"]
