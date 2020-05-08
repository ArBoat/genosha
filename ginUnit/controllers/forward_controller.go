package controllers

import (
	"bytes"
	"genosha/ginUnit"
	"genosha/utils/confs"
	"genosha/utils/myLogger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type ForwardController struct{}

func NewForwardController() *ForwardController {
	return &ForwardController{}
}

func (fc ForwardController) GetSimpleForward(c *gin.Context) {
	var sUrl, rqUrl string
	rPath := c.Request.URL.Path
	rQuery := c.Request.URL.RawQuery
	if rQuery != "" {
		sUrl = rPath + `?` + rQuery
		rqUrl = confs.ConfigMap["server_url"] + sUrl + `&`
	} else {
		sUrl = rPath
		rqUrl = confs.ConfigMap["server_url"] + sUrl + `?`
	}
	paras, err := ginUnit.keySign(c, "GET", nil, sUrl)
	if err != nil {
		myLogger.Log.Error("error", zap.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "internal err2",
			"code": 5000,
		})
		return
	}
	rqUrl = rqUrl + paras
	myLogger.Log.Info("request url of query reports' details:" + rqUrl)

	resp, err := http.Get(rqUrl)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		myLogger.Log.Error("error", zap.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "server request filed",
			"code": 5000,
		})
		myLogger.SendErrorEmail(err.Error())
		return
	}
	myLogger.Log.Info("resp.StatusCode", zap.Any("resp.StatusCode", resp.StatusCode))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		myLogger.Log.Error("error", zap.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "internal err",
			"code": 5000,
		})
		return
	}
	myLogger.Log.Info("resp.Body:" + string(body))
	ginUnit.handleResp(resp, c, body, resp.Header)
}

func (fc ForwardController) PostSimpleForward(c *gin.Context) {
	var sUrl, rqUrl string
	rPath := c.Request.URL.Path
	rBody := c.Request.Body

	rBodyByte, _ := ioutil.ReadAll(rBody)
	myLogger.Log.Info("requestBody" + string(rBodyByte))
	sUrl = rPath
	rqUrl = confs.ConfigMap["server_url"] + sUrl + `?`
	paras, err := ginUnit.keySign(c, "POST", rBodyByte, sUrl)
	if err != nil {
		myLogger.Log.Error("error", zap.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "internal err2",
			"code": 5000,
		})
		return
	}
	rqUrl = rqUrl + paras
	myLogger.Log.Info("request url of query reports' details:" + rqUrl)

	resp, err := http.Post(rqUrl, "application/json;charset=utf-8", bytes.NewBuffer(rBodyByte))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		myLogger.Log.Error("error", zap.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "server request filed",
			"code": 5000,
		})
		myLogger.SendErrorEmail(err.Error())
		return
	}
	myLogger.Log.Info("resp.StatusCode", zap.Any("resp.StatusCode", resp.StatusCode))
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		myLogger.Log.Error("error", zap.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "internal err",
			"code": 5000,
		})
		return
	}
	myLogger.Log.Info("resp.Body:" + string(body))
	ginUnit.handleResp(resp, c, body, resp.Header)
}

func (fc ForwardController) DeleteSimpleForward(c *gin.Context) {
	var sUrl, rqUrl string
	rPath := c.Request.URL.Path
	sUrl = rPath
	rqUrl = confs.ConfigMap["server_url"] + sUrl + `?`
	paras, err := ginUnit.keySign(c, "DELETE", nil, sUrl)
	if err != nil {
		myLogger.Log.Error("error", zap.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "internal err2",
			"code": 5000,
		})
		return
	}
	rqUrl = rqUrl + paras

	myLogger.Log.Info("request url of query reports' details:" + rqUrl)

	req, err := http.NewRequest("DELETE", rqUrl, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		myLogger.Log.Error("error", zap.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "server request filed",
			"code": 5000,
		})
		myLogger.SendErrorEmail(err.Error())
		return
	}
	myLogger.Log.Info("resp.StatusCode", zap.Any("resp.StatusCode", resp.StatusCode))
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		myLogger.Log.Error("error", zap.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "internal err",
			"code": 5000,
		})
		return
	}
	myLogger.Log.Info("resp.Body:" + string(body))
	ginUnit.handleResp(resp, c, body, resp.Header)
}

func sevenDayString(from string) string {
	tm, _ := time.Parse("2006-01-02", from)
	return tm.AddDate(0, 0, 6).Format("2006-01-02")
}
