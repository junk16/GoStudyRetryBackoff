package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("vim-go")
	///	retrier := NewRetryExponentialBackoff(3)
	//retrier := NewRetryConstantBackoff(3, 1)
	retrier := NewRetryLinerBackoff(3, 1)
	err := Run(retrier)
	if err != nil {
		panic(err)
	}

	fmt.Println("SUCCESS!!!!!!!")
}

////////////////////////////////////////////////
type Retrier interface {
	Fail() Retrier
	Wait()
	IsRetryable() bool
}

////////////////////////////////////////////////
type RetryExponentialBackoff struct {
	Retrier
	pos     int
	retried int
}

func NewRetryExponentialBackoff(retry int) Retrier {
	return &RetryExponentialBackoff{
		pos:     retry,
		retried: 0,
	}
}

func (r *RetryExponentialBackoff) Fail() Retrier {
	r.pos--
	r.retried++
	return r
}

func (r *RetryExponentialBackoff) Wait() {
	wait := int(math.Pow(2, float64(r.retried)))
	fmt.Printf("Sleep : %d\n", wait)
	time.Sleep(time.Duration(wait) * time.Second)
}

func (r *RetryExponentialBackoff) IsRetryable() bool {
	if r.pos == 0 {
		return false
	}
	return true
}

////////////////////////////////////////////////
type RetryConstantBackoff struct {
	Retrier
	pos     int
	retried int
	wait    int
}

func NewRetryConstantBackoff(retry int, wait int) Retrier {
	return &RetryConstantBackoff{
		pos:     retry,
		retried: 0,
		wait:    wait,
	}
}

func (r *RetryConstantBackoff) Fail() Retrier {
	r.pos--
	r.retried++
	return r
}

func (r *RetryConstantBackoff) Wait() {
	fmt.Printf("Sleep : %d\n", r.wait)
	time.Sleep(time.Duration(r.wait) * time.Second)
}

func (r *RetryConstantBackoff) IsRetryable() bool {
	if r.pos == 0 {
		return false
	}
	return true
}

/////////////////////////////////////////////////
type RetryLinerBackoff struct {
	Retrier
	pos        int
	retried    int
	additional int
}

func NewRetryLinerBackoff(retry int, additional int) Retrier {
	return &RetryLinerBackoff{
		pos:        retry,
		retried:    0,
		additional: additional,
	}
}

func (r *RetryLinerBackoff) Fail() Retrier {
	r.pos--
	r.retried++
	return r
}

func (r *RetryLinerBackoff) Wait() {
	wait := r.additional + (r.retried - 1)
	fmt.Printf("Sleep : %d\n", wait)
	time.Sleep(time.Duration(wait) * time.Second)
}

func (r *RetryLinerBackoff) IsRetryable() bool {
	if r.pos == 0 {
		return false
	}
	return true
}

////////////////////////////////////////////////
func Run(r Retrier) error {
	fmt.Println(r)
	if !r.IsRetryable() {
		return fmt.Errorf("can not scuccess")
	}

	err := test()
	if err != nil {
		r.Fail().Wait()
		return Run(r)
	}
	err = test2()
	if err != nil {
		r.Fail().Wait()
		return Run(r)
	}
	err = test3()
	if err != nil {

		r.Fail().Wait()
		return Run(r)
	}
	return nil
}

func test() error {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(10)

	fmt.Printf("test if even num, then error :%d\n", i)

	if i%2 == 0 {
		return fmt.Errorf("Error")
	}

	return nil
}

func test2() error {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(30)

	fmt.Printf("test2 if less than 10, then error :%d\n", i)

	if i <= 10 {
		return fmt.Errorf("Error")
	}

	return nil
}

func test3() error {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(30)

	fmt.Printf("test3 if less than 5, then error :%d\n", i)

	if i <= 5 {
		return fmt.Errorf("Error")
	}

	return nil
}
