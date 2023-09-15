# redis

## pub/sub usage


### Publish

如果發送資料給 redis channel 的時候可以使用 Publish 的功能

發布功能要實作：

1. service/repository/redis.go
2. service/delivery/delivery.go

service/repository/redis.go

```golang

const (
	FREE_TICKET_CHANNEL = "freeTicketChannel"
)

type RedisRepo struct {
	pool *redigo.Pool
}

func (r *RedisRepo) PublishFreeTicket(ctx context.Context, data interface{}) (err error) {
	err = redis.PUBLISH(r.pool, ctx, FREE_TICKET_CHANNEL, data)
	return
}

```

service/usecase/usecase.go

```golang

func (u *Usecase) updateFreeTicket() {
    // do something
    u.redisRepo.(ctx.Context, "message...")
}

```

### Subscribe

如果想從 redis 接收 pub/sub 的訊息的時候可以使用訂閱的功能

訂閱功能要實作：

1. service/repository/redis.go
2. service/delivery/delivery.go

service/repository/redis.go

```golang

const (
	FREE_TICKET_CHANNEL = "freeTicketChannel"
    LOGIN_CHANNEL = "loginChannel"
)

type RedisRepo struct {
	pool *redigo.Pool
}

func (r *RedisRepo) Subscribe(onMessage func(redigo.PubSubConn), channel ...interface{}) {
	redis.Subscribe(r.pool, onMessage, channel...)
}
```

service/delivery/delivery.go

```golang
type Delivery struct {
	redisRepo *repository.RedisRepo
	useCase   *usecase.UseCase
	channel   map[string]func(data []byte)
}

func New(
	redisRepo *repository.RedisRepo,
	useCase *usecase.UseCase,
) (d *Delivery) {
	d = &Delivery{
		redisRepo: redisRepo,
		useCase:   useCase,
	}

	d.channel = map[string]func(data []byte){
		repository.FREE_TICKET_CHANNEL: d.UpdatePlayerFreeTicket,
		repository.LOGIN_CHANNEL:       d.RecordPlayerLogin,
	}

	return
}

func (d *Delivery) Start() {
    var channel []interface{}
    for channelName :=range d.channel {
        channel = append(channel, channelName)
    }

	go d.redisRepo.Subscribe(d.onMessage, channel...)
}

func (d *Delivery) onMessage(pubSubConn redigo.PubSubConn) {
	for {
		switch v := pubSubConn.Receive().(type) {
		case redigo.Message:
            channelCallback, ok := d.channel[v.Channel]
            if !ok {
				zap.L().Warn("undefined channel")
                continue
            }
            
            channelCallback(v.Data)
			// fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
		case redigo.Subscription:
			// fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			zap.L().Error(v.Error())
			return
		}
	}
}

func (d *Delivery) UpdatePlayerFreeTicket(msg []byte) {
    ctx := ctx.Context{Context: context.Background()}
    data := UpdatePlayerFreeTicketEvent{}
    json.Unmarshal(msg, &data)

    d.useCase.UpdatePlayerFreeTicket(ctx, data)
}

func (d *Delivery) RecordPlayerLogin(msg []byte) {
    ctx := ctx.Context{Context: context.Background()}
    data := RecordPlayerLoginEvent{}
    json.Unmarshal(msg, &data)

    d.useCase.RecordPlayerLogin(ctx, data)
}

```
