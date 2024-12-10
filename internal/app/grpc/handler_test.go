package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/korol8484/shortener/internal/app/config"
	"github.com/korol8484/shortener/internal/app/grpc/service"
	"github.com/korol8484/shortener/internal/app/grpc/service/contract"
	"github.com/korol8484/shortener/internal/app/storage/memory"
	"github.com/korol8484/shortener/internal/app/usecase"
	"github.com/korol8484/shortener/internal/app/user/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

func server(userInt bool) (service.InternalClient, func()) {
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	uCase := usecase.NewUsecase(
		&config.App{BaseShortURL: "http://localhost"},
		memory.NewMemStore(),
		usecase.NewPingDummy(),
		zap.L(),
	)

	hand := NewHandler(uCase)

	var inters []grpc.UnaryServerInterceptor
	if userInt {
		inters = append(inters, JwtInterceptor(usecase.NewJwt(storage.NewMemoryStore(), zap.L(), "1234567891")))
	}
	inters = append(inters, IPInterceptor("127.0.0.1/24", []string{"/service.Internal/Stats"}))

	baseServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(inters...),
	)

	service.RegisterInternalServer(baseServer, hand)
	go func() {
		if err := baseServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	ccOpts := []grpc.DialOption{
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient("passthrough://bufnet", ccOpts...)
	if err != nil {
		log.Fatal(err)
	}

	closer := func() {
		listener.Close()
		baseServer.Stop()
		uCase.Close()
	}

	client := service.NewInternalClient(conn)

	return client, closer
}

func TestHandler_Ping(t *testing.T) {
	client, closer := server(true)
	defer closer()

	_, err := client.Ping(context.Background(), &empty.Empty{})
	require.NoError(t, err)
}

func TestHandler_HandleShort(t *testing.T) {
	client, closer := server(true)
	defer closer()

	res, err := client.HandleShort(context.Background(), &contract.RequestShort{
		Data: "http://ya.ru",
	})
	require.NoError(t, err)
	assert.Equal(t, res.GetData(), "http://localhost/zVFF0J")

	is, err := client.HandleGet(context.Background(), &contract.RequestFindByAlias{Alias: "zVFF0J"})
	require.NoError(t, err)
	assert.Equal(t, res.GetData(), is.GetData())

	_, err = client.HandleGet(context.Background(), &contract.RequestFindByAlias{Alias: "11"})
	require.Error(t, err)

	_, err = client.HandleShort(context.Background(), &contract.RequestShort{
		Data: "http__://ya.ru",
	})
	require.Error(t, err)
}

func TestHandler_HandleShortNotUser(t *testing.T) {
	client, closer := server(false)
	defer closer()

	_, err := client.HandleShort(context.Background(), &contract.RequestShort{
		Data: "http://ya.ru",
	})
	require.Error(t, err)

	_, err = client.HandleGet(context.Background(), &contract.RequestFindByAlias{Alias: "zVFF0J"})
	require.Error(t, err)

	_, err = client.HandleBatch(context.Background(), &contract.RequestBatch{Batch: []*contract.Batch{
		{
			CorrelationId: "1",
			OriginalUrl:   "http://ya.ru",
		},
	}})
	require.Error(t, err)

	_, err = client.UserURL(context.Background(), &empty.Empty{})
	require.Error(t, err)
}

func TestHandler_HandleBatch(t *testing.T) {
	client, closer := server(true)
	defer closer()

	batch, err := client.HandleBatch(context.Background(), &contract.RequestBatch{Batch: []*contract.Batch{
		{
			CorrelationId: "1",
			OriginalUrl:   "http://ya.ru",
		},
	}})
	require.NoError(t, err)
	assert.Len(t, batch.GetBatch(), 1)

	_, err = client.HandleBatch(context.Background(), &contract.RequestBatch{Batch: []*contract.Batch{
		{
			CorrelationId: "1",
			OriginalUrl:   "http__://ya.ru",
		},
	}})
	require.Error(t, err)
}
