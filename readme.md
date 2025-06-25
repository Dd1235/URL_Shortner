# README.md

Side project to get hands on with Go.

- Using Redis (free tier, gives you 30MB)

  - You could store locally using a map.
  - Using Redis here. Redis stores key-value pairs in RAM(of the remote Redis server). It is essentially a distributed cache(we are using it as a db). The Redis instance is managed by Redis cloud typcially hosted by AWS, Google Cloud, etc. But that means that its not persistent. Shift to something like Supabase or MongoDB if you want persistence.

- Will implement more sophisticated shortening algorithm later
- Deciding to make new short url even if key already exist. This is for better analytics.
- will make mcp server for this later.

# Some Go Syntax (mostly for myself!)

- Go is designed with composition over inheritance in mind. If a type has a method(corresponding to the interface), it implements that interface without explicitly declaring it.

```Go

type Reader interface {
  Read(p []byte) (n int, err error)
}

// now any type that has a Read methods is said to implement the Reader interface
// structural typing, if it fits the shape, it qualifies (why does this remind me of Ducktyping? It is walk and looks like that interface, it implements that interface? I do not know if this is right though)
```

- There is no method overloading in Go.
- Every interface method is virtual

- `time.Now().UnixNano()` returns current time in nanoseconds since Jan 1, 1970 UTC. Right now if there are about 19 digits in the number, it will take a few more years to read 20 digits.

## Pointer vs Value Receiver

```
type User struct {
  Name string
}

func (u User) valChangeName(n string){
  u.Name = n// this won't change anything as this makes a copy of the object it is called on
}

func (u *User) varChangeName(n string){
  u.Name = n// this will change the original object as it is a pointer receiver
}

// now for whatever reason you use a value receiver, it makes a copy of the object.
// you need to use pointers to call pointer receiver methods though, objects can call both (automatically deref)

```

## Context in go

- how do you abort a long running db query if the client disconnected? how you time out a nested operation after 5 seconds? pass some requestID down multiple layers of logic? context solves these for you.

```Go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}           // closed when cancelled or timed out
    Err() error                      // why context was cancelled
    Value(key any) any               // get stored value
}
```

`defer` keyword delays the execution of a function until the surrounding function returns.

Consider

```python
@app.route('/fetch')
def fetch_data():
    result = slow_database_query()
    return jsonify(result)
```

```javascript
app.get("/fetch", async (req, res) => {
  const data = await slowFetch();
  res.json(data);
});
```

now the user closes the tab because we have low low attention spans in our generation.
How to detect this cancellation?

What we need is propogation of cancellation.

```Go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context() // automatically canceled if client disconnects

    data, err := fetchData(ctx)
    if err != nil {
        if ctx.Err() == context.Canceled {
            log.Println("Request cancelled by client")
        }
        http.Error(w, "timeout or cancelled", http.StatusRequestTimeout)
        return
    }

    json.NewEncoder(w).Encode(data)
}

func fetchData(ctx context.Context) (Data, error) {
    select {
    case <-ctx.Done():         // cancel immediately if context ends
        return nil, ctx.Err()
    case result := <-doQuery(): // simulate query
        return result, nil
    }
}
```

Notice the `select`? Looks like switch, but this is for concurrent communication using channels. It lets you wait on multiple channels and react to whichever is ready first.

```Go
ch := make(chan string) // this is not a string but a "channel" that carries strings

go func() {
    time.Sleep(2 * time.Second) // wait two seconds
    ch <- "done" // send message to channel
}()

select {
case msg := <-ch:
    fmt.Println("Received:", msg)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout!")
}
```

- if multiple channels are ready, randomly pick to avoid starvation.
- a channel is a typed conduit for sending and receiving values between goroutines.
- it is blocking, typed, and acts like a pipe
- '<-' is used to receive from a channel, 'chan <-' is used to send to a channel.

Wow, the rabbit hole goes on.

- goroutines are managed by the go runtime, not by the os
- started via `go func() { ... }()` which runs the function concurrently
- they are multiplexed onto a smaller pool of threads
- m:n scheduler, many many goroutines, few threads
- go scheduler manages preemptive scheduling, stack resizing, blocking detection
- so context switch cost is lower for goroutines
- added benefit of cooperate + preemptive scheduling as opposed to just preemptive scheduling with processes
- so `chan` makes the communication easier, than typical message queue, pipes, fifo, shared memory etc.

- JS is single threaded, non-premetive, thats why the whole async/await thing.
- C++ and Java prelooms use os threads
- Java loop has fibers much like the goroutines.
- gives more context to the only one main.go with `func main()` with a package main thing, go runtime is to manage one scheduler with one go routine pool nad one stack, makes compiling simpler too

- so go does it best

# Dynamo db

- Managed NoSQL database
- cloud native, serverless, key-value and document store
- elastic scaling
- useful when you know your access patterns
- access through api/orm authorized through IAM
- integrates well with other aws services
- It is one of the most popular nosql database out there, provides consistent performance at any skills
- 2021 prime day sales, trillions of calls to dynamodb, high availability iwth digit millisecond performance
- consistent performance at any scale. (low single-digit millisecond latency)
- Workload pattern
  - multi tenant - load of one customer should not affect another
  - high resource utilization
  - boundless scale of tables
  - preditable performance
  - highly available (replication and recovery)
  - flexible usecase support

## Architecture

- multiple table, each is a collection of items
- each item - primary key
- primary key needs to be specified
- sort key + partition key(required) - determines where and how the data is stored
- item -> f(k) - determines which partition the data goes into. Hash based, range based routing etc.
- Also supports secondary indexes. Indexed value -> PK

# AWS Lambda

- Traditional company specific data centers were expensive, hard to scale, and required a lot of maintenance.
- EC2 - elastic compute cloud. Managed service. Use amazons servers insteads.
- Enhanced cloud infrastructure - we pay per execution in lambda, not for the hardware, so based on traffic you pay.
- Pay per use cloud computing
- Run code at scale without worrying about servers
- write functions - primary unit of lambda.

workflow:

- create a function
- write and upload your code
- run your function
- no servers to manage
- patching and security updates are handled by AWS
- autoscaling
- netflix, udemy, lyft, etc. use lambda

# API Gateway

- parameter validation
- allow-list/deny-list (and some rate limit checks)
- authentication/authorization(identity provider)
- high level rate limit
- dynamic routing
- service discovery
- protocol conversation
- should be deployed in multiple regions

# Useful ref

[this](https://codegangsta.gitbooks.io/building-web-apps-with-go/content/)
[wow, beautiful](https://go-proverbs.github.io)
