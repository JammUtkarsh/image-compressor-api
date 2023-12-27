# Zocket SDE Assignment

Use [Postman](https://www.postman.com/jammutkarsh/workspace/jammutkarsh-apis/request/29627550-7888fc91-1d04-41b8-9550-a2f750a5c8c3?ctx=documentation) or Curl to add sample data in the DB:

```bash
curl --location --request POST '127.0.0.1:8080/' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": 1,
    "product_name": "Intelligent Fresh Chips",
    "product_description": "Corporate",
    "product_price": 320.89,
    "product_images": [
        "http://placeimg.com/640/480",
        "http://placeimg.com/640/480",
        "http://placeimg.com/640/480",
        "http://placeimg.com/640/480",
        "http://placeimg.com/640/480"
    ]
}'
```
