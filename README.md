# menext-backend
backend for me-next

party - contains code for the party object. This includes various queues (playnext, suggest) as well as currently playing information. All changes to underlying data structures happen through the party class. 

server - contains code to manage interactions between clients and their parties via an http server. A key class in this file is the PartyManager which routes requests to the relevant party. The server functions are distributed across several files, with server.go containing administrative functions, nowPlayingAPI.go containing playing functions, and queueAPI.go containing queue functions. 



