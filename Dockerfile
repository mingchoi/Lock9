FROM plugins/base:multiarch

LABEL maintainer="Ming Choi <mingchoi.na@gmail.com>" \
	org.label-schema.name="Lock9" \
	org.label-schema.vendor="Ming Choi" \
	org.label-schema.schema="1.0"

ADD release/linux/amd64/lock9 /bin/

ENTRYPOINT ["/bin/lock9"]
