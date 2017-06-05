package measure

import (
	"context"
	"runtime"
	"sync"
	"time"
)

type ctxKey string

const key ctxKey = "measure_key"

type Result struct {
	Start, End time.Time
	Func       *runtime.Func
	File       string
	Line       int
}

type Results []Result

type memo struct {
	sync.Mutex
	res *Results
}

func doNothing() {}

func ContextWithMeasure(ctx context.Context) context.Context {
	v := make(Results, 0, 100)
	return context.WithValue(ctx, key, &memo{res: &v})
}

func GetResults(ctx context.Context) Results {
	res, ok := ctx.Value(key).(*memo)
	if !ok || res == nil {
		return nil
	}
	return *res.res
}

func Measure(ctx context.Context) func() {
	v, ok := ctx.Value(key).(*memo)
	if !ok || v == nil {
		return doNothing
	}

	var res Result
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return doNothing
	}
	res.Func = runtime.FuncForPC(pc)
	res.File, res.Line = file, line
	res.Start = time.Now()

	return func() {
		res.End = time.Now()
		v.Lock()
		*v.res = append(*v.res, res)
		v.Unlock()
	}
}
