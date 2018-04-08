package controllers

import (
	"fmt"
	"panda/models"
	"panda/types"

	"math/rand"
	"time"

	"bufio"
	"os"

	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

type GeneratesvgfileController struct {
	beego.Controller
}

/*get random number*/
func generate_rand(number int64) (nRand int64) {
	if number <= 0 {
		return -1
	}

	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	return int64(seed.Intn(int(number)))
}

/*Get svg_dtl from svg_info by s_id*/
func getSvgDetailBySId(s_id int64) (svg_detail string) {
	query := make(map[string]string, 0)
	query["s_id"] = fmt.Sprintf("%v", s_id)
	resultl, _ := models.GetAllSvginfo(query, []string{}, []string{}, []string{}, 0, 1)
	for _, v := range resultl {
		svg_detail = v.(models.Svg_info).Svg_dtl
		if v.(models.Svg_info).N_id != 0 {
			// RECURSION
			svg_detail += getSvgDetailBySId(v.(models.Svg_info).N_id)
		}
	}
	return svg_detail
}

/*Get svg information by catagory id, color(random, if have) and index number(random)*/
func getSvgDetail(catagory_id int64, color_flag int, color int64, index int64, ids []int64) (s_id int64, svg_detail string) {
	link_flag := 0
	svg_detail = ""

	query := make(map[string]string, 0)
	if color_flag == 1 {
		query["base_color"] = fmt.Sprintf("%v", color)
	}
	query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
	query["p_id"] = fmt.Sprintf("%v", 0)
	resultl, _ := models.GetAllSvginfo(query, []string{}, []string{"s_id"}, []string{"asc"}, index, 1)

	for _, v := range resultl {
		s_id = v.(models.Svg_info).S_id
		// This item can be used only "Link_id item" exists.
		if v.(models.Svg_info).Link_id != "" {
			for _, id := range ids {
				for _, vlink_id := range strings.Split(v.(models.Svg_info).Link_id, ",") {
					link_id, err := strconv.Atoi(vlink_id)
					if err != nil {
						fmt.Println("ATOI Failed")
					}
					if id == int64(link_id) { // MATCH! CAN BE USE
						link_flag = 1
						break
					}
				}
				if link_flag == 1 {
					break
				}
			}
			if link_flag != 1 {
				break
			}
		}
		svg_detail += v.(models.Svg_info).Svg_dtl
		// Link the next svg to be strcat
		if v.(models.Svg_info).N_id != 0 {
			svg_detail += getSvgDetailBySId(v.(models.Svg_info).N_id)
		}
	}
	return s_id, svg_detail
}

func (c *GeneratesvgfileController) HandlerGenerate() {

	conf := GetConfigData()

	path := c.Generate_svg(1, beego.AppConfig.String("svg_path"), "1")
	//path := c.Generate_svg(0, "c://", "1")

	c.Ctx.WriteString(conf.HostUrl + types.Svg_File_Path + "/" + path)
}

/*Generate svg file*/
func (c *GeneratesvgfileController) Generate_svg(flag int, basePath string, petID string) (svgPath string) {

	//svg head
	svg := "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"-100 25 800 800\">"

	//random for panda base color
	color := generate_rand(7)

	// ids
	cg_count, err0 := models.GetSvgCatagoryCountByField("id")
	if err0 != nil {
		fmt.Println(err0)
	}
	// memset
	var ids = make([]int64, cg_count, 20)

	// Svg Select_flag count
	cg_count_sf, err1 := models.GetSvgCatagoryCountByField("Select_flag")
	if err1 != nil {
		fmt.Println(err1)
	}
	// memset
	var selectflag_times = make([]int, cg_count_sf, 10)
	var selectflag_check = make([]int, cg_count_sf, 10)

	for i, _ := range selectflag_check {
		if i > 1 {
			count, _ := models.GetCountBySelectId(int64(i))
			selectflag_check[i] = int(generate_rand(count))
		}
	}

	//query structController
	query := make(map[string]string, 0)

	//Get svg_catatory order by rank
	ml, err := models.GetAllSvgcata(query, []string{}, []string{"rank"}, []string{"asc"}, 0, 20)
	if err != nil {
		fmt.Println(err)
	}

	for i, v := range ml {
		catagory_id := v.(models.Svg_catagory).Id
		select_flag := v.(models.Svg_catagory).Select_flag
		bodycolor_flag := v.(models.Svg_catagory).Bodycolor_flag
		percent := v.(models.Svg_catagory).Probability

		switch select_flag {
		case 0:
			//use random color
			if bodycolor_flag == 1 {
				// eye, ear, leg, hand
				rand := generate_rand(models.GetCountByCatagoryId(catagory_id) / 7)
				if rand == -1 {
					break
				}
				_id, _svg := getSvgDetail(catagory_id, 1, color, rand, ids)
				ids[i] = _id
				svg += _svg
			} else {
				// bodyline, mouth
				rand := generate_rand(models.GetCountByCatagoryId(catagory_id))
				if rand == -1 {
					break
				}
				_id, _svg := getSvgDetail(catagory_id, 0, 0, rand, ids)
				ids[i] = _id
				svg += _svg
			}

		case 1:
			//use random element
			//Be or not be
			if percent == 0 {
				goto platform
			}

			if percent > 50 && percent < 100 { // usual items
				rand_flag := generate_rand(int64(100 / (100 - percent)))
				if rand_flag == 0 {
					break
				}
			} else { //	unusual	items
				rand_flag := generate_rand(int64(100 / percent))
				if rand_flag != 0 {
					break
				}
			}

		platform:
			if (percent != 0) || (flag == 1 && percent == 0) {
				//get count of this item
				rand := generate_rand(models.GetCountByCatagoryId(catagory_id))
				if rand == -1 {
					break
				}
				_id, _svg := getSvgDetail(catagory_id, 0, 0, rand, ids)
				ids[i] = _id
				svg += _svg
			}

		default:
			//choose one; drop another
			//How many items where select_flag = this, calculator only first time.
			if selectflag_times[select_flag] == selectflag_check[select_flag] {
				if bodycolor_flag == 1 { //leg
					rand := generate_rand(models.GetCountByCatagoryId(catagory_id) / 7)
					if rand == -1 {
						break
					}
					_id, _svg := getSvgDetail(catagory_id, 1, color, rand, ids)
					ids[i] = _id
					svg += _svg
				} else { // hat && front_hair
					rand := generate_rand(models.GetCountByCatagoryId(catagory_id))
					if rand == -1 {
						break
					}
					_id, _svg := getSvgDetail(catagory_id, 0, 0, rand, ids)
					ids[i] = _id
					svg += _svg
				}
			}
			selectflag_times[select_flag]++
		}
	}

	svg += "</svg>"

	// Create svg file and write inside
	//	if basePath[0:1] != "\\" {
	//		basePath += "\\"
	//	}
	fileName := fmt.Sprintf("%s%v.svg", time.Now().Format("20060102150405"), petID)
	strFile := basePath + fileName
	f, err := os.Create(strFile)
	w := bufio.NewWriter(f)
	if _, err = w.WriteString(svg); err != nil {
		//err deal
		fmt.Println(err)
	}
	w.Flush()

	f.Close()

	return fileName
}
