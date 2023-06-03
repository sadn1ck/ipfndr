# godns

absolutely terrible (incomplete as well) implementation of dns resolving, following [Julia Evan's tutorial](https://implement-dns.wizardzines.com/index.html)

## why lol

- Learning go
- Wanted to write some code
- Writing go
- DNS is fun
- Implementing real life stuff
- DNS is not fun
- Might as well use an array + offset innit

## todo

### general

- [ ] Add support for other record types
- [ ] CLI with flags to query specific types
- [ ] Better (nonexistent) error handling

### [technical](https://implement-dns.wizardzines.com/book/exercises.html)

- [ ] the dns domain compression/decompression doesn't check for pointer loops
- [ ] caching? lmao
