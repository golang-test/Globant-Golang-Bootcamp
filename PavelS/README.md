# bookstore
RestAPI application on "Go" language

Setup "Bookstore".
Step 1. Database setup: configuration file "docker-compose.yml".
Command: docker-compose up
Step 2. Fill in database. 
Command: migrate -path migrations -database mysql://mysql:mysql@tcp(192.168.99.100)/bookstore_db?charset=utf8 up
Step 3. Build application images.
Command: docker build -t bookstore .     
Step 4. Start the application "Bookstore".
Command: docker run --name=go_bookstore -p 8080:8080 bookstore

