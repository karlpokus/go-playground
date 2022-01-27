# go-playground
Let's recreate https://go.dev/play. The challenge will be to create a proper isolated execution environment.

What do we need?
- a server
- a container to run `go run <code>`
- at least two shared volumes. One r/w for the code and a read-only for bind-mounting go

How should it work?
- A browser gets some html at /
- user adds go code to a textarea and then sends it to the server at /code
- server puts the code on disk
- server runs `docker run ... <code-path>` and returns stdout to the browser
- browser displays the result

# todos
- [ ] scaling concerns
- [ ] bonus: let the user save the code

# licence
MIT
