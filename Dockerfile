FROM scratch

ADD server static /

EXPOSE 8080

ENTRYPOINT [ "/server" ]
