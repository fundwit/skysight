package localize

import (
	"net/http"
	"net/http/httptest"
	"skysight/testinfra"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
)

func TestLocalize(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should return origin text if message key not exist", func(t *testing.T) {
		r := gin.Default()
		r.Use(LocalizeMiddleware("../i18n"))

		r.GET("/", func(c *gin.Context) {
			msg := MustGetMessage("running")
			c.String(http.StatusOK, msg)
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(Equal("running"))
	})

	t.Run("should return origin text if category not exist", func(t *testing.T) {
		r := gin.Default()
		r.Use(LocalizeMiddleware("../i18n"))

		r.GET("/", func(c *gin.Context) {
			msg := MustGetMessage("running")
			c.String(http.StatusOK, msg)
		})
		req := httptest.NewRequest(http.MethodGet, "/?lang=xxx", nil)
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(Equal("running"))
	})

	t.Run("should return default text if category not exist (query)", func(t *testing.T) {
		r := gin.Default()
		r.Use(LocalizeMiddleware("../i18n"))

		r.GET("/", func(c *gin.Context) {
			msg := MustGetMessage("status-running")
			c.String(http.StatusOK, msg)
		})
		req := httptest.NewRequest(http.MethodGet, "/?lang=xxx", nil)
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(Equal("running"))
	})

	t.Run("should return default text if category not exist (header)", func(t *testing.T) {
		r := gin.Default()
		r.Use(LocalizeMiddleware("../i18n"))

		r.GET("/", func(c *gin.Context) {
			msg := MustGetMessage("status-running")
			c.String(http.StatusOK, msg)
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("Accept-Language", "xxx")
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(Equal("running"))
	})
}

func TestLocalizeSpecifiedLanguage(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should return specified text by query string lang", func(t *testing.T) {
		r := gin.Default()
		r.Use(LocalizeMiddleware("../i18n"))

		r.GET("/", func(c *gin.Context) {
			msg := MustGetMessage("status-running")
			c.String(http.StatusOK, msg)
		})
		req := httptest.NewRequest(http.MethodGet, "/?lang=en", nil)
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(Equal("running"))
	})

	t.Run("should return specified text by header 'Accept-Language' (query first)", func(t *testing.T) {
		r := gin.Default()
		r.Use(LocalizeMiddleware("../i18n"))

		key := "status-running"
		r.GET("/", func(c *gin.Context) {
			msg, _ := GetMessage(key)
			c.String(http.StatusOK, msg)
		})
		req := httptest.NewRequest(http.MethodGet, "/?lang=en", nil)
		req.Header.Add("Accept-Language", "zh")
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(Equal("running"))

		req = httptest.NewRequest(http.MethodGet, "/?lang=zh", nil)
		req.Header.Add("Accept-Language", "en")
		status, body, _ = testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(Equal("运行中"))

		req = httptest.NewRequest(http.MethodGet, "/?lang=aaa", nil)
		req.Header.Add("Accept-Language", "bbb")
		status, body, _ = testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(Equal("running")) // default language
	})

	t.Run("should return specified text by header 'Accept-Language' with weight", func(t *testing.T) {
		r := gin.Default()
		r.Use(LocalizeMiddleware("../i18n"))

		key := "status-running"
		r.GET("/", func(c *gin.Context) {
			msg, _ := GetMessage(key)
			c.String(http.StatusOK, msg)
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("Accept-Language", "zh-CN,xx;q=0.5")
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(Equal("运行中"))

		req = httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("Accept-Language", "zh-CN;q=0.4,en;q=0.6")
		status, body, _ = testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(Equal("running"))

		req = httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("Accept-Language", "*")
		status, body, _ = testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(Equal("running")) // default language
	})
}
