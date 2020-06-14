# shorturl

# Test
Start redis server using `make start-redis-server`
Start server and then send a put request to create a short URL
Example : If you have cURL installed
```
curl -X PUT -H "Content-Type: application/json" -d '{"url":"www.google.com"}' localhost:5000/create
```
