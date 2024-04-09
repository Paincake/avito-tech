package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"strconv"
)

const (
	DefaultIdValue = -1
)

type GetBannerParams struct {
	FeatureId int64
	TagId     int64
	Limit     int
	Offset    int
}

func NewGetBannerParams(ctx echo.Context) (*GetBannerParams, error) {
	var err error
	var featureId int64
	featureId = DefaultIdValue
	var tagId int64
	tagId = DefaultIdValue
	limit := 0
	offset := 50
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

	return &GetBannerParams{
		FeatureId: featureId,
		TagId:     tagId,
		Limit:     limit,
		Offset:    offset,
	}, nil
}

type GetUserBannerParams struct {
	TagId        int64
	FeatureId    int64
	LastRevision bool
}

func NewGetUserBannerParams(ctx echo.Context) (*GetUserBannerParams, error) {
	var err error
	var featureId int64
	featureId = DefaultIdValue
	var tagId int64
	tagId = DefaultIdValue
	lastRevision := false
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

	return &GetUserBannerParams{
		TagId:        tagId,
		FeatureId:    featureId,
		LastRevision: lastRevision,
	}, nil
}

type PatchBannerIdParams struct {
	BannerId int64
}

func NewPatchBannerIdParams(ctx echo.Context) (*PatchBannerIdParams, error) {
	var err error
	var bannerId int64
	bannerId = DefaultIdValue
	param := ctx.QueryParams().Get("banner_id")
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
	param := ctx.QueryParams().Get("banner_id")
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
