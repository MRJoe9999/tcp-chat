# TCP server & client
## Instructions
1. Download as ZIP or clone
2. run the command ```go run server.go``` 
3. This is will start up the server.
4. Open another terminal and run ```go run client.go```
5. This will prompt you to enter a username
6. And then you can start chat
7. You can open another terminal and chat with each other

# Commands for testing
1. to delay messages ```sudo tc qdisc add dev lo root netem delay 100ms```
2. for packet loss ```sudo tc qdisc add dev lo root netem loss 30%```
# must remove the rule after done testing
1. here is the command to remove the rule ```sudo tc qdisc del dev lo root```

# Disconnect from the server
1. run this command on the terminal where the client is running
```\quit```


# link to video
https://youtu.be/K8WiN97yVhI

# link to presentation
https://www.canva.com/design/DAGnKL9B77Y/IJerSwG_jNwY-RvH-bOdwA/edit?utm_content=DAGnKL9B77Y&utm_campaign=designshare&utm_medium=link2&utm_source=sharebutton



