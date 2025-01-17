package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/smtp"
	"time"

	"github.com/Afomiat/PRODIGY_FULL-STACK_INTERNSHIP/config"
	"github.com/Afomiat/PRODIGY_FULL-STACK_INTERNSHIP/domain"
	"github.com/Afomiat/PRODIGY_FULL-STACK_INTERNSHIP/internal/employeeUtil"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson/primitive"
// 
)

type SignupUsecase struct {
	signupRepo     domain.UserRepository
	contextTimeout time.Duration
	otpRepo        domain.OtpRepository
	env            *config.Env
}

func NewSignupUsecase(signupRepo domain.UserRepository, otpRepo domain.OtpRepository, timeout time.Duration, env *config.Env) *SignupUsecase {
	return &SignupUsecase{
		signupRepo:     signupRepo,
		contextTimeout: timeout,
		otpRepo:        otpRepo,
		env:            env,
	}
}

func (su *SignupUsecase) GetUserByUserName(ctx context.Context, username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	user, err := su.signupRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (su *SignupUsecase) GetUserByEmail(ctx context.Context, Email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	user, err := su.signupRepo.GetUserByEmail(ctx, Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (su *SignupUsecase) SendOtp(c context.Context, user *domain.User, smtpusername, smtppassword string) error {
    storedOTP, err := su.GetOtpByEmail(c, user.Email)
    if err != nil && err != mongo.ErrNoDocuments {
        return err
    }

    if storedOTP != nil {
        if time.Now().Before(storedOTP.ExpiresAt) {
            return errors.New("OTP already sent")
        }

        if err := su.otpRepo.DeleteOTP(c, storedOTP.Email); err != nil {
            return err
        }
    }

    otp := domain.OTP{
        Value:     employeeUtil.GenerateOTP(),
        Username:  user.Username,
        Email:     user.Email,
        Password:  user.Password,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(time.Minute * 5),
    }

    if err := su.otpRepo.SaveOTP(c, &otp); err != nil {
        return err
    }
	fmt.Println("OTP saved ******************************:", otp)

    if err := su.SendEmail(user.Email, otp.Value, smtpusername, smtppassword); err != nil {
        return err
    }

    return nil
}

func (su *SignupUsecase) GetOtpByEmail(ctx context.Context, email string) (*domain.OTP, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()
	
	return su.otpRepo.GetOtpByEmail(ctx, email)
}

func (su *SignupUsecase) SendEmail(email string, otpValue, smtpusername string, smtppassword string) error {
    from := smtpusername
    password := smtppassword

    to := []string{email}
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"
    message := []byte("Your OTP is " + otpValue)

    auth := smtp.PlainAuth("", from, password, smtpHost)
    val_emal := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)

	fmt.Println("Email sent:", val_emal)

	return val_emal
}

func (su *SignupUsecase) VerifyOtp(ctx context.Context, otp *domain.VerifyOtp) (*domain.OTP, error) {
    ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
    defer cancel()

    storedOTP, err := su.GetOtpByEmail(ctx, otp.Email)
    if err != nil || storedOTP == nil {
        if storedOTP == nil {
            return nil, errors.New("otp not found for the provided email. Please signup again.")
        }
        return nil, err
    }

    if storedOTP.Value != otp.Value {
        return nil, errors.New("invalid OTP.")
    }

    if time.Now().After(storedOTP.ExpiresAt) {
        return nil, errors.New("otp has expired.")
    }

    err = su.otpRepo.DeleteOTP(ctx, storedOTP.Email)
    if err != nil {
        return nil, err
    }

    return storedOTP, nil
}


func (su *SignupUsecase) RegisterUser(ctx context.Context, user *domain.SignupForm)(*primitive.ObjectID, error){
    ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
    defer cancel()

    hashedPass, err := employeeUtil.HassPassword(user.Password)

    if err != nil{
        return nil, err
    }

    addUser := domain.User{
        ID: primitive.NewObjectID(),
        Username: user.Username,
        Password: hashedPass,
        Email: user.Email,
        Role: domain.EmployeeRole,
    }
    err = su.signupRepo.CreateUser(ctx, &addUser)
    
    return &addUser.ID, err
}