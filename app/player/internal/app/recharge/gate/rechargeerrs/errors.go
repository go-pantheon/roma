package rechargeerrs

import "errors"

var (
	ErrProfileDev   = errors.New("profile error")
	ErrRechargeType = errors.New("recharge type error")
	ErrSignature    = errors.New("signature error")
	ErrUnmarshal    = errors.New("unmarshal error")
	ErrPackageName  = errors.New("package name error")
	ErrApiVerify    = errors.New("api verify error")
	ErrNotPurchased = errors.New("not purchased")
	ErrExisted      = errors.New("order existed")
	ErrAlreadyAck   = errors.New("already ack")
	ErrPending      = errors.New("order is pending")
	ErrCanceled     = errors.New("order canceled")
	ErrProductId    = errors.New("product id error")
	ErrAmountTooBig = errors.New("amount too big")
)
