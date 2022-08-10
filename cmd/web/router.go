package main

import (
	. "github.com/gin-gonic/gin"
	"gomixer/pkg/audio"
	"net/http"
)

func InitRouter() *Engine {
	r := Default()

	r.GET("/audio/devices", func(c *Context) {
		deviceType := audio.DeviceType(c.DefaultQuery("type", "all"))

		if err := audio.ValidateDeviceType(deviceType); err != nil {
			c.JSON(http.StatusBadRequest, H{
				"error": err.Error(),
			})
			return
		}

		var err error
		var deviceList []audio.Device
		if deviceList, err = audio.ListDevices(deviceType); err != nil {
			c.JSON(http.StatusInternalServerError, H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, H{
			"devices": deviceList,
		})
	})

	r.GET("/audio/:device", func(c *Context) {
		deviceId := c.Param("device")

		var device audio.Device
		var err error

		device, err = audio.GetDevice(deviceId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, device)
	})

	r.POST("/audio/:device", func(c *Context) {
		deviceId := c.Param("device")

		var device audio.Device
		var data audio.DeviceUpdate
		var err error

		if err := c.Bind(&data); err != nil {
			c.JSON(http.StatusInternalServerError, H{
				"error": err.Error(),
			})
			return
		}

		if device, err = audio.UpdateDevice(deviceId, data); err != nil {
			c.JSON(http.StatusInternalServerError, H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, device)
	})

	return r
}
