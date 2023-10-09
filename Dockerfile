FROM krakend/builder:2.1.3 as pluginsBuilder

WORKDIR /app

COPY ./plugins/ /app

RUN go build -buildmode=plugin -o plugins.so .

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
COPY --from=pluginsBuilder /app/plugins.so /etc/krakend/plugins/.bin/

RUN mkdir -p /etc/krakend/out && chmod -R 777 /etc/krakend/out

CMD ["krakend", "run", "-d", "-c", "/etc/krakend/krakend.tmpl"]

EXPOSE 8080
