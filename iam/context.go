package iam

import (
	"context"
)

type contextkey string

const (
	contextKeyM2MProfile  = contextkey("m2m-profile")
	contextKeyUserProfile = contextkey("user-profile")
)

func SetM2MProfileInContext(ctx context.Context, profile M2MProfile) context.Context {
	return context.WithValue(ctx, contextKeyM2MProfile, profile)
}

func GetM2MProfileFromContext(ctx context.Context) M2MProfile {
	if p, ok := ctx.Value(contextKeyM2MProfile).(M2MProfile); ok {
		return p
	}

	return M2MProfile{}
}

func SetUserProfileInContext(ctx context.Context, profile UserProfile) context.Context {
	return context.WithValue(ctx, contextKeyUserProfile, profile)
}

func GetUserProfileFromContext(ctx context.Context) UserProfile {
	if p, ok := ctx.Value(contextKeyUserProfile).(UserProfile); ok {
		return p
	}

	return UserProfile{}
}
