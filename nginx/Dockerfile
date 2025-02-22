FROM nginx:1.27-alpine
EXPOSE 80

COPY backend.conf /etc/nginx/nginx.conf

CMD ["nginx"]
HEALTHCHECK --interval=5s --timeout=5s --retries=5 CMD curl --fail http://localhost/ping