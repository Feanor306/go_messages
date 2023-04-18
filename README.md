Golang Messages queue using Gin, RabbitMQ and Redis

Application 1: API 
POST Endpoint: /message
POST body: { sender: String, receiver: String, message: String }
-> pushes received information to a RabbitMQ Queue
Return OK Status if everything is there, otherwise Bad Request

Application 2: MessageProcessor
Subscribes to queue from RabbitMQ and processes the message
Processing of message means saving the message to Redis in a way that application 3
works.

Application 3: Reporting API
GET Endpoint: /message/list
Parameters: sender: String, receiver String
Returns an array of objects with sender, receiver and message content that were
exchanged between sender and receiver in chronological descending order
