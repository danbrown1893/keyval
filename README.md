#Cache Library

##Main
The main file can be used to start the cache and a server that will allow a user to make \
requests for storing, retrieving, and deleting items from the cache.

The deadline is expected to be a Unix timestamp for this implementation. It will also accept \
and empty string in which case the cached value will not expire. \
The endpoint will be `localhost:3000/store` \
A sample POST request is:
```
{
    "key": "foo",
    "value": "value",
    "deadline": "1739331922"
}
```

A sample GET is:\
`localhost:3000/store?foo`

A sample DELETE is:\
`localhost:3000/store?foo`

##Cache
The cache uses a mutex which will allow it to be used in a process with multiple threads.\
There are also tests and benchmarks which can be used to check the performance targets.\
The `Get` method benchmark for a custom data structure stored as the value on my machine shows:
```
goos: darwin
goarch: amd64
pkg: keyval/cache
cpu: Intel(R) Core(TM) i7-7820HQ CPU @ 2.90GHz
BenchmarkCache_GetCustomData
BenchmarkCache_GetCustomData-8   	10365808	       106.5 ns/op
```

The `Set` method benchmark on my machine shows:
```
goos: darwin
goarch: amd64
pkg: keyval/cache
cpu: Intel(R) Core(TM) i7-7820HQ CPU @ 2.90GHz
BenchmarkCache_SetCustomData
BenchmarkCache_SetCustomData-8   	 1806247	       779.3 ns/op
```

There are comments on each public method in the cache library that describe the\
functionality.

The 10,000,000 Key/Value pairs performance target will depend on the size of the\
values that are stored as well. As long as the machine/pod/instance has enough memory\
it will work and will just increase the cost.
