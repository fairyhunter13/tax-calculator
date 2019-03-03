FROM nginx:alpine


COPY ./scripts/nginxwrapper.sh .

RUN chmod +x nginxwrapper.sh ; \    
    apk update; \
    apk add --no-cache bash curl; \
    rm -rf /var/cache/apk/*; 

CMD ./nginxwrapper.sh