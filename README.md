go-ljson-conf
=============

A powerful configuration package for go using Loose JSON as the format.
([godoc](http://godoc.org/github.com/daviddengcn/go-ljson-conf))

Features
--------
**Loose JSON format**

Full compatible with JSON format but more easier for writing by hand.
Visit [Loose JSON project](https://github.com/daviddengcn/ljson) for more details.

**Dot-seperated hierarchical key**

For configure:

```javascript
{
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

LICENSE
-------
BSD license.
