// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             (unknown)
// source: account/interface/v1/account.proto

package interfacev1

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

const OperationAccountInterfaceAppleLogin = "/account.interface.v1.AccountInterface/AppleLogin"
const OperationAccountInterfaceAppleLoginCallback = "/account.interface.v1.AccountInterface/AppleLoginCallback"
const OperationAccountInterfaceDevPing = "/account.interface.v1.AccountInterface/DevPing"
const OperationAccountInterfaceFacebookLogin = "/account.interface.v1.AccountInterface/FacebookLogin"
const OperationAccountInterfaceGoogleLogin = "/account.interface.v1.AccountInterface/GoogleLogin"
const OperationAccountInterfaceLogin = "/account.interface.v1.AccountInterface/Login"
const OperationAccountInterfaceRefresh = "/account.interface.v1.AccountInterface/Refresh"
const OperationAccountInterfaceRegister = "/account.interface.v1.AccountInterface/Register"
const OperationAccountInterfaceToken = "/account.interface.v1.AccountInterface/Token"

type AccountInterfaceHTTPServer interface {
	// AppleLogin Apple login
	AppleLogin(context.Context, *AppleLoginRequest) (*AppleLoginResponse, error)
	// AppleLoginCallback Apple login callback
	AppleLoginCallback(context.Context, *AppleLoginCallbackRequest) (*AppleLoginCallbackResponse, error)
	// DevPing Connection test
	DevPing(context.Context, *DevPingRequest) (*DevPingResponse, error)
	// FacebookLogin Facebook login
	FacebookLogin(context.Context, *FacebookLoginRequest) (*FacebookLoginResponse, error)
	// GoogleLogin Google login
	GoogleLogin(context.Context, *GoogleLoginRequest) (*GoogleLoginResponse, error)
	// Login Login
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	// Refresh Session renewal
	Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error)
	// Register Register
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	// Token Get TCP handshake token
	Token(context.Context, *TokenRequest) (*TokenResponse, error)
}

func RegisterAccountInterfaceHTTPServer(s *http.Server, srv AccountInterfaceHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/dev/ping", _AccountInterface_DevPing0_HTTP_Handler(srv))
	r.POST("/v1/register", _AccountInterface_Register0_HTTP_Handler(srv))
	r.POST("/v1/login", _AccountInterface_Login0_HTTP_Handler(srv))
	r.POST("/v1/refresh", _AccountInterface_Refresh0_HTTP_Handler(srv))
	r.POST("/v1/token", _AccountInterface_Token0_HTTP_Handler(srv))
	r.POST("/v1/apple/login", _AccountInterface_AppleLogin0_HTTP_Handler(srv))
	r.POST("/v1/apple/login/callback", _AccountInterface_AppleLoginCallback0_HTTP_Handler(srv))
	r.POST("/v1/google/login", _AccountInterface_GoogleLogin0_HTTP_Handler(srv))
	r.POST("/v1/fb/login", _AccountInterface_FacebookLogin0_HTTP_Handler(srv))
}

func _AccountInterface_DevPing0_HTTP_Handler(srv AccountInterfaceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DevPingRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAccountInterfaceDevPing)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DevPing(ctx, req.(*DevPingRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DevPingResponse)
		return ctx.Result(200, reply)
	}
}

func _AccountInterface_Register0_HTTP_Handler(srv AccountInterfaceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in RegisterRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAccountInterfaceRegister)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Register(ctx, req.(*RegisterRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RegisterResponse)
		return ctx.Result(200, reply)
	}
}

func _AccountInterface_Login0_HTTP_Handler(srv AccountInterfaceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in LoginRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAccountInterfaceLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Login(ctx, req.(*LoginRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LoginResponse)
		return ctx.Result(200, reply)
	}
}

func _AccountInterface_Refresh0_HTTP_Handler(srv AccountInterfaceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in RefreshRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAccountInterfaceRefresh)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Refresh(ctx, req.(*RefreshRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RefreshResponse)
		return ctx.Result(200, reply)
	}
}

func _AccountInterface_Token0_HTTP_Handler(srv AccountInterfaceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in TokenRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAccountInterfaceToken)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Token(ctx, req.(*TokenRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*TokenResponse)
		return ctx.Result(200, reply)
	}
}

func _AccountInterface_AppleLogin0_HTTP_Handler(srv AccountInterfaceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AppleLoginRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAccountInterfaceAppleLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AppleLogin(ctx, req.(*AppleLoginRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AppleLoginResponse)
		return ctx.Result(200, reply)
	}
}

func _AccountInterface_AppleLoginCallback0_HTTP_Handler(srv AccountInterfaceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AppleLoginCallbackRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAccountInterfaceAppleLoginCallback)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AppleLoginCallback(ctx, req.(*AppleLoginCallbackRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AppleLoginCallbackResponse)
		return ctx.Result(200, reply)
	}
}

func _AccountInterface_GoogleLogin0_HTTP_Handler(srv AccountInterfaceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GoogleLoginRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAccountInterfaceGoogleLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GoogleLogin(ctx, req.(*GoogleLoginRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GoogleLoginResponse)
		return ctx.Result(200, reply)
	}
}

func _AccountInterface_FacebookLogin0_HTTP_Handler(srv AccountInterfaceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in FacebookLoginRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAccountInterfaceFacebookLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.FacebookLogin(ctx, req.(*FacebookLoginRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*FacebookLoginResponse)
		return ctx.Result(200, reply)
	}
}

type AccountInterfaceHTTPClient interface {
	AppleLogin(ctx context.Context, req *AppleLoginRequest, opts ...http.CallOption) (rsp *AppleLoginResponse, err error)
	AppleLoginCallback(ctx context.Context, req *AppleLoginCallbackRequest, opts ...http.CallOption) (rsp *AppleLoginCallbackResponse, err error)
	DevPing(ctx context.Context, req *DevPingRequest, opts ...http.CallOption) (rsp *DevPingResponse, err error)
	FacebookLogin(ctx context.Context, req *FacebookLoginRequest, opts ...http.CallOption) (rsp *FacebookLoginResponse, err error)
	GoogleLogin(ctx context.Context, req *GoogleLoginRequest, opts ...http.CallOption) (rsp *GoogleLoginResponse, err error)
	Login(ctx context.Context, req *LoginRequest, opts ...http.CallOption) (rsp *LoginResponse, err error)
	Refresh(ctx context.Context, req *RefreshRequest, opts ...http.CallOption) (rsp *RefreshResponse, err error)
	Register(ctx context.Context, req *RegisterRequest, opts ...http.CallOption) (rsp *RegisterResponse, err error)
	Token(ctx context.Context, req *TokenRequest, opts ...http.CallOption) (rsp *TokenResponse, err error)
}

type AccountInterfaceHTTPClientImpl struct {
	cc *http.Client
}

func NewAccountInterfaceHTTPClient(client *http.Client) AccountInterfaceHTTPClient {
	return &AccountInterfaceHTTPClientImpl{client}
}

func (c *AccountInterfaceHTTPClientImpl) AppleLogin(ctx context.Context, in *AppleLoginRequest, opts ...http.CallOption) (*AppleLoginResponse, error) {
	var out AppleLoginResponse
	pattern := "/v1/apple/login"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAccountInterfaceAppleLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *AccountInterfaceHTTPClientImpl) AppleLoginCallback(ctx context.Context, in *AppleLoginCallbackRequest, opts ...http.CallOption) (*AppleLoginCallbackResponse, error) {
	var out AppleLoginCallbackResponse
	pattern := "/v1/apple/login/callback"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAccountInterfaceAppleLoginCallback))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *AccountInterfaceHTTPClientImpl) DevPing(ctx context.Context, in *DevPingRequest, opts ...http.CallOption) (*DevPingResponse, error) {
	var out DevPingResponse
	pattern := "/v1/dev/ping"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAccountInterfaceDevPing))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *AccountInterfaceHTTPClientImpl) FacebookLogin(ctx context.Context, in *FacebookLoginRequest, opts ...http.CallOption) (*FacebookLoginResponse, error) {
	var out FacebookLoginResponse
	pattern := "/v1/fb/login"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAccountInterfaceFacebookLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *AccountInterfaceHTTPClientImpl) GoogleLogin(ctx context.Context, in *GoogleLoginRequest, opts ...http.CallOption) (*GoogleLoginResponse, error) {
	var out GoogleLoginResponse
	pattern := "/v1/google/login"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAccountInterfaceGoogleLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *AccountInterfaceHTTPClientImpl) Login(ctx context.Context, in *LoginRequest, opts ...http.CallOption) (*LoginResponse, error) {
	var out LoginResponse
	pattern := "/v1/login"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAccountInterfaceLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *AccountInterfaceHTTPClientImpl) Refresh(ctx context.Context, in *RefreshRequest, opts ...http.CallOption) (*RefreshResponse, error) {
	var out RefreshResponse
	pattern := "/v1/refresh"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAccountInterfaceRefresh))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *AccountInterfaceHTTPClientImpl) Register(ctx context.Context, in *RegisterRequest, opts ...http.CallOption) (*RegisterResponse, error) {
	var out RegisterResponse
	pattern := "/v1/register"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAccountInterfaceRegister))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *AccountInterfaceHTTPClientImpl) Token(ctx context.Context, in *TokenRequest, opts ...http.CallOption) (*TokenResponse, error) {
	var out TokenResponse
	pattern := "/v1/token"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAccountInterfaceToken))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
