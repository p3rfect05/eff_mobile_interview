FROM alpine:latest
RUN mkdir /app

COPY ./carAPI /app


# Run the server executable
CMD [ "/app/carAPI" ]