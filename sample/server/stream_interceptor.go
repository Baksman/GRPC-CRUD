package main

// import (
// 	// "context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/peer"
// )

// func StreamLoggingInterceptor(
// 	srv interface{},
// 	ss grpc.ServerStream,
// 	info *grpc.StreamServerInfo,
// 	handler grpc.StreamHandler,
// ) error {
// 	// Extract client info (IP address)
// 	p, ok := peer.FromContext(ss.Context())
// 	if ok {
// 		log.Printf("Client IP: %s - Started streaming: %s", p.Addr, info.FullMethod)
// 	} else {
// 		log.Printf("Started streaming: %s", info.FullMethod)
// 	}

// 	// Call the handler to proceed with the normal processing of the streaming RPC
// 	err := handler(srv, ss)

// 	// Log the end of the streaming call
// 	log.Printf("Finished streaming: %s", info.FullMethod)

// 	return err
// }

// // func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
// // 	o := evaluateOptions(opts)
// // 	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
// // 		panicked := true

// // 		defer func() {
// // 			if r := recover(); r != nil || panicked {
// // 				err = recoverFrom(stream.Context(), r, o.recoveryHandlerFunc)
// // 			}
// // 		}()

// // 		err = handler(srv, stream)
// // 		panicked = false
// // 		return err
// // 	}
// // }

// func LoggingStreamInterceptorText(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
// 	log.Printf("Incoming stream for %s", stream.Context().Value("grpc.service"))

// 	// wrapped := &loggingServerStream{ss}

// 	ctx := stream.Context()
// 	fmt.Println(info.FullMethod)

// 	newStream := &customServerStream{
// 		ServerStream: stream,
// 		ctx:          ctx,
// 	}
// 	err := handler(srv, newStream)
// 	log.Printf("done with request %s", stream.Context().Value("grpc.service"))
// 	return err
// }

// // func(srv any, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error
// func LoggingStreamInterceptor(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
// 	log.Printf("Incoming stream for %s", stream.Context().Value("grpc.service"))

// 	// wrapped := &loggingServerStream{ss}

// 	ctx := stream.Context()
// 	fmt.Println(info.FullMethod)

// 	newStream := &customServerStream{
// 		ServerStream: stream,
// 		ctx:          ctx,
// 	}
// 	err := handler(srv, newStream)
// 	log.Printf("done with request %s", stream.Context().Value("grpc.service"))
// 	return err

// }

// // type customServerStream struct {
// // 	grpc.ServerStream
// // 	ctx context.Context
// // }
// type loggingServerStream struct {
// 	grpc.ServerStream
// }

// func (s *loggingServerStream) Send(m interface{}) error {
// 	start := time.Now()
// 	err := s.ServerStream.SendMsg(m)
// 	log.Printf("Sent message in %v", time.Since(start))
// 	return err
// }

// func (s *loggingServerStream) RecvMsg(m interface{}) error {
// 	start := time.Now()
// 	err := s.ServerStream.RecvMsg(m)
// 	log.Printf("Received message in %v", time.Since(start))
// 	return err
// }
