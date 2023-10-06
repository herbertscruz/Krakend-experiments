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

ENV FC_ENABLE=1
ENV FC_PARTIALS="/etc/krakend/partials"
ENV FC_SETTINGS="/etc/krakend/settings"
ENV FC_TEMPLATES="/etc/krakend/templates"
ENV FC_OUT="/etc/krakend/out/krakend.json"

COPY krakend.tmpl /etc/krakend/
COPY partials /etc/krakend/partials
COPY settings /etc/krakend/settings
COPY templates /etc/krakend/templates
COPY --from=httpClientBuilder /app/http-client.so /etc/krakend/plugins/.bin/
COPY --from=httpServerBuilder /app/http-server.so /etc/krakend/plugins/.bin/
COPY --from=reqRespModifierBuilder /app/req-resp-modifier.so /etc/krakend/plugins/.bin/

RUN mkdir -p /etc/krakend/out && chmod -R 777 /etc/krakend/out

CMD ["krakend", "run", "-d", "-c", "/etc/krakend/krakend.tmpl"]

EXPOSE 8080
