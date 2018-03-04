package controllers

import (
	"encoding/json"

	"panda/models"
	"strconv"

	"github.com/astaxie/beego"
)

// AttributeController operations for Attribute
type AttributeController struct {
	beego.Controller
}

// URLMapping ...
func (c *AttributeController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Attribute
// @Param	body		body 	models.Attribute	true		"body for Attribute content"
// @Success 201 {int} models.Attribute
// @Failure 403 body is empty
// @router / [post]
func (c *AttributeController) Post() {
	var v models.Attribute
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddAttribute(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Attribute by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Attribute
// @Failure 403 :id is empty
// @router /:id [get]
func (c *AttributeController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	v, err := models.GetAttributeById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Attribute
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Attribute	true		"body for Attribute content"
// @Success 200 {object} models.Attribute
// @Failure 403 :id is not int
// @router /:id [put]
func (c *AttributeController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	v := models.Attribute{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateAttributeById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Attribute
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *AttributeController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := models.DeleteAttribute(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
