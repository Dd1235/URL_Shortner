# README.md

Side project to get hands on with Go.

- Using Redis (free tier, gives you 30MB)

  - You could store locally using a map.
  - Using Redis here. Redis stores key-value pairs in RAM(of the remote Redis server). It is essentially a distributed cache(not db). The Redis instance is managed by Redis cloud typcially hosted by AWS, Google Cloud, etc. But that means that its not persistent. Shift to something like Supabase or MongoDB if you want persistence.
