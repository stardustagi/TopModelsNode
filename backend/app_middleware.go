package backend

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/jwt"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/utils"
	"github.com/stardustagi/TopModelsNode/constants"
	"go.uber.org/zap"
)

// NodeAccess 节点访问
func NodeAccess() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//在这里处理拦截请求的逻辑

			redisCmd := redis.GetRedisDb()
			nodeIdString := c.Request().Header.Get("nodeId")
			jwtStr := c.Request().Header.Get("jwt")
			accessKey := c.Request().Header.Get("accessKey")
			once := c.Request().Header.Get("once")
			// 取nodeKey
			nodeId, err := strconv.ParseInt(nodeIdString, 10, 64)
			if err != nil {
				return c.JSON(401, map[string]any{
					"errcode": 2,
					"errmsg":  "Node ID格式错误",
				})
			}
			nodeAccessKey := fmt.Sprintf("%s:%s", constants.ApplicationPrefix, constants.NodeAccessModelsKey(nodeId))
			secretKey, err := redisCmd.HGet(context.Background(), nodeAccessKey, "securityKey").Result()
			if err != nil && errors.Is(err, redis.Nil) {
				return c.JSON(401, map[string]any{
					"errcode": 2,
					"errmsg":  "Node 不存在，请重新注册",
				})
			}
			if secretKey == "" {
				return c.JSON(401, map[string]any{
					"errcode": 2,
					"errmsg":  "Node 信息错误，请联系管理员",
				})
			}
			secret := fmt.Sprintf("%s-%s-%s-%s", accessKey, secretKey, once, nodeIdString)
			if jwtStr != "" {
				jwtObj, ok := jwt.JWTDecrypt(jwtStr, secret)
				if !ok || jwtObj == nil || jwtObj["token"] == nil || jwtObj["id"] == nil {
					return c.JSON(401, map[string]any{
						"errcode": 2,
						"errmsg":  "jwt解析错误",
					})
				}

				id, ok := jwtObj["id"].(string)
				if !ok {
					return c.JSON(401, map[string]any{
						"errcode": 2,
						"errmsg":  "用户信息获取失败",
					})
				}
				oldToken, err1 := redisCmd.HGet(context.Background(), nodeAccessKey, "token").Result()
				if err1 != nil {
					return c.JSON(401, map[string]any{
						"errcode": 2,
						"errmsg":  "获取token失败",
					})
				}
				if jwtObj["token"] != oldToken {
					return c.JSON(401, map[string]any{
						"errcode": 2,
						"errmsg":  "账户已经在其他终端登录",
					})
				}
				c.QueryParams().Add("id", id)
			} else {
				ip := echo.ExtractIPFromXFFHeader()(c.Request())
				// 判断是否内网IP
				ok := utils.IsInnerIp(ip)
				if !ok {
					return c.JSON(401, map[string]any{
						"errcode": 2,
						"errmsg":  "服务器繁忙,请稍后再试!",
					})
				}
			}
			return next(c)
		}
	}
}

// NodeUserAccess 节点用户访问
func NodeUserAccess() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		//在这里处理拦截请求的逻辑
		return func(c echo.Context) error {
			//在这里处理拦截请求的逻辑
			jwtstr := c.Request().Header.Get("jwt")
			nodeUserId := c.Request().Header.Get("id")
			secret := fmt.Sprintf("%s-%s-%s", constants.AppName, constants.AppVersion, nodeUserId)
			if jwtstr != "" {
				jwtobj, ok := jwt.JWTDecrypt(jwtstr, secret)
				if !ok || jwtobj == nil || jwtobj["token"] == nil || jwtobj["id"] == nil {
					return c.JSON(401, map[string]any{
						"errcode": 2,
						"errmsg":  "jwt解析错误",
					})
				}

				id, ok := jwtobj["id"].(string)
				if !ok {
					return c.JSON(401, map[string]any{
						"errcode": 2,
						"errmsg":  "用户信息获取失败",
					})
				}
				intId, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					return c.JSON(401, map[string]any{
						"errcode": 2,
						"errmsg":  "用户信息获取失败",
					})
				}
				tokenKey := fmt.Sprintf("%s:%s:user:%s", constants.AppName, constants.AppVersion, constants.NodeUserTokenKey(intId))
				jwtKey := fmt.Sprintf("%s:%s", constants.ModelsKeyPrefix, constants.NodeUserJwtKey(intId))
				logs.Info("redis key is:", zap.String("tokenKey", tokenKey), zap.String("jwtKey", jwtKey))
				redisCmd := redis.GetRedisDb()
				oldToken, err1 := redisCmd.Get(context.Background(), tokenKey).Result()
				if err1 != nil {
					return c.JSON(401, map[string]any{
						"errcode": 2,
						"errmsg":  "获取token失败",
					})
				}
				if jwtobj["token"] != oldToken {
					return c.JSON(401, map[string]any{
						"errcode": 2,
						"errmsg":  "账户已经在其他终端登录",
					})
				}
				//}
				c.QueryParams().Add("id", id)
			} else {
				ip := echo.ExtractIPFromXFFHeader()(c.Request())
				// 判断是否内网IP
				ok := utils.IsInnerIp(ip)
				if !ok {
					return c.JSON(401, map[string]any{
						"errcode": 2,
						"errmsg":  "服务器繁忙,请稍后再试!",
					})
				}
			}
			return next(c)
		}
	}
}
