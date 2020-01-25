FROM scratch
ADD sessionsrv.elf /app/
CMD ["/app/sessionsrv.elf"]
EXPOSE 8090