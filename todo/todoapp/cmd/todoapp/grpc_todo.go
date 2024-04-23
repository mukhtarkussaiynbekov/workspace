package main

import (
	"context"

	"github.com/mukhtarkv/workspace/api/errdetails"
	pb "github.com/mukhtarkv/workspace/api/todo/todoapp/v1beta1"
	"github.com/mukhtarkv/workspace/kit/errors"
	"github.com/mukhtarkv/workspace/kit/log"
	"github.com/mukhtarkv/workspace/todo/todoapp"
	"google.golang.org/grpc/codes"
)

// grpcToDo represent the grpc todo item service implementation.
type grpcToDo struct {
	pb.UnsafeToDoAppServer
	service *todoapp.ToDoService
	log     *log.Logger
}

func newGrpcToDo(service *todoapp.ToDoService) (*grpcToDo, error) {

	if service == nil {
		return nil, errors.New("nil todo item service")
	}

	return &grpcToDo{
		service: service,
		log:     log.L(),
	}, nil
}

func (u *grpcToDo) List(ctx context.Context, request *pb.ListRequest) (*pb.ListResponse, error) {

	items, err := u.service.List(ctx)
	if err != nil {
		u.log.Error(ctx, "listing todo", log.Error(err))

		return nil, errors.Status(
			codes.Unknown,
			err.Error(), // log the full error to provide enough context
			&errdetails.ErrorInfo{
				Reason: "UNKNOWN_ERROR",
			})
	}
	res := []*pb.ListResponse_ToDoItem{}
	for _, item := range items {
		res = append(res, &pb.ListResponse_ToDoItem{
			Id: item.Id,
			Title: item.Title,
			Details: item.Details,
		})
	}
	return &pb.ListResponse{
		TodoItems: res,
	}, nil
}

func (u *grpcToDo) Create(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {

	if err := request.Validate(); err != nil {
		u.log.Error(ctx, "creating todo item", log.Error(err))

		return nil, errors.Status(
			codes.InvalidArgument,
			err.Error(),
			&errdetails.ErrorInfo{
				Reason: "INVALID_REQUEST",
				Metadata: map[string]string{
					"request": request.String(),
				},
			})
	}

	todo := todoapp.ToDoItem{
		Title: request.Title,
		Details: request.Details,
	}
	if err := u.service.Create(ctx, &todo); err != nil {
		u.log.Info(ctx, "creating todo item", log.Error(err))

		if errors.Is(err, todoapp.ErrToDoItemAlreadyExist) {
			return nil, errors.Status(
				codes.AlreadyExists,
				"todo item already exists in the system",
				&errdetails.ErrorInfo{
					Reason: "INVALID_REQUEST",
					Metadata: map[string]string{
						"request": request.String(),
					},
				})
		}

		return nil, errors.Status(
			codes.InvalidArgument,
			err.Error(),
			&errdetails.ErrorInfo{
				Reason: "FAIL_CREATE_TODO",
				Metadata: map[string]string{
					"request": request.String(),
				},
			})
	}

	return &pb.CreateResponse{
		Id:   todo.Id,
		Title: todo.Title,
		Details: todo.Details,
	}, nil
}

func (u *grpcToDo) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	if err := request.Validate(); err != nil {
		u.log.Error(ctx, "updating todo item", log.Error(err))

		return nil, errors.Status(
			codes.InvalidArgument,
			err.Error(),
			&errdetails.ErrorInfo{
				Reason: "INVALID_REQUEST",
				Metadata: map[string]string{
					"request": request.String(),
				},
			})
	}

	todo := todoapp.ToDoItem{
		Id: request.Id,
		Title: request.Item.Title,
		Details: request.Item.Details,
	}
	if err := u.service.Update(ctx, &todo, request.UpdateMask.Paths); err != nil {
		u.log.Info(ctx, "updating todo item", log.Error(err))

		if errors.Is(err, todoapp.ErrToDoItemNotFound) {
			return nil, errors.Status(
				codes.AlreadyExists,
				"todo item not found",
				&errdetails.ErrorInfo{
					Reason: "INVALID_REQUEST",
					Metadata: map[string]string{
						"request": request.String(),
					},
				})
		}

		return nil, errors.Status(
			codes.InvalidArgument,
			err.Error(),
			&errdetails.ErrorInfo{
				Reason: "FAIL_UPDATE_TODO",
				Metadata: map[string]string{
					"request": request.String(),
				},
			})
	}

	return &pb.UpdateResponse{}, nil
}

func (u *grpcToDo) Delete(ctx context.Context, request *pb.DeleteRequest) (*pb.DeleteResponse, error) {

	if err := u.service.Delete(ctx, request.Id); err != nil {
		u.log.Error(ctx, "deleting todo item", log.Error(err), log.String("todo.id", request.Id))

		return nil, errors.Status(
			codes.InvalidArgument,
			err.Error(),
			&errdetails.ErrorInfo{
				Reason: "FAIL_DELETE_TODO",
				Metadata: map[string]string{
					"request": request.String(),
				},
			})
	}
	return &pb.DeleteResponse{}, nil
}
