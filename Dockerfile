FROM alpine:3.3
ADD bin/chef-synk-linux /usr/bin/chef-synk
RUN chmod +x /usr/bin/chef-synk \
  && apk add --update -t deps ca-certificates \
  && apk del --purge deps \
  && rm /var/cache/apk/*
CMD "chef-synk"
