package dto

import (
	"fmt"
	//"github.com/Paincake/avito-tech/internal/server"
	"github.com/labstack/echo/v4"
	"strconv"
)

const (
	DefaultIdValue      = -1
	TokenRoleContextKey = "Token"
)

type GetBannerParams struct {
	FeatureId int64
	TagId     int64
	Limit     int
	Offset    int
	UseActive bool
}

func NewGetBannerParams(ctx echo.Context) (*GetBannerParams, error) {
	var err error
	featureId := int64(DefaultIdValue)
	tagId := int64(DefaultIdValue)
	limit := 50
	offset := 0
	useActive := true
	param := ctx.QueryParams().Get("feature_id")
	if param != "" {
		featureId, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid feature_id format: %s", err)
		}
	}
	param = ctx.QueryParams().Get("tag_id")
	if param != "" {
		tagId, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid tag_id format: %s", err)
		}
	}
	param = ctx.QueryParams().Get("limit")
	if param != "" {
		limit, err = strconv.Atoi(param)
		if err != nil {
			return nil, fmt.Errorf("invalid limit format: %s", err)
		}
	}
	param = ctx.QueryParams().Get("offset")
	if param != "" {
		offset, err = strconv.Atoi(param)
		if err != nil {
			return nil, fmt.Errorf("invalid offset format: %s", err)
		}
	}

	param, ok := ctx.Get(TokenRoleContextKey).(string)
	if ok {
		if param == "admin" {
			useActive = false
		}
	}

	return &GetBannerParams{
		FeatureId: featureId,
		TagId:     tagId,
		Limit:     limit,
		Offset:    offset,
		UseActive: useActive,
	}, nil
}

type GetUserBannerParams struct {
	TagId        int64
	FeatureId    int64
	LastRevision bool
	UseActive    bool
}

func NewGetUserBannerParams(ctx echo.Context) (*GetUserBannerParams, error) {
	var err error
	featureId := int64(DefaultIdValue)
	tagId := int64(DefaultIdValue)
	lastRevision := false
	useActive := true
	param := ctx.QueryParams().Get("feature_id")
	if param == "" {
		return nil, fmt.Errorf("missed required query param: feature_id")
	} else {
		featureId, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid feature_id format: %s", err)
		}
	}
	param = ctx.QueryParams().Get("tag_id")
	if param == "" {
		return nil, fmt.Errorf("missed required query param: tag_id")
	} else {
		tagId, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid tag_id format: %s", err)
		}
	}
	param = ctx.QueryParams().Get("use_last_revision")
	if param != "" {
		lastRevision, err = strconv.ParseBool(param)
		if err != nil {
			return nil, fmt.Errorf("invalid use_last_revision format")
		}
	}
	param, ok := ctx.Get(TokenRoleContextKey).(string)
	if ok {
		if param == "admin" {
			useActive = false
		}
	}
	return &GetUserBannerParams{
		TagId:        tagId,
		FeatureId:    featureId,
		LastRevision: lastRevision,
		UseActive:    useActive,
	}, nil
}

type PatchBannerIdParams struct {
	BannerId int64
}

func NewPatchBannerIdParams(ctx echo.Context) (*PatchBannerIdParams, error) {
	var err error
	var bannerId int64
	bannerId = DefaultIdValue
	param := ctx.Param("id")
	if param == "" {
		return nil, fmt.Errorf("missed required query param: banner_id")
	} else {
		bannerId, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid banner_id format: %s", err)
		}
	}

	return &PatchBannerIdParams{
		BannerId: bannerId,
	}, nil
}

type DeleteBannerIdParams struct {
	BannerId int64
}

func NewDeleteBannerIdParams(ctx echo.Context) (*DeleteBannerIdParams, error) {
	var err error
	var bannerId int64
	bannerId = DefaultIdValue
	param := ctx.Param("id")
	if param == "" {
		return nil, fmt.Errorf("missed required query param: delete_id")
	} else {
		bannerId, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid delete_id format: %s", err)
		}
	}

	return &DeleteBannerIdParams{
		BannerId: bannerId,
	}, nil
}
