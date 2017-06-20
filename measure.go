// Package measure provide context timer.
package measure

import (
	"context"
	"runtime"
	"sync"
	"time"
)

type ctxKey string

const key ctxKey = "measure_key"

// StopFunc stop measurement.
type StopFunc func()

// Result has measurement result.
type Result struct {
	Start, End time.Time
	// Func, File, and Line are target function info
	Func *runtime.Func
	File string
	Line int
}

// Results is array of Reuslt.
type Results []Result

type memo struct {
	sync.Mutex
	res *Results
}

func doNothing() {}

// ContextWithMeasure returns copy of the parent context with measure value.
func ContextWithMeasure(ctx context.Context) context.Context {
	v := make(Results, 0, 100)
	return context.WithValue(ctx, key, &memo{res: &v})
}

// GetResults returns measurement result.
func GetResults(ctx context.Context) Results {
	res, ok := ctx.Value(key).(*memo)
	if !ok || res == nil {
		return nil
	}
	return *res.res
}

// Measure start measurement.
// Measurement stops when the returned StopFunc is called.
func Measure(ctx context.Context) StopFunc {
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
