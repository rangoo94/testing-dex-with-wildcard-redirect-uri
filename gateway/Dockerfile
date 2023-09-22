FROM nginx:alpine

COPY ./index.html /usr/share/nginx/html/index.tpl

CMD [ \
  "sh", "-c", \
  "envsubst '$WHITELIST_REGEX' < /usr/share/nginx/html/index.tpl > /usr/share/nginx/html/index.html && nginx -g \"daemon off;\"" \
]
