# README.md

Side project to get hands on with Go.

- Using Redis (free tier, gives you 30MB)

  - You could store locally using a map.
  - Using Redis here. Redis stores key-value pairs in RAM(of the remote Redis server). It is essentially a distributed cache(not db). The Redis instance is managed by Redis cloud typcially hosted by AWS, Google Cloud, etc. But that means that its not persistent. Shift to something like Supabase or MongoDB if you want persistence.

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

# Usefule ref

[this](https://codegangsta.gitbooks.io/building-web-apps-with-go/content/)
