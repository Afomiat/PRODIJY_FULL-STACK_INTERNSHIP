package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignupUsecase interface {
	VerifyOtp(ctx context.Context, otp *VerifyOtp) (*OTP, error)
	RegisterUser(ctx context.Context, user *SignupForm) (*primitive.ObjectID, error)
	GetUserByUserName(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, Email string) (*User, error)
	SendOtp(cxt context.Context, user *User, stmpName, stmpPass string) error
}



type OtpRepository interface {
	GetOtpByEmail(ctx context.Context, email string) (*OTP, error)
	DeleteOTP(ctx context.Context, email string) error
	SaveOTP(ctx context.Context, otp *OTP) error
}

type LoginUsecase interface {
	AuthenticateUser(c context.Context, login *AuthLogin) (*User, error)
	CreateAccessToken(user *User, secret string, expiry int) (string, error)
	CreateRefreshToken(user *User, secret string, expiry int) (string, error)
	SaveRefreshToken(c context.Context, token *Token) error
	CheckRefreshToken(c context.Context, refreshToken string) (*Token, error)
}

type UserRepository interface {
	GetUserByID(c context.Context, id primitive.ObjectID) (*User, error)
	GetUserByEmail(c context.Context, email string) (*User, error)
	GetUserByUsername(c context.Context, username string) (*User, error)
	// GetUserByUserName(ctx context.Context, username string) (*SignupForm, error)

	GetAllUsers(c context.Context) ([]*User, error)
	CreateUser(c context.Context, user *User) error
	UpdateUser(c context.Context, user *User) error

	DeleteUser(c context.Context, id primitive.ObjectID) error
}

type UserUsecase interface {
	GetUserByID(c context.Context, id primitive.ObjectID) (*User, error)
	GetUserByEmail(c context.Context, email string) (*User, error)
	GetUserByUsername(c context.Context, username string) (*User, error)
	GetAllUsers(c context.Context) ([]*User, error)
	CreateUser(c context.Context, user *SignupForm, claims *JwtCustomClaims) error

	UpdateUser(c context.Context, user *User,claims *JwtCustomClaims, existinguser *User) (*User, error)
	DeleteUser(c context.Context, id primitive.ObjectID,claims *JwtCustomClaims) error
}
type TokenRepository interface {
	SaveToken(ctx context.Context, token *Token) error

	FindTokenByRefreshToken(ctx context.Context, refreshToken string) (*Token, error)
}




type AttendanceUsecase interface {
    ClockIn(ctx context.Context, userID primitive.ObjectID) error
    ClockOut(ctx context.Context, userID primitive.ObjectID) error
    GetAllAttendanceRecords(ctx context.Context) ([]AttendanceRecord, error)
}

type AttendanceRepository interface {
    InsertAttendanceRecord(ctx context.Context, record *AttendanceRecord) error
    FindLatestClockInRecord(ctx context.Context, userID primitive.ObjectID) (*AttendanceRecord, error)
    UpdateAttendanceRecord(ctx context.Context, record *AttendanceRecord) error
    GetAllRecords(ctx context.Context) ([]AttendanceRecord, error)
}
