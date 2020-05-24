package dao

import "fmt"

func pipeline() (interface{}, error) {
	c.Send("SET", "foo", "bar")
	c.Send("GET", "foo")
	c.Flush()
	c.Receive()           // reply from SET
	v, err := c.Receive() // reply from GET
	return v, err
}

func pipelinedTransactions() {
	c.Send("MULTI")
	c.Send("INCR", "foo")
	c.Send("INCR", "bar")
	r, _ := c.Do("EXEC")
	fmt.Println(r) // prints [1, 1]
}
