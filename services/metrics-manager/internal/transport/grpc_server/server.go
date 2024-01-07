package grpc_server

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"log"
	"metricsTerrarium/lib"
	"metricsTerrarium/services/metrics-manager/internal/database"
	"metricsTerrarium/services/metrics-manager/internal/general_types"
	"metricsTerrarium/services/metrics-manager/internal/transport/grpc_server/metrics"
	"metricsTerrarium/services/metrics-manager/pkg/api"
	"net"
)

type GrpcServerImpl struct {
	api.UnimplementedMetricsGetterServer

	dbConnection *sql.DB
	rawMetrics   map[string]general_types.RawMetric
}

type Server struct {
	service *GrpcServerImpl
	server  *grpc.Server
}

type ServerProperties struct {
	Config       *lib.Config
	DbConnection *database.Db
	RawMetrics   map[string]general_types.RawMetric
}

func (s *Server) Start(properties ServerProperties) {
	s.server = grpc.NewServer()
	s.service = &GrpcServerImpl{
		UnimplementedMetricsGetterServer: api.UnimplementedMetricsGetterServer{},
		dbConnection:                     properties.DbConnection.Connection,
		rawMetrics:                       properties.RawMetrics,
	}

	api.RegisterMetricsGetterServer(s.server, s.service)

	lis, err := net.Listen("tcp", properties.Config.MetricsManagerPort)
	if err != nil {
		log.Fatalf("TCP Connection creation error. Err: %s", err)
	} else {
		log.Printf("Created GRPC listener at %s", properties.Config.MetricsManagerPort)
	}

	err = s.server.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve. Err: %s", err)
	}
}

// GRPC handlers

func (s GrpcServerImpl) GetRawMetrics(context.Context, *api.RawMetricsRequestMessage) (*api.MetricsResponse, error) {
	return metrics.GetRawMetrics(s.rawMetrics)
}

func (s GrpcServerImpl) GetPreparedMetrics(ctx context.Context, req *api.PreparedMetricsRequestMessage) (*api.MetricsResponse, error) {
	return metrics.GetPreparedMetrics(req, s.dbConnection)
}
