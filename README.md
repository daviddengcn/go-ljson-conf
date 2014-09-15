go-ljson-conf [![GoSearch](http://go-search.org/badge?id=github.com%2Fdaviddengcn%2Fgo-ljson-conf)](http://go-search.org/view?id=github.com%2Fdaviddengcn%2Fgo-ljson-conf)
=============

A powerful configuration package for go using Loose JSON as the format.
([godoc](http://godoc.org/github.com/daviddengcn/go-ljson-conf))

Features
--------
**Data format**

The data format is the Loose JSON, which is full compatible with JSON format but more easier for writing by hand.
Visit [Loose JSON project](https://github.com/daviddengcn/ljson) for more details.

Lines starting with `//` or `;` (excluding the leading spaces) are considered as comments and thus ignored.

**Dot-seperated hierarchical key**

For configure:

```javascript
{
	// http settings
	http: {
		addr: "www.example.com"
		ports: [80, 8080]
	}
}
```
You can fetch values by different ways:

code                   |value              |type
-----------------------|-------------------|---------
`String("http.addr")`  |`"www.example.com"`|`string`
`IntList("http.ports")`|`[80, 8080]`       |`[]int`
`Int("http.ports[1]")` |`8080`             |`int`
`Object("http")`       |`map[addr:www...]` |`map[string]interface{}`

**Include function**

`main.conf`

```javascript
{
	http: {
		#include#: "addr.conf"
	}
}
```

`addr.conf`

```javascript
{
	addr: "www.example.com"
	ports: [80, 8080]
}
```

we got:

```javascript
{
	http: {
		addr: "www.example.com"
		ports: [80, 8080]
	}
}
```

**Multi-path searching**
We calling `Load` function, it searches the following path in order:
1. For absolute path, it is directly used,
1. Current directory,
1. Directory of the executable, and
1. User's home directory

LICENSE
-------
BSD license.
