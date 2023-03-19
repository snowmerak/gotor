# gotor

gotor is a simple tool to help you to generate actor model in golang.

## Installation

```bash
go install github.com/snowmerak/gotor
```

## Usage

```bash
# help
gotor -h

# generate template yaml file
gotor -init <file-path.yaml>

# generate actor from yaml file
gotor -gen <file-path.yaml>
```

## Example

```yaml
actors:
  - path: test
    package_name: ActorMap
    actor_name: Map
    channels:
      - name: get
        type: tuple.Tuple[string, chan string]
      - name: set
        type: string
      - name: delete
        type: string
```

```go
package ActorMap

import "context"

type Map struct {
	ctx      context.Context
	cancel   context.CancelFunc
	setCh    chan string
	deleteCh chan string
	getChCh  chan tuple.Tuple[string, chan string]
	// TODO: Write your actor states here
}

func NewMap(ctx context.Context, queueSize int) *Map {
	ctx, cancel := context.WithCancel(ctx)
	setCh := make(chan string, queueSize)
	deleteCh := make(chan string, queueSize)
	getChCh := make(chan tuple.Tuple[string, chan string], queueSize)

	go func() {
		for {
			select {
			case <-ctx.Done():
				// TODO: Write your actor stop logic here
				return
			case value := <-getChCh:
			// TODO: Write your actor logic here
			case value := <-setCh:
			// TODO: Write your actor logic here
			case value := <-deleteCh:
				// TODO: Write your actor logic here
			}
		}
	}()
	return &Map{
		ctx:      ctx,
		cancel:   cancel,
		setCh:    setCh,
		deleteCh: deleteCh,
		getChCh:  getChCh,
	}
}

func (m *Map) DeleteCh(value string) {
	m.deleteCh <- value
}

func (m *Map) GetChCh(value tuple.Tuple[string, chan string]) {
	m.getChCh <- value
}

func (m *Map) SetCh(value string) {
	m.setCh <- value
}

func (m *Map) Stop() {
	m.cancel()
}
```
