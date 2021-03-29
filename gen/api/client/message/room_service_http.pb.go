// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             (unknown)
// source: message/room_service.proto

package climsg

import (
	context "context"
	http "github.com/go-kratos/kratos/transport/http"
	binding "github.com/go-kratos/kratos/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationRoomServiceAgreeToInviteJoinRoom = "/message.RoomService/AgreeToInviteJoinRoom"
const OperationRoomServiceApproveRequestToJoinRoom = "/message.RoomService/ApproveRequestToJoinRoom"
const OperationRoomServiceCloseRoom = "/message.RoomService/CloseRoom"
const OperationRoomServiceCreateRoom = "/message.RoomService/CreateRoom"
const OperationRoomServiceInviteToJoinRoom = "/message.RoomService/InviteToJoinRoom"
const OperationRoomServiceKickUserFromRoom = "/message.RoomService/KickUserFromRoom"
const OperationRoomServiceLeaveRoom = "/message.RoomService/LeaveRoom"
const OperationRoomServiceRequestToJoinRoom = "/message.RoomService/RequestToJoinRoom"
const OperationRoomServiceRoomDetail = "/message.RoomService/RoomDetail"
const OperationRoomServiceRoomList = "/message.RoomService/RoomList"

type RoomServiceHTTPServer interface {
	// AgreeToInviteJoinRoom Agree to invite to join room
	AgreeToInviteJoinRoom(context.Context, *CSAgreeToInviteJoinRoom) (*SCAgreeToInviteJoinRoom, error)
	// ApproveRequestToJoinRoom Approve request to join room
	ApproveRequestToJoinRoom(context.Context, *CSApproveRequestToJoinRoom) (*SCApproveRequestToJoinRoom, error)
	// CloseRoom Close room
	CloseRoom(context.Context, *CSCloseRoom) (*SCCloseRoom, error)
	// CreateRoom Create room
	CreateRoom(context.Context, *CSCreateRoom) (*SCCreateRoom, error)
	// InviteToJoinRoom Invite to join room
	InviteToJoinRoom(context.Context, *CSInviteToJoinRoom) (*SCInviteToJoinRoom, error)
	// KickUserFromRoom Kick user from room
	KickUserFromRoom(context.Context, *CSKickUserFromRoom) (*SCKickUserFromRoom, error)
	// LeaveRoom Leave room
	LeaveRoom(context.Context, *CSLeaveRoom) (*SCLeaveRoom, error)
	// RequestToJoinRoom Request to join room
	RequestToJoinRoom(context.Context, *CSRequestToJoinRoom) (*SCRequestToJoinRoom, error)
	// RoomDetail Room detail
	RoomDetail(context.Context, *CSRoomDetail) (*SCRoomDetail, error)
	// RoomList Room List
	RoomList(context.Context, *CSRoomList) (*SCRoomList, error)
}

func RegisterRoomServiceHTTPServer(s *http.Server, srv RoomServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/room/list", _RoomService_RoomList0_HTTP_Handler(srv))
	r.POST("/room/detail", _RoomService_RoomDetail0_HTTP_Handler(srv))
	r.POST("/room/create", _RoomService_CreateRoom0_HTTP_Handler(srv))
	r.POST("/room/invite_to_join", _RoomService_InviteToJoinRoom0_HTTP_Handler(srv))
	r.POST("/room/agree_to_invite_join", _RoomService_AgreeToInviteJoinRoom0_HTTP_Handler(srv))
	r.POST("/room/request_to_join", _RoomService_RequestToJoinRoom0_HTTP_Handler(srv))
	r.POST("/room/approve_request_to_join", _RoomService_ApproveRequestToJoinRoom0_HTTP_Handler(srv))
	r.POST("/room/kick_user", _RoomService_KickUserFromRoom0_HTTP_Handler(srv))
	r.POST("/room/leave", _RoomService_LeaveRoom0_HTTP_Handler(srv))
	r.POST("/room/close", _RoomService_CloseRoom0_HTTP_Handler(srv))
}

func _RoomService_RoomList0_HTTP_Handler(srv RoomServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSRoomList
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomServiceRoomList)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.RoomList(ctx, req.(*CSRoomList))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCRoomList)
		return ctx.Result(200, reply)
	}
}

func _RoomService_RoomDetail0_HTTP_Handler(srv RoomServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSRoomDetail
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomServiceRoomDetail)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.RoomDetail(ctx, req.(*CSRoomDetail))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCRoomDetail)
		return ctx.Result(200, reply)
	}
}

func _RoomService_CreateRoom0_HTTP_Handler(srv RoomServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSCreateRoom
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomServiceCreateRoom)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateRoom(ctx, req.(*CSCreateRoom))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCCreateRoom)
		return ctx.Result(200, reply)
	}
}

func _RoomService_InviteToJoinRoom0_HTTP_Handler(srv RoomServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSInviteToJoinRoom
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomServiceInviteToJoinRoom)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.InviteToJoinRoom(ctx, req.(*CSInviteToJoinRoom))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCInviteToJoinRoom)
		return ctx.Result(200, reply)
	}
}

func _RoomService_AgreeToInviteJoinRoom0_HTTP_Handler(srv RoomServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSAgreeToInviteJoinRoom
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomServiceAgreeToInviteJoinRoom)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AgreeToInviteJoinRoom(ctx, req.(*CSAgreeToInviteJoinRoom))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCAgreeToInviteJoinRoom)
		return ctx.Result(200, reply)
	}
}

func _RoomService_RequestToJoinRoom0_HTTP_Handler(srv RoomServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSRequestToJoinRoom
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomServiceRequestToJoinRoom)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.RequestToJoinRoom(ctx, req.(*CSRequestToJoinRoom))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCRequestToJoinRoom)
		return ctx.Result(200, reply)
	}
}

func _RoomService_ApproveRequestToJoinRoom0_HTTP_Handler(srv RoomServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSApproveRequestToJoinRoom
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomServiceApproveRequestToJoinRoom)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ApproveRequestToJoinRoom(ctx, req.(*CSApproveRequestToJoinRoom))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCApproveRequestToJoinRoom)
		return ctx.Result(200, reply)
	}
}

func _RoomService_KickUserFromRoom0_HTTP_Handler(srv RoomServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSKickUserFromRoom
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomServiceKickUserFromRoom)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.KickUserFromRoom(ctx, req.(*CSKickUserFromRoom))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCKickUserFromRoom)
		return ctx.Result(200, reply)
	}
}

func _RoomService_LeaveRoom0_HTTP_Handler(srv RoomServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSLeaveRoom
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomServiceLeaveRoom)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.LeaveRoom(ctx, req.(*CSLeaveRoom))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCLeaveRoom)
		return ctx.Result(200, reply)
	}
}

func _RoomService_CloseRoom0_HTTP_Handler(srv RoomServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSCloseRoom
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomServiceCloseRoom)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CloseRoom(ctx, req.(*CSCloseRoom))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCCloseRoom)
		return ctx.Result(200, reply)
	}
}

type RoomServiceHTTPClient interface {
	AgreeToInviteJoinRoom(ctx context.Context, req *CSAgreeToInviteJoinRoom, opts ...http.CallOption) (rsp *SCAgreeToInviteJoinRoom, err error)
	ApproveRequestToJoinRoom(ctx context.Context, req *CSApproveRequestToJoinRoom, opts ...http.CallOption) (rsp *SCApproveRequestToJoinRoom, err error)
	CloseRoom(ctx context.Context, req *CSCloseRoom, opts ...http.CallOption) (rsp *SCCloseRoom, err error)
	CreateRoom(ctx context.Context, req *CSCreateRoom, opts ...http.CallOption) (rsp *SCCreateRoom, err error)
	InviteToJoinRoom(ctx context.Context, req *CSInviteToJoinRoom, opts ...http.CallOption) (rsp *SCInviteToJoinRoom, err error)
	KickUserFromRoom(ctx context.Context, req *CSKickUserFromRoom, opts ...http.CallOption) (rsp *SCKickUserFromRoom, err error)
	LeaveRoom(ctx context.Context, req *CSLeaveRoom, opts ...http.CallOption) (rsp *SCLeaveRoom, err error)
	RequestToJoinRoom(ctx context.Context, req *CSRequestToJoinRoom, opts ...http.CallOption) (rsp *SCRequestToJoinRoom, err error)
	RoomDetail(ctx context.Context, req *CSRoomDetail, opts ...http.CallOption) (rsp *SCRoomDetail, err error)
	RoomList(ctx context.Context, req *CSRoomList, opts ...http.CallOption) (rsp *SCRoomList, err error)
}

type RoomServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewRoomServiceHTTPClient(client *http.Client) RoomServiceHTTPClient {
	return &RoomServiceHTTPClientImpl{client}
}

func (c *RoomServiceHTTPClientImpl) AgreeToInviteJoinRoom(ctx context.Context, in *CSAgreeToInviteJoinRoom, opts ...http.CallOption) (*SCAgreeToInviteJoinRoom, error) {
	var out SCAgreeToInviteJoinRoom
	pattern := "/room/agree_to_invite_join"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomServiceAgreeToInviteJoinRoom))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomServiceHTTPClientImpl) ApproveRequestToJoinRoom(ctx context.Context, in *CSApproveRequestToJoinRoom, opts ...http.CallOption) (*SCApproveRequestToJoinRoom, error) {
	var out SCApproveRequestToJoinRoom
	pattern := "/room/approve_request_to_join"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomServiceApproveRequestToJoinRoom))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomServiceHTTPClientImpl) CloseRoom(ctx context.Context, in *CSCloseRoom, opts ...http.CallOption) (*SCCloseRoom, error) {
	var out SCCloseRoom
	pattern := "/room/close"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomServiceCloseRoom))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomServiceHTTPClientImpl) CreateRoom(ctx context.Context, in *CSCreateRoom, opts ...http.CallOption) (*SCCreateRoom, error) {
	var out SCCreateRoom
	pattern := "/room/create"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomServiceCreateRoom))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomServiceHTTPClientImpl) InviteToJoinRoom(ctx context.Context, in *CSInviteToJoinRoom, opts ...http.CallOption) (*SCInviteToJoinRoom, error) {
	var out SCInviteToJoinRoom
	pattern := "/room/invite_to_join"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomServiceInviteToJoinRoom))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomServiceHTTPClientImpl) KickUserFromRoom(ctx context.Context, in *CSKickUserFromRoom, opts ...http.CallOption) (*SCKickUserFromRoom, error) {
	var out SCKickUserFromRoom
	pattern := "/room/kick_user"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomServiceKickUserFromRoom))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomServiceHTTPClientImpl) LeaveRoom(ctx context.Context, in *CSLeaveRoom, opts ...http.CallOption) (*SCLeaveRoom, error) {
	var out SCLeaveRoom
	pattern := "/room/leave"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomServiceLeaveRoom))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomServiceHTTPClientImpl) RequestToJoinRoom(ctx context.Context, in *CSRequestToJoinRoom, opts ...http.CallOption) (*SCRequestToJoinRoom, error) {
	var out SCRequestToJoinRoom
	pattern := "/room/request_to_join"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomServiceRequestToJoinRoom))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomServiceHTTPClientImpl) RoomDetail(ctx context.Context, in *CSRoomDetail, opts ...http.CallOption) (*SCRoomDetail, error) {
	var out SCRoomDetail
	pattern := "/room/detail"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomServiceRoomDetail))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomServiceHTTPClientImpl) RoomList(ctx context.Context, in *CSRoomList, opts ...http.CallOption) (*SCRoomList, error) {
	var out SCRoomList
	pattern := "/room/list"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoomServiceRoomList))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
