# Walkernews

```
               _ _                                                 
              | | |                                                
__      ____ _| | | _____ _ __ _ __   _____      _____             
\ \ /\ / / _` | | |/ / _ \ '__| '_ \ / _ \ \ /\ / / __|            
 \ V  V / (_| | |   <  __/ |  | | | |  __/\ V  V /\__ \            
  \_/\_/ \__,_|_|_|\_\___|_|  |_| |_|\___| \_/\_/ |___/
  ~ A * HACKERNEWS * API * CLONE ~
```

Walkernews is a Hackernews API clone. The primary technologies this uses:

* [Go](https://golang.org) - an open source, strongly-typed language developed by Google
* [GraphQL](https://graphql.org) - An API paradigm that gives the client more control over their calls
* [FaunaDB](https://fauna.com) - A NoSQL, serverless database.

## Why though?

I started this project specifically to learn fundamentals of the Go programming language; which has grown quite a bit in
popularity since its release in 2010. If this wasn't a learning project, I probably would have opted for a serverless
configuration. Good news is: lots of serverless services support Go, and it's very performant!

In the course of this project, I also wound up learning quite a bit about FaunaDB. Fauna was much easier to get rolling
with in Go than DynamoDB. Its query language is much more expressive, the driver library is well-documented, and the
service as a whole is much friendlier than Dynamo! Fun fact - Fauna handles almost all
the [authentication](https://docs.fauna.com/fauna/current/tutorials/authentication/) for this project!

## Side note on Go

As a side note - Go as a whole, while overall a great tool, struggles to be approachable for beginners. Documentation is
less verbose, and more dependent on the developer reading the library's code themselves.

## Credit where credit is due

Though I deviated a lot, I did use this tutorial to get started building the
API: [Link](https://www.howtographql.com/graphql-go/0-introduction/)