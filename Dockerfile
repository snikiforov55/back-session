FROM scratch
#ENV REDIS_HOSTNAME=redis
#ENV REDIS_PORT=6379
#ENV REDIS_PASSWORD=""
#ENV REDIS_DB=0
##ENV SESSION_EXP_SEC", session.DefaultSessionExpirationSec)
#ENV SERVICE_PORT=8090

ADD sessionsrv.elf64 /app/
CMD ["/app/sessionsrv.elf64"]
EXPOSE 8090