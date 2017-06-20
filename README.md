# measure

## usage

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mizkei/measure"
)

func A(ctx context.Context) {
	defer measure.Measure(ctx)()

	B(ctx)
	time.Sleep(2 * time.Second)
}

func B(ctx context.Context) {
	defer measure.Measure(ctx)()

	C(ctx)
	time.Sleep(3 * time.Second)
}

func C(ctx context.Context) {
	defer measure.Measure(ctx)()

	time.Sleep(4 * time.Second)
}

func main() {
	ctx := measure.ContextWithMeasure(context.Background())

	A(ctx)

	res := measure.GetResults(ctx)
	for _, r := range res {
		fmt.Printf("[%s] start:%s, end:%s\n", r.Func.Name(), r.Start.Format(time.Stamp), r.End.Format(time.Stamp))
	}
	// Output:
	// [main.C] start:Jun  5 12:34:00, end:Jun  5 12:34:04
	// [main.B] start:Jun  5 12:34:00, end:Jun  5 12:34:07
	// [main.A] start:Jun  5 12:34:00, end:Jun  5 12:34:09
}
```
