package grpcclient

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/viebiz/lit/grpcclient/testdata"
)

type streamWeatherServer struct {
	testdata.UnimplementedWeatherServiceServer
}

func (s *streamWeatherServer) StreamWeather(req *testdata.WeatherRequest, stream grpc.ServerStreamingServer[testdata.WeatherDetail]) error {
	details := []*testdata.WeatherDetail{
		{
			Location:    "Hive City, Necromunda",
			Date:        req.GetDate(),
			Description: "Toxic smog with occasional acid rain",
			Temperature: 42.7,
		},
		{
			Location:    "Macragge's Northern Hemisphere",
			Date:        "M41.874.21",
			Description: "Freezing winds with snowstorms",
			Temperature: -20.5,
		},
	}

	for _, d := range details {
		if err := stream.Send(d); err != nil {
			return err
		}
	}
	return nil
}

func TestClientConn_NewStream_ServerStreaming(t *testing.T) {
	// Given: start a gRPC server with a server-streaming method
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := lis.Addr().String()

	grpcServer := grpc.NewServer()
	testdata.RegisterWeatherServiceServer(grpcServer, &streamWeatherServer{})

	done := make(chan struct{})
	go func() {
		_ = grpcServer.Serve(lis)
		close(done)
	}()
	t.Cleanup(func() {
		grpcServer.Stop()
		<-done
	})

	// And: create an unauthenticated client connection (our wrapper type)
	conn, err := NewUnauthenticatedConnection(context.Background(), addr)
	require.NoError(t, err)

	// When: call a server streaming RPC which uses NewStream under the hood
	client := testdata.NewWeatherServiceClient(conn)
	stream, err := client.StreamWeather(context.Background(), &testdata.WeatherRequest{Date: "M41.993.32"})
	require.NoError(t, err)

	var got []*testdata.WeatherDetail
	for {
		d, rerr := stream.Recv()
		if rerr == io.EOF {
			break
		}
		require.NoError(t, rerr)
		got = append(got, d)
	}

	// Then: we received the expected stream of messages
	require.Len(t, got, 2)
	require.Equal(t, "Hive City, Necromunda", got[0].GetLocation())
	require.Equal(t, "M41.993.32", got[0].GetDate())
	require.Equal(t, "Macragge's Northern Hemisphere", got[1].GetLocation())
	require.Equal(t, "M41.874.21", got[1].GetDate())
}
