FROM krakend/builder:2.1.3 as httpClientBuilder

WORKDIR /app

COPY ./plugins/shared/ /shared

COPY ./plugins/http-client/ /app

RUN go build -buildmode=plugin -o http-client.so .

FROM krakend/builder:2.1.3 as httpServerBuilder

WORKDIR /app

COPY ./plugins/shared/ /shared

COPY ./plugins/http-server/ /app

RUN go build -buildmode=plugin -o http-server.so .

FROM krakend/builder:2.1.3 as reqRespModifierBuilder

WORKDIR /app

COPY ./plugins/shared/ /shared

COPY ./plugins/req-resp-modifier/ /app

RUN go build -buildmode=plugin -o req-resp-modifier.so .

FROM devopsfaith/krakend:2.1.3

COPY krakend.json /etc/krakend/
COPY --from=httpClientBuilder /app/http-client.so /etc/krakend/plugins/.bin/
COPY --from=httpServerBuilder /app/http-server.so /etc/krakend/plugins/.bin/
COPY --from=reqRespModifierBuilder /app/req-resp-modifier.so /etc/krakend/plugins/.bin/

CMD ["krakend", "run", "-d", "-c", "/etc/krakend/krakend.json"]

EXPOSE 8080
